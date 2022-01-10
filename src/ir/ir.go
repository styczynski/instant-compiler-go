package ir

import (
	"fmt"
	"reflect"

	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type IRGeneratorState struct {
	tempCounter      int
	tempCounterBlock map[int]int
	temps            map[string]struct{}
}

func (ir *IRGeneratorState) isTempVar(name string) bool {
	if _, ok := ir.temps[name]; ok {
		return true
	}
	return false
}

func (ir *IRGeneratorState) NextVar(blockID int, prefix string) string {
	if _, ok := ir.tempCounterBlock[blockID]; !ok {
		ir.tempCounterBlock[blockID] = 0
	}
	newName := fmt.Sprintf("%s_%d_%d", prefix, blockID, ir.tempCounterBlock[blockID])
	ir.tempCounterBlock[blockID]++
	return newName
}

func (ir *IRGeneratorState) NextTempVar() string {
	ir.tempCounter++
	newName := fmt.Sprintf("temp_%d", ir.tempCounter)
	ir.temps[newName] = struct{}{}
	return newName
}

func translateType(t hindley_milner.Type) IRType {
	resolvedType := IR_UNKNOWN
	if _, ok := t.(*hindley_milner.FunctionType); ok {
		resolvedType = IR_FN
	} else if primitive, ok := t.(ast.PrimitiveType); ok {
		if primitive.Name() == "int" {
			resolvedType = IR_INT32
		} else if primitive.Name() == "boolean" {
			resolvedType = IR_BIT
		}
	}
	return resolvedType
}

func generateIRExpr(c *context.ParsingContext, ir *IRGeneratorState, node generic_ast.Expression) ([]*IRStatement, IRType, string) {
	ret := []*IRStatement{}
	resultVar := ir.NextTempVar()

	if e, ok := node.(*ast.Primary); ok {
		if e.IsVariable() {
			return []*IRStatement{}, translateType(e.ResolvedType), *e.Variable
		} else if e.IsInt() {
			ret = append(ret, WrapIRExpression(&IRExpression{
				BaseASTNode: e.BaseASTNode,
				TargetName:  resultVar,
				Operation:   "Load",
				Type:        translateType(e.ResolvedType),
				Arguments: []string{
					e.Print(c),
				},
				ArgumentsTypes: []IRType{translateType(e.ResolvedType)},
			}))
			return ret, translateType(e.ResolvedType), resultVar
		}
	} else if e, ok := node.(*ast.Index); ok {
		if !e.HasIndexingExpr() {
			return generateIRExpr(c, ir, e.Primary)
		}
	} else if e, ok := node.(*ast.UnaryApplication); ok {
		if e.IsIndex() {
			return generateIRExpr(c, ir, e.Index)
		} else if e.IsApplication() {
			argList := []string{}
			argListT := []IRType{}
			sTarget, vType, vTarget := generateIRExpr(c, ir, e.Index)
			ret = append(ret, sTarget...)
			argList = append(argList, vTarget)
			argListT = append(argListT, vType)
			for _, arg := range e.Arguments {
				s, t, v := generateIRExpr(c, ir, arg)
				ret = append(ret, s...)
				argList = append(argList, v)
				argListT = append(argListT, t)
			}
			ret = append(ret, WrapIRExpression(&IRExpression{
				BaseASTNode:    e.BaseASTNode,
				TargetName:     resultVar,
				Operation:      "Call",
				Type:           translateType(e.ResolvedType),
				Arguments:      argList,
				ArgumentsTypes: argListT,
			}))
			return ret, translateType(e.ResolvedType), resultVar
		}
	} else if e, ok := node.(*ast.Unary); ok {
		if e.IsOperation() {
			s, t, v := generateIRExpr(c, ir, e.Unary)
			ret = append(ret, s...)
			ret = append(ret, WrapIRExpression(&IRExpression{
				BaseASTNode:    e.BaseASTNode,
				TargetName:     resultVar,
				Operation:      e.Op,
				Type:           translateType(e.ResolvedType),
				Arguments:      []string{v},
				ArgumentsTypes: []IRType{t},
			}))
			return ret, translateType(e.ResolvedType), resultVar
		} else if e.IsUnaryApplication() {
			return generateIRExpr(c, ir, e.UnaryApplication)
		}
	} else if e, ok := node.(*ast.Multiplication); ok {
		if e.HasNext() {
			ls, lt, lv := generateIRExpr(c, ir, e.Unary)
			rs, rt, rv := generateIRExpr(c, ir, e.Next)
			ret = append(ret, ls...)
			ret = append(ret, rs...)
			ret = append(ret, WrapIRExpression(&IRExpression{
				BaseASTNode:    e.BaseASTNode,
				TargetName:     resultVar,
				Operation:      e.Op,
				Type:           translateType(e.ResolvedType),
				Arguments:      []string{lv, rv},
				ArgumentsTypes: []IRType{lt, rt},
			}))
			return ret, translateType(e.ResolvedType), resultVar
		} else {
			return generateIRExpr(c, ir, e.Unary)
		}
	} else if e, ok := node.(*ast.Addition); ok {
		if e.HasNext() {
			ls, lt, lv := generateIRExpr(c, ir, e.Multiplication)
			rs, rt, rv := generateIRExpr(c, ir, e.Next)
			ret = append(ret, ls...)
			ret = append(ret, rs...)
			ret = append(ret, WrapIRExpression(&IRExpression{
				BaseASTNode:    e.BaseASTNode,
				TargetName:     resultVar,
				Operation:      e.Op,
				Type:           translateType(e.ResolvedType),
				Arguments:      []string{lv, rv},
				ArgumentsTypes: []IRType{lt, rt},
			}))
			return ret, translateType(e.ResolvedType), resultVar
		} else {
			return generateIRExpr(c, ir, e.Multiplication)
		}
	} else if e, ok := node.(*ast.Comparison); ok {
		if e.HasNext() {
			ls, lt, lv := generateIRExpr(c, ir, e.Addition)
			rs, rt, rv := generateIRExpr(c, ir, e.Next)
			ret = append(ret, ls...)
			ret = append(ret, rs...)
			ret = append(ret, WrapIRExpression(&IRExpression{
				BaseASTNode:    e.BaseASTNode,
				TargetName:     resultVar,
				Operation:      e.Op,
				Type:           translateType(e.ResolvedType),
				Arguments:      []string{lv, rv},
				ArgumentsTypes: []IRType{lt, rt},
			}))
			return ret, translateType(e.ResolvedType), resultVar
		} else {
			return generateIRExpr(c, ir, e.Addition)
		}
	} else if e, ok := node.(*ast.Equality); ok {
		if e.HasNext() {
			ls, lt, lv := generateIRExpr(c, ir, e.Comparison)
			rs, rt, rv := generateIRExpr(c, ir, e.Next)
			ret = append(ret, ls...)
			ret = append(ret, rs...)
			ret = append(ret, WrapIRExpression(&IRExpression{
				BaseASTNode:    e.BaseASTNode,
				TargetName:     resultVar,
				Operation:      e.Op,
				Type:           translateType(e.ResolvedType),
				Arguments:      []string{lv, rv},
				ArgumentsTypes: []IRType{lt, rt},
			}))
			return ret, translateType(e.ResolvedType), resultVar
		} else {
			return generateIRExpr(c, ir, e.Comparison)
		}
	} else if e, ok := node.(*ast.LogicalOperation); ok {
		if e.HasNext() {
			ls, lt, lv := generateIRExpr(c, ir, e.Equality)
			rs, rt, rv := generateIRExpr(c, ir, e.Next)
			ret = append(ret, ls...)
			ret = append(ret, rs...)
			ret = append(ret, WrapIRExpression(&IRExpression{
				BaseASTNode:    e.BaseASTNode,
				TargetName:     resultVar,
				Operation:      e.Op,
				Type:           translateType(e.ResolvedType),
				Arguments:      []string{lv, rv},
				ArgumentsTypes: []IRType{lt, rt},
			}))
			return ret, translateType(e.ResolvedType), resultVar
		} else {
			return generateIRExpr(c, ir, e.Equality)
		}
	} else if expr, ok := node.(*ast.Expression); ok {
		if expr.IsLogicalOperation() {
			return generateIRExpr(c, ir, expr.LogicalOperation)
		}
	}
	return ret, IR_UNKNOWN, resultVar
}

func genrateIR(graph *cfg.CFG, c *context.ParsingContext, ir *IRGeneratorState) {
	//func(cfg *CFG, block *Block) []generic_ast.NormalNode
	MapEntireGraph(graph, func(g *cfg.CFG, block *cfg.Block, node cfg.CFGCodeNode) []*IRStatement {
		fmt.Printf("=> NODE[%d]: %s\n", block.ID, reflect.TypeOf(node))
		ret := []*IRStatement{}

		if e, ok := node.(*ast.If); ok {
			s, t, v := generateIRExpr(c, ir, e.Condition)

			ret = append(ret, s...)
			ret = append(ret, WrapIRIf(&IRIf{
				BaseASTNode:   e.BaseASTNode,
				Condition:     v,
				ConditionType: t,
				BlockThen:     block.GetSuccs()[0],
				BlockElse:     block.GetSuccs()[1],
			}))
			return ret
		} else if expr, ok := node.(generic_ast.Expression); ok {
			if _, ok := (expr.(*ast.Empty)); ok {
				//ret = append(ret, WrapIREmpty())
				return ret
			}
			if assStmt, ok := (expr.(*ast.Assignment)); ok {
				exprIR, varType, varName := generateIRExpr(c, ir, assStmt.Value)
				ret = append(ret, exprIR...)
				ret = append(ret, WrapIRExpression(&IRExpression{
					BaseASTNode:    assStmt.BaseASTNode,
					TargetName:     assStmt.TargetName,
					Operation:      "Copy",
					Type:           varType,
					ArgumentsTypes: []IRType{varType},
					Arguments:      []string{varName},
				}))
			} else if declStmt, ok := (expr.(*ast.Declaration)); ok {
				for _, item := range declStmt.Items {
					if item.HasInitializer() {
						exprIR, varType, varName := generateIRExpr(c, ir, item.Initializer)
						ret = append(ret, exprIR...)
						ret = append(ret, WrapIRExpression(&IRExpression{
							BaseASTNode:    item.BaseASTNode,
							TargetName:     item.Name,
							Operation:      "Copy",
							Type:           varType,
							Arguments:      []string{varName},
							ArgumentsTypes: []IRType{varType},
						}))
					}
				}
				return ret
			}
			if retStmt, ok := (expr.(*ast.Return)); ok {
				if retStmt.HasExpression() {
					exprIR, varType, varName := generateIRExpr(c, ir, retStmt.Expression)
					ret = append(ret, exprIR...)
					ret = append(ret, WrapIRExit(&IRExit{
						BaseASTNode: retStmt.BaseASTNode,
						Value:       &varName,
						Type:        varType,
					}))
					return ret
				} else {
					ret = append(ret, WrapIRExit(&IRExit{
						BaseASTNode: retStmt.BaseASTNode,
						Value:       nil,
					}))
					return ret
				}
			}
			//ret = append(ret, WrapIREmpty())
			return ret
		}

		//ret = append(ret, WrapIREmpty())
		return ret
	})
}

func convertToSSA(graph *cfg.CFG, ir *IRGeneratorState) {
	visitedIDs := map[int]struct{}{}
	nameMapping := cfg.VariableSubstitutionMap{}
	allBlockOutputMappings := map[int]cfg.VariableSubstitutionMap{}
	allBlockInputMappings := map[int]cfg.VariableSubstitutionMap{}
	graph.VisitGraph(graph.Entry, func(g *cfg.CFG, block *cfg.Block, next func(blockID int)) {
		if _, wasVisited := visitedIDs[block.ID]; wasVisited {
			return
		}
		visitedIDs[block.ID] = struct{}{}
		code := graph.GetBlockCode(block.ID)
		if codeBlock, ok := code.(*IRBlock); ok {
			blockOutputsMapping := cfg.VariableSubstitutionMap{}
			allBlockInputMappings[block.ID] = nameMapping.Copy()
			for _, stmt := range codeBlock.Statements {
				vars := graph.ReferencedVars(stmt)
				newNameMapping := cfg.VariableSubstitutionMap{}
				for varName, _ := range vars.Assigned() {
					if !ir.isTempVar(varName) {
						newNameMapping[varName] = ir.NextVar(block.ID, varName)
					}
				}

				// for varName, _ := range vars.Used() {
				// 	if !nameMapping.Has(varName) && !ir.isTempVar(varName) {
				// 		nameMapping[varName] = ir.NextVar(varName)
				// 	}
				// }
				cfg.ReplaceVariables(stmt, nameMapping, newNameMapping, map[generic_ast.TraversableNode]struct{}{})
				nameMapping.Join(newNameMapping)
				blockOutputsMapping.Join(nameMapping)
			}
			allBlockOutputMappings[block.ID] = blockOutputsMapping
		}

		for _, stmt := range block.GetSuccs() {
			if _, wasVisited := visitedIDs[stmt]; wasVisited {
				continue
			}
			next(stmt)
		}
	})

	visitedIDs = map[int]struct{}{}
	graph.VisitGraph(graph.Entry, func(g *cfg.CFG, block *cfg.Block, next func(blockID int)) {
		if _, wasVisited := visitedIDs[block.ID]; wasVisited {
			return
		}
		visitedIDs[block.ID] = struct{}{}
		code := graph.GetBlockCode(block.ID)
		if codeBlock, ok := code.(*IRBlock); ok {
			blockPreds := block.GetPreds()
			if len(blockPreds) > 1 {
				//blockOutputMappings := allBlockOutputMappings[block.ID]
				blockInputMappings := allBlockInputMappings[block.ID]
				subst := cfg.VariableSubstitutionMap{}

				headers := []*IRStatement{}

				for origVarName, varName := range blockInputMappings {
					phiTarget := fmt.Sprintf("%s_phi_%d", origVarName, block.ID)
					phiBlocks := map[int]string{}
					phiValues := map[string]struct{}{}

					for _, pred := range blockPreds {
						predOutputs := allBlockOutputMappings[pred]
						if predVarName, ok := predOutputs[origVarName]; ok {
							phiBlocks[pred] = predVarName
							phiValues[predVarName] = struct{}{}
						}
					}

					if len(phiBlocks) > 0 && len(phiValues) > 1 {
						headers = append(headers, WrapIRPhi(CreateIRPhi(
							phiTarget,
							phiBlocks,
						)))
						subst[varName] = phiTarget
					}
				}

				if len(headers) > 0 {
					headers = append(headers, codeBlock.Statements...)
					codeBlock.Statements = headers
				}
				cfg.ReplaceVariables(code, subst, cfg.VariableSubstitutionMap{}, map[generic_ast.TraversableNode]struct{}{})
			}
		}
		for _, stmt := range block.GetSuccs() {
			if _, wasVisited := visitedIDs[stmt]; wasVisited {
				continue
			}
			next(stmt)
		}
	})
}

type ControlFlowGraphMapper func(cfg *cfg.CFG, block *cfg.Block, node cfg.CFGCodeNode) []*IRStatement

func MapEntireGraph(graph *cfg.CFG, mapper ControlFlowGraphMapper) {
	visitedIDs := map[int]struct{}{}

	blockContents := map[int]cfg.CFGCodeNode{}

	graph.VisitGraph(graph.Entry, func(cfg *cfg.CFG, block *cfg.Block, next func(blockID int)) {
		if _, wasVisited := visitedIDs[block.ID]; wasVisited {
			return
		}

		fmt.Printf("KURWA MAC PIERDOLONA W DUPE: %d\n", block.ID)

		visitedIDs[block.ID] = struct{}{}

		blockContents[block.ID] = &IRBlock{
			BlockID:    block.ID,
			Label:      "block",
			Statements: mapper(cfg, block, graph.GetBlockCode(block.ID)),
		}

		for _, stmt := range block.GetSuccs() {
			stmt := graph.ResolveID(stmt)
			if _, wasVisited := visitedIDs[stmt]; wasVisited {
				continue
			}
			next(stmt)
		}
	})

	for id, content := range blockContents {
		graph.OverrideBlockCode(id, content)
	}
}

func collapseToSimpleBlocks(graph *cfg.CFG) bool {
	visitedIDs := map[int]struct{}{}
	idsToRemove := map[int]struct{}{}
	mergedAnything := false
	visitor := func(block *cfg.Block) {
		if _, wasVisited := visitedIDs[block.ID]; wasVisited {
			return
		}
		blockPreds := block.GetPreds()
		blockSuccs := block.GetSuccs()
		if len(blockPreds) == 1 && len(blockPreds) == len(blockSuccs) && block.ID != graph.Entry && block.ID != graph.Exit {
			// Good candidate to merge
			sibling := graph.GetBlock(blockPreds[0])
			siblingPreds := sibling.GetPreds()
			siblingSuccs := sibling.GetSuccs()
			if len(siblingPreds) == 1 && len(siblingPreds) == len(siblingSuccs) && sibling.ID != graph.Entry && sibling.ID != graph.Exit {
				if predBlock, ok := graph.GetBlockCode(sibling.ID).(*IRBlock); ok {
					if curBlock, ok := graph.GetBlockCode(block.ID).(*IRBlock); ok {
						fmt.Printf("?> Merge %d into %d\n", block.ID, sibling.ID)
						mergedAnything = true
						predBlock.Join(curBlock)
						// rewire
						idPos := -1
						for index, id := range siblingSuccs {
							if id == block.ID {
								idPos = index
								break
							}
						}
						siblingSuccs[idPos] = blockSuccs[0]
						graph.ShadowBlock(block.ID, sibling)
						idsToRemove[block.ID] = struct{}{}
						return
					} else {
						fmt.Printf("!> (block type is %s) CANNOT Merge %d into %d\n", reflect.TypeOf(graph.GetBlockCode(block.ID)), block.ID, sibling.ID)
					}
				} else {
					fmt.Printf("!> (sibling type is %s) CANNOT Merge %d into %d\n", reflect.TypeOf(graph.GetBlockCode(sibling.ID)), block.ID, sibling.ID)
				}
			} else {
				fmt.Printf("!> (sibling pred/succ %d/%d) CANNOT Merge %d into %d\n", len(siblingPreds), len(siblingSuccs), block.ID, sibling.ID)
			}
		} else {
			fmt.Printf("!> (block pred/succ %d/%d) CANNOT Merge %d into ANY\n", len(blockPreds), len(blockSuccs), block.ID)
		}
	}

	graph.VisitGraph(graph.Entry, func(cfg *cfg.CFG, block *cfg.Block, next func(blockID int)) {
		if _, wasVisited := visitedIDs[block.ID]; wasVisited {
			return
		}
		visitedIDs[block.ID] = struct{}{}
		succs := block.GetSuccs()
		visitor(block)
		for _, stmt := range succs {
			if _, wasVisited := visitedIDs[graph.ResolveID(stmt)]; wasVisited {
				continue
			}
			visitor(graph.GetBlock(stmt))
			next(stmt)
		}
	})
	graph.RemoveBlocks(idsToRemove)
	return mergedAnything
}

func outputIR(root generic_ast.Expression, graph *cfg.CFG, c *context.ParsingContext) *IRFunction {
	rootNode := root.(*ast.TopDef).Function
	ret := &IRFunction{
		FunctionBody: []*IRBlock{},
		BaseASTNode:  rootNode.BaseASTNode,
		ReturnType:   translateType(rootNode.ResolvedType),
		Args:         []string{},
		ArgsTypes:    []IRType{},
		Name:         rootNode.Name,
	}
	visitedIDs := map[int]struct{}{}
	graph.VisitGraph(graph.Entry, func(g *cfg.CFG, block *cfg.Block, next func(blockID int)) {
		if _, wasVisited := visitedIDs[block.ID]; wasVisited {
			return
		}
		visitedIDs[block.ID] = struct{}{}
		b := graph.GetBlockCode(block.ID).(*IRBlock)
		if block.ID != graph.Entry && block.ID != graph.Exit {
			ret.FunctionBody = append(ret.FunctionBody, b)
		}
		for _, stmt := range block.GetSuccs() {
			if _, wasVisited := visitedIDs[graph.ResolveID(stmt)]; wasVisited {
				continue
			}
			next(stmt)
		}
	})

	return ret
}

func CreateIR(root generic_ast.Expression, flow cfg.FlowAnalysis, c *context.ParsingContext) *IRFunction {
	ir := &IRGeneratorState{
		tempCounter:      0,
		tempCounterBlock: map[int]int{},
		temps:            map[string]struct{}{},
	}
	graph := flow.Graph()
	genrateIR(graph, c, ir)
	for i := 0; i < 20; i++ {
		if !collapseToSimpleBlocks(graph) {
			break
		}
	}
	convertToSSA(graph, ir)

	return outputIR(root, graph, c)
}
