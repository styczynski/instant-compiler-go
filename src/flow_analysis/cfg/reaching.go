package cfg

import (
	"fmt"
	"strings"

	"github.com/willf/bitset"

	"github.com/styczynski/latte-compiler/src/generic_ast"
)

type ReachedBlocks map[int]struct{}

type ReachingVariablesInfo struct {
	usagesMap map[int]ReachedBlocks
}

func (lvi ReachingVariablesInfo) ReachedBlocks(blockID int) ReachedBlocks {
	return lvi.usagesMap[blockID]
}

func (rb ReachedBlocks) Print(cfg *CFG) string {
	blockStrs := []string{}
	for block, _ := range rb {
		blockStrs = append(blockStrs, fmt.Sprintf("%d", block))
	}
	return fmt.Sprintf("{%s}", strings.Join(blockStrs, ", "))
}

func DefUse(cfg *CFG) ReachingVariablesInfo {
	blocks, gen, kill := generateVariableInfoKillBitsets(cfg)
	ins, _ := defineVariableReachingBitsets(cfg, gen, kill)
	return ReachingVariablesInfo{
		usagesMap: defUseResultSet(blocks, ins),
	}
}

func DefsReaching(blockID int, cfg *CFG) map[int]struct{} {
	blocks, gen, kill := generateVariableInfoKillBitsets(cfg)
	ins, _ := defineVariableReachingBitsets(cfg, gen, kill)
	return defineReachingVariablesBlocksBitsets(blockID, blocks, ins)
}

func generateVariableInfoKillBitsets(cfg *CFG) (blocks []*Block, gen, kill map[int]*bitset.BitSet) {
	okills := make(map[Variable]*bitset.BitSet)
	gen = make(map[int]*bitset.BitSet)
	kill = make(map[int]*bitset.BitSet)
	blocks = cfg.Blocks()

	for _, b := range blocks { // prime
		gen[b.ID] = new(bitset.BitSet)
		kill[b.ID] = new(bitset.BitSet)
	}

	for i := 0; i < 2; i++ {
		for j, block := range blocks {
			j := uint(j)

			def := GetAllDeclaredVariables(cfg.codeMapping[block.ID], map[generic_ast.TraversableNode]struct{}{})

			for _, d := range def {
				if _, ok := okills[d]; !ok {
					okills[d] = new(bitset.BitSet)
				}
				gen[block.ID].Set(j)
				okills[d].Set(j)
				kill[block.ID] = kill[block.ID].Union(okills[d]).Difference(gen[block.ID])
			}
		}
	}
	return blocks, gen, kill
}

func defineVariableReachingBitsets(cfg *CFG, gen, kill map[int]*bitset.BitSet) (input, output map[int]*bitset.BitSet) {
	input = make(map[int]*bitset.BitSet)
	output = make(map[int]*bitset.BitSet)
	blocks := cfg.Blocks()
	for i := 0; i < len(blocks); i++ {
		block := blocks[i]
		input[block.ID] = new(bitset.BitSet)
		output[block.ID] = new(bitset.BitSet)
		if block.ID == cfg.Entry {
			blocks = append(blocks[:i], blocks[i+1:]...)
			i--
		}
	}
	for {
		var anythingChanged bool
		for _, block := range blocks {
			for _, p := range cfg.BlockPredecessors(block.ID) {
				input[block.ID].InPlaceUnion(output[p.ID])
			}
			old := output[block.ID].Clone()
			output[block.ID] = gen[block.ID].Union(input[block.ID].Difference(kill[block.ID]))
			anythingChanged = anythingChanged || !old.Equal(output[block.ID])
		}
		if !anythingChanged {
			break
		}
	}
	return input, output
}

func defUseResultSet(blocks []*Block, ins map[int]*bitset.BitSet) map[int]ReachedBlocks {
	du := make(map[int]ReachedBlocks)
	for _, block := range blocks {
		du[block.ID] = make(map[int]struct{})
	}
	for _, block := range blocks {
		for i, ok := uint(0), true; ok; i++ {
			if i, ok = ins[block.ID].NextSet(i); ok {
				du[blocks[i].ID][block.ID] = struct{}{}
			}
		}
	}
	return du
}

func defineReachingVariablesBlocksBitsets(blockID int, blocks []*Block, ins map[int]*bitset.BitSet) map[int]struct{} {
	result := make(map[int]struct{})
	insStmt, found := ins[blockID]
	if !found {
		panic("Node is missing from CFG.")
	}
	for i, ok := uint(0), true; ok; i++ {
		if i, ok = insStmt.NextSet(i); ok {
			result[blocks[i].ID] = struct{}{}
		}
	}
	return result
}
