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
	fmt.Printf("Translate? %v\n", t)
	resolvedType := IR_UNKNOWN
	if _, ok := t.(*hindley_milner.FunctionType); ok {
		resolvedType = IR_FN
	} else if primitive, ok := t.(ast.PrimitiveType); ok {
		if primitive.Name() == "int" {
			resolvedType = IR_INT32
		} else if primitive.Name() == "boolean" {
			resolvedType = IR_BIT
		} else if primitive.Name() == "string" {
			resolvedType = IR_STRING
		} else if primitive.Name() == "void" {
			resolvedType = IR_VOID
		}
	}
	return resolvedType
}

func generateIRExpr(c *context.ParsingContext, ir *IRGeneratorState, node generic_ast.Expression) ([]*IRStatement, IRType, string) {
	fmt.Printf("  -> Expr: %s{%s}\n", reflect.TypeOf(node), node.(generic_ast.PrintableNode).Print(c))

	ret := []*IRStatement{}
	resultVar := ir.NextTempVar()

	if block, ok := node.(*ast.Block); ok {
		// Nested block
		for _, stmt := range block.Statements {
			code, _, _ := generateIRExpr(c, ir, stmt)
			ret = append(ret, code...)
		}
		return ret, IR_VOID, ""
	} else if syscallExpr, ok := node.(*ast.Syscall); ok {
		argList := []string{}
		argListT := []IRType{}
		fmt.Printf("SYSCALL TARGET %v\n", syscallExpr.Target)
		for _, arg := range syscallExpr.Arguments {
			s, t, v := generateIRExpr(c, ir, arg)
			ret = append(ret, s...)
			argList = append(argList, v)
			argListT = append(argListT, t)
		}
		ret = append(ret, WrapIRCall(&IRCall{
			BaseASTNode:    syscallExpr.BaseASTNode,
			TargetName:     resultVar,
			CallTarget:     syscallExpr.Target,
			CallTargetType: IR_FN,
			Type:           translateType(syscallExpr.ResolvedType),
			Arguments:      argList,
			ArgumentsTypes: argListT,
			IsBuiltin:      true,
		}))
		return ret, translateType(syscallExpr.ResolvedType), resultVar
	} else if e, ok := node.(*ast.Primary); ok {
		if e.IsVariable() {
			return []*IRStatement{}, translateType(e.ResolvedType), *e.Variable
		} else if e.IsInt() {
			ret = append(ret, WrapIRConst(&IRConst{
				BaseASTNode: e.BaseASTNode,
				TargetName:  resultVar,
				Type:        translateType(e.ResolvedType),
				Value:       *e.Int,
			}).SetComment("Const int %s", e.Print(c)))
			return ret, translateType(e.ResolvedType), resultVar
		} else if e.IsBool() {
			v := int64(0)
			if *e.Bool {
				v = 1
			}
			ret = append(ret, WrapIRConst(&IRConst{
				BaseASTNode: e.BaseASTNode,
				TargetName:  resultVar,
				Type:        translateType(e.ResolvedType),
				Value:       v,
			}).SetComment("Const boolean %s", e.Print(c)))
			return ret, translateType(e.ResolvedType), resultVar
		} else if e.IsString() {
			ret = append(ret, WrapIRConst(&IRConst{
				BaseASTNode: e.BaseASTNode,
				TargetName:  resultVar,
				Type:        translateType(e.ResolvedType),
				StringValue: e.String,
			}).SetComment("Const string %s", *e.String))
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
			fmt.Printf("CALL TARGET %v\n", vTarget)
			ret = append(ret, sTarget...)
			for _, arg := range e.Arguments {
				s, t, v := generateIRExpr(c, ir, arg)
				ret = append(ret, s...)
				argList = append(argList, v)
				argListT = append(argListT, t)
			}
			ret = append(ret, WrapIRCall(&IRCall{
				BaseASTNode:    e.BaseASTNode,
				TargetName:     resultVar,
				CallTarget:     vTarget,
				CallTargetType: vType,
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
				Operation:      CreateIROperator(e.Op, 1, IR_OP_KIND_ANY),
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
				Operation:      CreateIROperator(e.Op, 2, IR_OP_KIND_NUMERIC),
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
			if lt == IR_STRING && rt == IR_STRING {
				ret = append(ret, WrapIRCall(&IRCall{
					BaseASTNode:    e.BaseASTNode,
					TargetName:     resultVar,
					CallTarget:     "AddStrings",
					CallTargetType: IR_FN,
					Type:           IR_STRING,
					Arguments:      []string{lv, rv},
					ArgumentsTypes: []IRType{lt, rt},
				}))
			} else {
				ret = append(ret, WrapIRExpression(&IRExpression{
					BaseASTNode:    e.BaseASTNode,
					TargetName:     resultVar,
					Operation:      CreateIROperator(e.Op, 2, IR_OP_KIND_NUMERIC),
					Type:           translateType(e.ResolvedType),
					Arguments:      []string{lv, rv},
					ArgumentsTypes: []IRType{lt, rt},
				}))
			}
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
				Operation:      CreateIROperator(e.Op, 2, IR_OP_KIND_LOGIC),
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
				Operation:      CreateIROperator(e.Op, 2, IR_OP_KIND_LOGIC),
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
				Operation:      CreateIROperator(e.Op, 2, IR_OP_KIND_LOGIC),
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
		} else if expr.IsSyscall() {
			return generateIRExpr(c, ir, expr.Syscall)
		}
	}
	return ret, IR_UNKNOWN, resultVar
}

func extractIfBlockJumpID(ifBlock *ast.Statement, graph *cfg.CFG) int {
	ifBlockChildren := ifBlock.GetChildren()
	if len(ifBlockChildren) > 0 {
		child := ifBlockChildren[0]
		if block, ok := child.(*ast.Block); ok {
			if len(block.Statements) > 0 {
				blockTopChildren := block.Statements[0].GetChildren()
				if len(blockTopChildren) > 0 {
					if node, ok := blockTopChildren[0].(cfg.CFGCodeNode); ok {
						jumpBlock := graph.FindBlock(node)
						if jumpBlock != nil {
							return jumpBlock.ID
						}
					}
				}
			}
		} else {
			// Single statement
			jumpBlock := graph.FindBlock(child.(cfg.CFGCodeNode))
			if jumpBlock != nil {
				return jumpBlock.ID
			}
		}
	}
	return -1
}

func genrateIR(graph *cfg.CFG, c *context.ParsingContext, ir *IRGeneratorState) {

	blocksLoopBackJumps := map[int]int{}

	//func(cfg *CFG, block *Block) []generic_ast.NormalNode
	MapEntireGraph(graph, func(g *cfg.CFG, block *cfg.Block, node cfg.CFGCodeNode) []*IRStatement {
		if node != nil {
			fmt.Printf("=> NODE[%d]: %s{%s}\n", block.ID, reflect.TypeOf(node), node.Print(c))
		}
		ret := []*IRStatement{}

		if bl, ok := node.(*ast.Block); ok {
			exprIR, _, _ := generateIRExpr(c, ir, bl)
			ret = append(ret, exprIR...)
			return ret
		} else if e, ok := node.(*ast.While); ok {
			sCond, tCond, vCond := generateIRExpr(c, ir, e.Condition)
			//doNode := e.Do.GetChildren()[0].(*ast.Block).Statements[0].GetChildren()[0].(cfg.CFGCodeNode)
			//doBlockID := extractIfBlockJumpID(e.Do, graph)
			ret = append(ret, sCond...)

			doBlockID := extractIfBlockJumpID(e.Do, graph)
			skipBlockID := -1
			for _, succ := range block.GetSuccs() {
				if succ != doBlockID {
					skipBlockID = succ
				}
			}

			traceBlockID := doBlockID
			lastDoBlockID := doBlockID
			started := false
			for traceBlockID != block.ID || !started {
				traceBlock := g.GetBlock(traceBlockID)
				traceNext := traceBlock.GetSuccs()
				if len(traceNext) != 1 {
					panic("While failure")
				}
				traceBlockID, lastDoBlockID = traceNext[0], traceBlock.ID
				started = true
				fmt.Printf("--> While points to [%s]\n", g.GetBlockCode(traceBlockID).Print(c))
			}

			fmt.Printf("WHILE LOOP\n")
			ret = append(ret, WrapIRIf(&IRIf{
				BaseASTNode:   e.BaseASTNode,
				Condition:     vCond,
				ConditionType: tCond,
				BlockThen:     skipBlockID,
				BlockElse:     -1,
				Negated:       true,
			}).SetComment("While condition"))

			// Append loop instruction
			blocksLoopBackJumps[lastDoBlockID] = block.ID

			return ret
		} else if e, ok := node.(*ast.If); ok {
			s, t, v := generateIRExpr(c, ir, e.Condition)
			// thenNode := e.Then.GetChildren()[0].(*ast.Block).Statements[0].GetChildren()[0].(cfg.CFGCodeNode)
			// fmt.Printf("[?] If contents: %s{%s}\n", reflect.TypeOf(thenNode), thenNode.Print(c))

			thenBlockID := extractIfBlockJumpID(e.Then, graph)
			elseBlockID := -1
			if e.HasElseBlock() {
				elseBlockID = extractIfBlockJumpID(e.Else, graph)
			} else {
				for _, succ := range block.GetSuccs() {
					if succ != thenBlockID {
						elseBlockID = succ
					}
				}
			}

			ret = append(ret, s...)
			ret = append(ret, WrapIRIf(&IRIf{
				BaseASTNode:   e.BaseASTNode,
				Condition:     v,
				ConditionType: t,
				BlockThen:     thenBlockID,
				BlockElse:     elseBlockID,
			}).SetComment("If condition"))
			return ret
		} else if expr, ok := node.(generic_ast.Expression); ok {
			if exprStmt, ok := (expr.(*ast.Expression)); ok {
				exprIR, _, _ := generateIRExpr(c, ir, exprStmt)
				ret = append(ret, exprIR...)
				return ret
			} else if _, ok := (expr.(*ast.Empty)); ok {
				//ret = append(ret, WrapIREmpty())
				return ret
			} else if assStmt, ok := (expr.(*ast.Assignment)); ok {
				exprIR, varType, varName := generateIRExpr(c, ir, assStmt.Value)
				ret = append(ret, exprIR...)
				ret = append(ret, WrapIRCopy(&IRCopy{
					BaseASTNode: assStmt.BaseASTNode,
					TargetName:  assStmt.TargetName,
					Type:        varType,
					Var:         varName,
				}).SetComment("Assign variable %s", assStmt.TargetName))
				return ret
			} else if declStmt, ok := (expr.(*ast.Declaration)); ok {
				for _, item := range declStmt.Items {
					if item.HasInitializer() {
						exprIR, varType, varName := generateIRExpr(c, ir, item.Initializer)
						ret = append(ret, exprIR...)
						ret = append(ret, WrapIRCopy(&IRCopy{
							BaseASTNode: item.BaseASTNode,
							TargetName:  item.Name,
							Type:        varType,
							Var:         varName,
						}).SetComment("Assign variable %s", item.Name))
					}
				}
				return ret
			} else if retStmt, ok := (expr.(*ast.Return)); ok {
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

	for blockID, targetBlockID := range blocksLoopBackJumps {
		bl := graph.GetBlockCode(blockID).(*IRBlock)
		bl.Statements = append(bl.Statements, WrapIRJump(&IRJump{
			BlockTarget: targetBlockID,
		}).SetComment("While loop return to block_%d", targetBlockID))
	}
}

type ControlFlowGraphMapper func(cfg *cfg.CFG, block *cfg.Block, node cfg.CFGCodeNode) []*IRStatement

func MapEntireGraph(graph *cfg.CFG, mapper ControlFlowGraphMapper) {
	visitedIDs := map[int]struct{}{}

	blockContents := map[int]cfg.CFGCodeNode{}

	graph.VisitGraph(graph.Entry, func(cfg *cfg.CFG, block *cfg.Block, next func(blockID int)) {
		if _, wasVisited := visitedIDs[block.ID]; wasVisited {
			return
		}

		//fmt.Printf("KURWA MAC PIERDOLONA W DUPE: %d\n", block.ID)

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

	// outputCodeIR := outputIR(root, graph, c)
	// fmt.Printf("\n\nSSA:\n\n%s", outputCodeIR.Print(c))

	phiElim(graph, c)
	regSplit(graph, c)
	ifOrder(graph, c)

	outputCodeIR := outputIR(root, graph, c)
	irAnalysis := cfg.CreateFlowAnalysis(outputCodeIR)

	irCFG := irAnalysis.Graph()
	irAnalysis.Optimize(c)

	irLiveness := irAnalysis.Liveness()
	reaching := irAnalysis.Reaching()
	for _, blockID := range irCFG.ListBlockIDs() {
		//fmt.Printf("CFG TYPE %s\n", reflect.TypeOf(irCFG.GetBlockCode(blockID)))
		parent := irCFG.GetBlockCode(blockID).Parent()
		if parent != nil {
			stmt := parent.(*IRStatement)
			stmt.SetFlowAnalysisProps(
				irLiveness.BlockIn(blockID),
				irLiveness.BlockOut(blockID),
				reaching.ReachedBlocks(blockID),
			)
		}
	}

	subst := cfg.VariableSubstitutionMap{}
	cont := true
	for cont {
		subst, cont = copyCollaps(graph, c, subst)
	}

	// fmt.Printf("Fold done (IR):\n")
	// fmt.Printf("\n\nENTIRE CODE IR:\n\n%s", outputCodeIR.Print(c))
	// fmt.Printf("\n\nENTIRE GRAPH IR:\n\n")
	// fmt.Print(irAnalysis.Print(c))
	// fmt.Printf("Yeah.\n")

	return outputCodeIR
}
