package cfg

import (
	"fmt"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
)

func (cfg *CFG) PrintCFG(c *context.ParsingContext) string {
	blocks := []context.SelectionBlock{}
	i := 0

	var describer func(src context.SelectionBlock, id int, mappingID func(selectionBlock context.SelectionBlock) int) []string
	describer = func(src context.SelectionBlock, id int, mappingID func(selectionBlock context.SelectionBlock) int) []string {
		items := []string{}
		srcNode := src.(generic_ast.NormalNodeSelection).GetNode()
		var srcBlock *block = nil
		for _, b := range cfg.blocks {
			if b.stmt == srcNode {
				srcBlock = b
				break
			}
		}
		for _, pred := range srcBlock.preds {
			if pred != nil {
				predID := mappingID(generic_ast.NewNormalNodeSelection(pred, describer))
				items = append(items, fmt.Sprintf("<-%d", predID))
			}
		}
		for _, succ := range srcBlock.succs {
			if succ != nil {
				predID := mappingID(generic_ast.NewNormalNodeSelection(succ, describer))
				items = append(items, fmt.Sprintf("->%d", predID))
			}
		}
		return items
	}

	for _, b := range cfg.blocks {
		block := b
		if block.stmt != nil {
			blocks = append(blocks, generic_ast.NewNormalNodeSelection(block.stmt, describer))
			i++
		}
	}

	return fmt.Sprintf("%s%s",
		c.PrintSelectionBlocksList(context.SelectionBlocks(blocks)),
		c.PrintSelectionBlocks(context.SelectionBlocks(blocks)),
		)
}
