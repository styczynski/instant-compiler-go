package cfg

import (
	"fmt"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
)

func (flow *FlowAnalysisImpl) Print(c *context.ParsingContext) string {
	cfg := flow.Graph()
	liveness := flow.Liveness()
	reaching := flow.Reaching()

	blocks := []context.SelectionBlock{}
	i := 0

	var describer func(src context.SelectionBlock, id int, mappingID func(selectionBlock context.SelectionBlock) int) []string
	describer = func(src context.SelectionBlock, id int, mappingID func(selectionBlock context.SelectionBlock) int) []string {
		items := []string{}
		srcNode := src.(generic_ast.NormalNodeSelection).GetNode().(interface{})
		var srcBlock *Block = nil
		for _, b := range cfg.blocks {
			if cfg.codeMapping[b.ID] == srcNode {
				srcBlock = b
				break
			}
		}
		for _, pred := range srcBlock.preds {
			//predID := mappingID(generic_ast.NewNormalNodeSelection(cfg.codeMapping[pred], -1, describer))
			items = append(items, fmt.Sprintf("<-%d", pred))
		}
		for _, succ := range srcBlock.succs {
			//predID := mappingID(generic_ast.NewNormalNodeSelection(cfg.codeMapping[succ], -1, describer))
			items = append(items, fmt.Sprintf("->%d", succ))
		}
		items = append(items, cfg.ReferencedVars(srcNode.(CFGCodeNode)).Print())
		items = append(items, liveness.BlockIn(srcBlock.ID).String())
		items = append(items, reaching.ReachedBlocks(src.GetID()).Print(cfg))
		return items
	}

	for _, b := range cfg.blocks {
		block := b
		if cfg.codeMapping[block.ID] != nil {
			blocks = append(blocks, generic_ast.NewNormalNodeSelection(cfg.codeMapping[block.ID], block.ID, describer))
			i++
		}
	}

	return fmt.Sprintf("%s%s",
		c.PrintSelectionBlocksList(context.SelectionBlocks(blocks)),
		c.PrintSelectionBlocks(context.SelectionBlocks(blocks)),
	)
}
