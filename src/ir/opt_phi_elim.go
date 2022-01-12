package ir

import (
	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/parser/context"
)

func phiElim(graph *cfg.CFG, c *context.ParsingContext) {
	visitedIDs := map[int]struct{}{}
	graph.VisitGraph(graph.Entry, func(g *cfg.CFG, block *cfg.Block, next func(blockID int)) {
		if _, wasVisited := visitedIDs[block.ID]; wasVisited {
			return
		}
		visitedIDs[block.ID] = struct{}{}
		b := graph.GetBlockCode(block.ID).(*IRBlock)
		for i, stmt := range b.Statements {
			if stmt.IsPhi() {
				phi := stmt.Phi
				for i, pred := range phi.Blocks {
					predCode := graph.GetBlockCode(pred).(*IRBlock)
					predCode.Statements = append(predCode.Statements, WrapIRCopy(&IRCopy{
						BaseASTNode: phi.BaseASTNode,
						TargetName:  phi.TargetName,
						Type:        phi.Type,
						Var:         phi.Values[i],
					}))
				}
				b.Statements[i] = WrapIREmpty()
			}
		}

		for _, stmt := range block.GetSuccs() {
			if _, wasVisited := visitedIDs[graph.ResolveID(stmt)]; wasVisited {
				continue
			}
			next(stmt)
		}
	})
}
