package cfg

import (
	"github.com/willf/bitset"

	"github.com/styczynski/latte-compiler/src/generic_ast"
)

func mapResultsSetsIntoVariables(cfg *CFG, vars []string, varsNames map[string]Variable, ins, outs map[int]*bitset.BitSet) LiveVariablesInfo {
	blocks := cfg.Blocks()
	output := make(map[int]VariableSet, len(blocks))
	input := make(map[int]VariableSet, len(blocks))

	for _, block := range blocks {
		input[block.ID] = VariableSet{}
		output[block.ID] = VariableSet{}
		for i := uint(0); i < outs[block.ID].Len(); i++ {
			if outs[block.ID].Test(i) {
				output[block.ID].Add(varsNames[vars[i]])
			}
		}
		for i := uint(0); i < ins[block.ID].Len(); i++ {
			if ins[block.ID].Test(i) {
				input[block.ID].Add(varsNames[vars[i]])
			}
		}
	}
	return LiveVariablesInfo{
		in:  input,
		out: output,
	}
}

func createUsageBitsets(cfg *CFG) (vars []string, varsNames map[string]Variable, def, use map[int]*bitset.BitSet) {
	blocks := cfg.Blocks()

	def = make(map[int]*bitset.BitSet, len(blocks))
	use = make(map[int]*bitset.BitSet, len(blocks))
	varIndices := make(map[string]uint) // map var to its index in vars
	varsNames = make(map[string]Variable)

	for _, block := range blocks {
		def[block.ID] = new(bitset.BitSet)
		use[block.ID] = new(bitset.BitSet)

		d := GetAllDeclaredVariables(cfg.codeMapping[block.ID], map[generic_ast.TraversableNode]struct{}{}).Copy()
		u := GetAllUsagesVariables(cfg.codeMapping[block.ID], map[generic_ast.TraversableNode]struct{}{}).Copy()

		if block.ID == cfg.Exit {
			for _, dfr := range cfg.OutOfFlowBlocks {
				u.Insert(GetAllUsagesVariables(cfg.codeMapping[dfr], map[generic_ast.TraversableNode]struct{}{}))
			}
		}

		for _, d := range d {
			k, ok := varIndices[d.Name()]
			if !ok {
				k = uint(len(vars))
				varIndices[d.Name()] = k
				varsNames[d.Name()] = d
				vars = append(vars, d.Name())
			}

			def[block.ID].Set(k)
		}

		for _, u := range u {
			k, ok := varIndices[u.Name()]
			if !ok {
				k = uint(len(vars))
				varIndices[u.Name()] = k
				varsNames[u.Name()] = u
				vars = append(vars, u.Name())
			}

			use[block.ID].Set(k)
		}
	}
	return vars, varsNames, def, use
}
