package ir

import (
	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
)

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
						
					}
				} else {
					
				}
			} else {
				
			}
		} else {
			
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
