package ir

import (
	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/parser/context"
)

var REG_SELF_OP_MAPPING = map[IROperator]IROperator{
	IR_OP_ADD: IR_OP_SELF_ADD,
	IR_OP_SUB: IR_OP_SELF_SUB,
	IR_OP_MUL: IR_OP_SELF_MUL,
	IR_OP_DIV: IR_OP_SELF_DIV,
	// IR_OP_EQ: IR_OP_EQ,
	// IR_OP_GT:  IR_OP_GT,
	// IR_OP_LT:  ,
	// IR_OP_GTEQ: "gte",
	// IR_OP_LTEQ: "lte",
}

func regSplit(graph *cfg.CFG, c *context.ParsingContext) {
	visitedIDs := map[int]struct{}{}
	graph.VisitGraph(graph.Entry, func(g *cfg.CFG, block *cfg.Block, next func(blockID int)) {
		if _, wasVisited := visitedIDs[block.ID]; wasVisited {
			return
		}
		visitedIDs[block.ID] = struct{}{}
		b := graph.GetBlockCode(block.ID).(*IRBlock)

		newStmts := []*IRStatement{}
		for _, stmt := range b.Statements {
			if stmt.IsExpression() {
				expr := stmt.Expression

				if len(expr.Arguments) == 2 {
					hasOp := false
					for op, selfOp := range REG_SELF_OP_MAPPING {
						if expr.Operation == op {
							newStmts = append(newStmts, WrapIRCopy(&IRCopy{
								BaseASTNode: expr.BaseASTNode,
								Type:        expr.ArgumentsTypes[0],
								TargetName:  expr.TargetName,
								Var:         expr.Arguments[0],
							}))
							newStmts = append(newStmts, WrapIRExpression(&IRExpression{
								BaseASTNode:    expr.BaseASTNode,
								Type:           expr.Type,
								TargetName:     expr.TargetName,
								Operation:      selfOp,
								ArgumentsTypes: []IRType{expr.ArgumentsTypes[1]},
								Arguments:      []string{expr.Arguments[1]},
							}))
							hasOp = true
							break
						}
					}
					if hasOp {
						continue
					}
				}
			}
			newStmts = append(newStmts, stmt)
		}
		b.Statements = newStmts

		for _, stmt := range block.GetSuccs() {
			if _, wasVisited := visitedIDs[graph.ResolveID(stmt)]; wasVisited {
				continue
			}
			next(stmt)
		}
	})
}
