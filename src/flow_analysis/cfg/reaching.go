package cfg

import (
	"fmt"
	"strings"

	"github.com/willf/bitset"

	"github.com/styczynski/latte-compiler/src/generic_ast"
)

type ReachedBlocks map[generic_ast.NormalNode]struct{}

type ReachingVariablesInfo struct {
	usagesMap map[generic_ast.NormalNode]ReachedBlocks
}

func (lvi ReachingVariablesInfo) ReplaceBlock(old generic_ast.NormalNode, new generic_ast.NormalNode) {
	for block, usages := range lvi.usagesMap {
		key := block
		if block == old {
			delete(lvi.usagesMap, block)
			lvi.usagesMap[new] = usages
			key = new
		}
		for subBlock, _ := range lvi.usagesMap[key] {
			if subBlock == old {
				delete(lvi.usagesMap[key], old)
				lvi.usagesMap[key][new] = struct{}{}
			}
		}
	}
}


func (lvi ReachingVariablesInfo) ReachedBlocks(block generic_ast.NormalNode) ReachedBlocks {
	return lvi.usagesMap[block]
}

func (rb ReachedBlocks) Print(cfg *CFG) string {
	blockStrs := []string{}
	for block, _ := range rb {
		blockStrs = append(blockStrs, fmt.Sprintf("%d", cfg.blocksIDs[block]))
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

func DefsReaching(stmt generic_ast.NormalNode, cfg *CFG) map[generic_ast.NormalNode]struct{} {
	blocks, gen, kill := generateVariableInfoKillBitsets(cfg)
	ins, _ := defineVariableReachingBitsets(cfg, gen, kill)
	return defineReachingVariablesBlocksBitsets(stmt, blocks, ins)
}

func generateVariableInfoKillBitsets(cfg *CFG) (blocks []generic_ast.NormalNode, gen, kill map[generic_ast.NormalNode]*bitset.BitSet) {
	okills := make(map[Variable]*bitset.BitSet)
	gen = make(map[generic_ast.NormalNode]*bitset.BitSet)
	kill = make(map[generic_ast.NormalNode]*bitset.BitSet)
	blocks = cfg.Blocks()

	for _, b := range blocks { // prime
		gen[b] = new(bitset.BitSet)
		kill[b] = new(bitset.BitSet)
	}

	for i := 0; i < 2; i++ {
		for j, block := range blocks {
			j := uint(j)

			def := GetAllDeclaredVariables(block, map[generic_ast.TraversableNode]struct{}{})

			for _, d := range def {
				if _, ok := okills[d]; !ok {
					okills[d] = new(bitset.BitSet)
				}
				gen[block].Set(j)
				okills[d].Set(j)
				kill[block] = kill[block].Union(okills[d]).Difference(gen[block])
			}
		}
	}
	return blocks, gen, kill
}

func defineVariableReachingBitsets(cfg *CFG, gen, kill map[generic_ast.NormalNode]*bitset.BitSet) (input, output map[generic_ast.NormalNode]*bitset.BitSet) {
	input = make(map[generic_ast.NormalNode]*bitset.BitSet)
	output = make(map[generic_ast.NormalNode]*bitset.BitSet)
	blocks := cfg.Blocks()
	for i := 0; i < len(blocks); i++ {
		block := blocks[i]
		input[block] = new(bitset.BitSet)
		output[block] = new(bitset.BitSet)
		if block == cfg.Entry {
			blocks = append(blocks[:i], blocks[i+1:]...)
			i--
		}
	}
	for {
		var anythingChanged bool
		for _, block := range blocks {
			for _, p := range cfg.BlockPredecessors(block) {
				input[block].InPlaceUnion(output[p])
			}
			old := output[block].Clone()
			output[block] = gen[block].Union(input[block].Difference(kill[block]))
			anythingChanged = anythingChanged || !old.Equal(output[block])
		}
		if !anythingChanged {
			break
		}
	}
	return input, output
}

func defUseResultSet(blocks []generic_ast.NormalNode, ins map[generic_ast.NormalNode]*bitset.BitSet) map[generic_ast.NormalNode]ReachedBlocks {
	du := make(map[generic_ast.NormalNode]ReachedBlocks)
	for _, block := range blocks {
		du[block] = make(map[generic_ast.NormalNode]struct{})
	}
	for _, block := range blocks {
		for i, ok := uint(0), true; ok; i++ {
			if i, ok = ins[block].NextSet(i); ok {
				du[blocks[i]][block] = struct{}{}
			}
		}
	}
	return du
}

func defineReachingVariablesBlocksBitsets(stmt generic_ast.NormalNode, blocks []generic_ast.NormalNode, ins map[generic_ast.NormalNode]*bitset.BitSet) map[generic_ast.NormalNode]struct{} {
	result := make(map[generic_ast.NormalNode]struct{})
	insStmt, found := ins[stmt]
	if !found {
		panic("Node is missing from CFG.")
	}
	for i, ok := uint(0), true; ok; i++ {
		if i, ok = insStmt.NextSet(i); ok {
			result[blocks[i]] = struct{}{}
		}
	}
	return result
}