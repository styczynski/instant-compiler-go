package ir

import (
	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
)

func copyCollaps(graph *cfg.CFG, c *context.ParsingContext, subst cfg.VariableSubstitutionMap) (cfg.VariableSubstitutionMap, bool) {
	visitedIDs := map[int]struct{}{}
	newSubstFound := false
	substFrom := ""
	substTo := ""
	graph.VisitGraph(graph.Entry, func(g *cfg.CFG, block *cfg.Block, next func(blockID int)) {
		if _, wasVisited := visitedIDs[block.ID]; wasVisited {
			return
		}
		if newSubstFound {
			return
		}
		visitedIDs[block.ID] = struct{}{}
		b := graph.GetBlockCode(block.ID).(*IRBlock)
		//cfg.ReplaceVariables(b, subst, subst, map[generic_ast.TraversableNode]struct{}{})

		newStmts := []*IRStatement{}
		if len(b.Statements) > 0 {
			for i, stmt := range b.Statements {
				if stmt.IsCopy() && i != len(b.Statements)-1 {
					next := b.Statements[i+1]
					varSet := next.VarOut
					cpy := stmt.Copy
					vName := cpy.Var
					isOk := true
					// for {
					// 	if vNewName, ok := subst[vName]; ok {
					// 		fmt.Printf("   %s -> %s\n", vName, vNewName)
					// 		vName = vNewName
					// 	} else {
					// 		break
					// 	}
					// }
					if varSet.HasVariable(vName) {
						isOk = false
					}
					//fmt.Printf("Ignore %s because it has %s (%s)\n", cpy.Print(c), vName, next.Print(c))
					if isOk {
						if !newSubstFound {
							// We can eliminate variable
							//fmt.Printf("Subst %s => %s because of %s and next %s\n", cpy.TargetName, cpy.Var, stmt.Print(c), next.Print(c))
							subst[cpy.TargetName] = cpy.Var
							newSubstFound = true
							substFrom = cpy.TargetName
							substTo = cpy.Var
							continue
						}
					}
				}
				newStmts = append(newStmts, stmt)
			}
			//newStmts = append(newStmts, b.Statements[len(b.Statements)-1])
			b.Statements = newStmts
		}
		//cfg.ReplaceVariables(b, subst, subst, map[generic_ast.TraversableNode]struct{}{})

		for _, stmt := range block.GetSuccs() {
			if _, wasVisited := visitedIDs[graph.ResolveID(stmt)]; wasVisited {
				continue
			}
			next(stmt)
		}
	})

	if !newSubstFound {
		return subst, false
	}

	subst[substFrom] = substTo

	visitedIDs = map[int]struct{}{}
	graph.VisitGraph(graph.Entry, func(g *cfg.CFG, block *cfg.Block, next func(blockID int)) {
		if _, wasVisited := visitedIDs[block.ID]; wasVisited {
			return
		}
		visitedIDs[block.ID] = struct{}{}
		b := graph.GetBlockCode(block.ID).(*IRBlock)
		cfg.ReplaceVariables(b, subst, subst, map[generic_ast.TraversableNode]struct{}{})

		for _, stmt := range block.GetSuccs() {
			if _, wasVisited := visitedIDs[graph.ResolveID(stmt)]; wasVisited {
				continue
			}
			next(stmt)
		}

	})
	//fmt.Printf("EXIT subs\n")
	return subst, true
}
