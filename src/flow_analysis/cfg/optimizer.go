package cfg

import (
	"github.com/styczynski/latte-compiler/src/parser/context"
)

func (flow *FlowAnalysisImpl) Optimize(c *context.ParsingContext) {
	cfg := flow.Graph()
	liveness := flow.Liveness()

	for _, b := range cfg.blocks {
		block := b
		if block.stmt != nil {
			//fmt.Printf("OPTIMIZE %v\n", reflect.TypeOf(block.stmt))
			if rmb, ok := block.stmt.(NodeWithRemovableVariableAsignment); ok {
				liveVars := liveness.BlockIn(block.stmt)
				refVars := flow.graph.ReferencedVars(block.stmt)
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
				cfg.ReplaceBlock(block.stmt, rmb.RemoveVariableAssignment(varsToRemove))
			}

			//items = append(items, cfg.ReferencedVars(srcNode).Print())
			//items = append(items, liveness.BlockIn(srcBlock.stmt).String())
		}
	}
}
