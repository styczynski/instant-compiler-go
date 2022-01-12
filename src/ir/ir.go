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
			ret = append(ret, WrapIRConst(&IRConst{
				BaseASTNode: e.BaseASTNode,
				TargetName:  resultVar,
				Type:        translateType(e.ResolvedType),
				Value:       e.Print(c),
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
				ret = append(ret, WrapIRCopy(&IRCopy{
					BaseASTNode: assStmt.BaseASTNode,
					TargetName:  assStmt.TargetName,
					Type:        varType,
					Var:         varName,
				}))
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
	phiElim(graph, c)
	regSplit(graph, c)

	outputCodeIR := outputIR(root, graph, c)

	irAnalysis := cfg.CreateFlowAnalysis(outputCodeIR)

	irCFG := irAnalysis.Graph()
	irAnalysis.Optimize(c)

	irLiveness := irAnalysis.Liveness()
	reaching := irAnalysis.Reaching()
	for _, blockID := range irCFG.ListBlockIDs() {
		fmt.Printf("CFG TYPE %s\n", reflect.TypeOf(irCFG.GetBlockCode(blockID)))
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

	fmt.Printf("Fold done (IR):\n")
	fmt.Printf("\n\nENTIRE CODE IR:\n\n%s", outputCodeIR.Print(c))
	fmt.Printf("\n\nENTIRE GRAPH IR:\n\n")
	fmt.Print(irAnalysis.Print(c))
	fmt.Printf("Yeah.\n")

	return outputCodeIR
}
