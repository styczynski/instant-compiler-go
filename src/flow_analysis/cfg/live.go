package cfg

import (
	"github.com/willf/bitset"

	"github.com/styczynski/latte-compiler/src/generic_ast"
)

type NodeVariablesMapping map[generic_ast.NormalNode]VariableSet

type LiveVariablesInfo struct {
	in NodeVariablesMapping
	out NodeVariablesMapping
}

func (lvi LiveVariablesInfo) ReplaceBlock(old generic_ast.NormalNode, new generic_ast.NormalNode) {
	for block, vars := range lvi.in {
		key := block
		if block == old {
			delete(lvi.in, block)
			lvi.in[new] = vars
			key = new
		}
		lvi.in[key].ReplaceBlock(old, new)
	}
	for block, vars := range lvi.out {
		key := block
		if block == old {
			delete(lvi.out, block)
			lvi.out[new] = vars
			key = new
		}
		lvi.out[key].ReplaceBlock(old, new)
	}
}

func (lvi LiveVariablesInfo) BlockIn(block generic_ast.NormalNode) VariableSet {
	return lvi.in[block]
}

func (lvi LiveVariablesInfo) BlockOut(block generic_ast.NormalNode) VariableSet {
	return lvi.out[block]
}

func LiveVars(cfg *CFG) LiveVariablesInfo {
	vars, varsNames, def, use := createUsageBitsets(cfg)
	ins, outs := generateLivenessVariablesBits(cfg, def, use)
	return mapResultsSetsIntoVariables(cfg, vars, varsNames, ins, outs)
}

func generateLivenessVariablesBits(cfg *CFG, def, use map[generic_ast.NormalNode]*bitset.BitSet) (input, output map[generic_ast.NormalNode]*bitset.BitSet) {
	blocks := cfg.Blocks()
	input = make(map[generic_ast.NormalNode]*bitset.BitSet, len(blocks))
	output = make(map[generic_ast.NormalNode]*bitset.BitSet, len(blocks))
	for _, block := range blocks {
		input[block] = new(bitset.BitSet)
		output[block] = new(bitset.BitSet)
	}

	for {
		var change bool
		for _, block := range blocks {
			for _, s := range cfg.BlockSuccessors(block) {
				output[block].InPlaceUnion(input[s])
			}
			old := input[block].Clone()
			input[block] = use[block].Union(output[block].Difference(def[block]))
			change = change || !old.Equal(input[block])
		}
		if !change {
			break
		}
	}
	return input, output
}