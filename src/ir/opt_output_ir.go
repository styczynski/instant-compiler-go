package ir

import (
	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
)

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

		// Set parents
		b.OverrideParent(ret)
		for _, stmt := range b.Statements {
			stmt.OverrideParent(b)
		}
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
