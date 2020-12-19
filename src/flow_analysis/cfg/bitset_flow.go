package cfg

import (
	"github.com/willf/bitset"

	"github.com/styczynski/latte-compiler/src/generic_ast"
)

func mapResultsSetsIntoVariables(cfg *CFG, vars []string, varsNames map[string]Variable, ins, outs map[generic_ast.NormalNode]*bitset.BitSet) LiveVariablesInfo {
	blocks := cfg.Blocks()
	output := make(map[generic_ast.NormalNode]VariableSet, len(blocks))
	input := make(map[generic_ast.NormalNode]VariableSet, len(blocks))

	for _, block := range blocks {
		input[block] = VariableSet{}
		output[block] = VariableSet{}
		for i := uint(0); i < outs[block].Len(); i++ {
			if outs[block].Test(i) {
				output[block].Add(varsNames[vars[i]])
			}
		}
		for i := uint(0); i < ins[block].Len(); i++ {
			if ins[block].Test(i) {
				input[block].Add(varsNames[vars[i]])
			}
		}
	}
	return LiveVariablesInfo{
		in:  input,
		out: output,
	}
}


func createUsageBitsets(cfg *CFG) (vars []string, varsNames map[string]Variable, def, use map[generic_ast.NormalNode]*bitset.BitSet) {
	blocks := cfg.Blocks()

	def = make(map[generic_ast.NormalNode]*bitset.BitSet, len(blocks))
	use = make(map[generic_ast.NormalNode]*bitset.BitSet, len(blocks))
	varIndices := make(map[string]uint) // map var to its index in vars
	varsNames = make(map[string]Variable)

	for _, block := range blocks {
		def[block] = new(bitset.BitSet)
		use[block] = new(bitset.BitSet)

		d := GetAllDeclaredVariables(block, map[generic_ast.TraversableNode]struct{}{}).Copy()
		u := GetAllUsagesVariables(block, map[generic_ast.TraversableNode]struct{}{}).Copy()

		if block == cfg.Exit {
			for _, dfr := range cfg.OutOfFlowBlocks {
				u.Insert(GetAllUsagesVariables(dfr, map[generic_ast.TraversableNode]struct{}{}))
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

			def[block].Set(k)
		}

		for _, u := range u {
			k, ok := varIndices[u.Name()]
			if !ok {
				k = uint(len(vars))
				varIndices[u.Name()] = k
				varsNames[u.Name()] = u
				vars = append(vars, u.Name())
			}

			use[block].Set(k)
		}
	}
	return vars, varsNames, def, use
}



