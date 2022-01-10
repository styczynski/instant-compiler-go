package cfg

import (
	"fmt"

	"github.com/styczynski/latte-compiler/src/parser/context"
)

func (flow *FlowAnalysisImpl) Optimize(c *context.ParsingContext) {
	cfg := flow.Graph()
	liveness := flow.Liveness()

	for _, b := range cfg.blocks {
		block := b
		if cfg.codeMapping[block.ID] != nil {
			//fmt.Printf("OPTIMIZE %v\n", reflect.TypeOf(block.stmt))
			if rmb, ok := cfg.codeMapping[block.ID].(NodeWithRemovableVariableAsignment); ok {
				liveVars := liveness.BlockIn(block.ID)
				refVars := flow.graph.ReferencedVars(cfg.codeMapping[block.ID])
				varsToRemove := map[string]struct{}{}
				for defVarName, _ := range refVars.decl {
					if !liveVars.HasVariable(defVarName) {
						varsToRemove[defVarName] = struct{}{}
					}
				}
				for asgtVarName, _ := range refVars.asgt {
					if !liveVars.HasVariable(asgtVarName) {
						varsToRemove[asgtVarName] = struct{}{}
					}
				}
				fmt.Printf("REMOVE %v from block %d (live vars %v)\n", varsToRemove, block.ID, liveVars)
				cfg.ReplaceBlockCode(block.ID, rmb.RemoveVariableAssignment(varsToRemove))
			}

			//items = append(items, cfg.ReferencedVars(srcNode).Print())
			//items = append(items, liveness.BlockIn(srcBlock.stmt).String())
		}
	}
}
