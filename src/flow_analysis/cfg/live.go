package cfg

import (
	"github.com/willf/bitset"
)

type NodeVariablesMapping map[int]VariableSet

type LiveVariablesInfo struct {
	in  NodeVariablesMapping
	out NodeVariablesMapping
}

func (lvi LiveVariablesInfo) BlockIn(blockID int) VariableSet {
	return lvi.in[blockID]
}

func (lvi LiveVariablesInfo) BlockOut(blockID int) VariableSet {
	return lvi.out[blockID]
}

func LiveVars(cfg *CFG) LiveVariablesInfo {
	vars, varsNames, def, use := createUsageBitsets(cfg)
	ins, outs := generateLivenessVariablesBits(cfg, def, use)
	return mapResultsSetsIntoVariables(cfg, vars, varsNames, ins, outs)
}


func generateLivenessVariablesBits(cfg *CFG, def, use map[int]*bitset.BitSet) (input, output map[int]*bitset.BitSet) {
	blocks := cfg.Blocks()
	input = make(map[int]*bitset.BitSet, len(blocks))
	output = make(map[int]*bitset.BitSet, len(blocks))
	for _, block := range blocks {
		input[block.ID] = new(bitset.BitSet)
		output[block.ID] = new(bitset.BitSet)
	}

	for {
		var change bool
		for _, block := range blocks {
			for _, s := range block.succs {
				output[block.ID].InPlaceUnion(input[s])
			}
			old := input[block.ID].Clone()
			input[block.ID] = use[block.ID].Union(output[block.ID].Difference(def[block.ID]))
			change = change || !old.Equal(input[block.ID])
		}
		if !change {
			break
		}
	}
	return input, output
}
