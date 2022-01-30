package ir

import (
	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/parser/context"
)

func ifOrder(graph *cfg.CFG, c *context.ParsingContext) {
	visitedIDs := map[int]struct{}{}
	graph.VisitGraph(graph.Entry, func(g *cfg.CFG, block *cfg.Block, next func(blockID int)) {
		if _, wasVisited := visitedIDs[block.ID]; wasVisited {
			return
		}
		visitedIDs[block.ID] = struct{}{}
		b := graph.GetBlockCode(block.ID).(*IRBlock)

		if len(block.GetSuccs()) == 0 {
			return
		}

		firstNextBlock := block.GetSuccs()[0]

		newStmts := []*IRStatement{}
		for _, stmt := range b.Statements {
			if stmt.IsIf() {
				ifStmt := stmt.If
				if ifStmt.BlockThen == firstNextBlock && ifStmt.HasElseBlock() {
					ifStmt.BlockThen, ifStmt.BlockElse = ifStmt.BlockElse, -1
					ifStmt.Negated = true
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
