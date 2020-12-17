package cfg

import (
	"github.com/willf/bitset"

	"github.com/styczynski/latte-compiler/src/generic_ast"
)

// liveVarsResultsSets maps in and out bitsets back to their respective vars, such that
// each statement has a set of variables that are live upon entry and a set of
// variables that are live upon exit.
func mapVarsResultSets(cfg *CFG, vars []string, varsNames map[string]Variable, ins, outs map[generic_ast.NormalNode]*bitset.BitSet) LiveVariablesInfo {
	blocks := cfg.Blocks()
	in := make(map[generic_ast.NormalNode]VariableSet, len(blocks))
	out := make(map[generic_ast.NormalNode]VariableSet, len(blocks))

	for _, block := range blocks {
		in[block] = VariableSet{}
		out[block] = VariableSet{}

		for i := uint(0); i < ins[block].Len(); i++ {
			if ins[block].Test(i) {
				in[block].Add(varsNames[vars[i]])
			}
		}

		for i := uint(0); i < outs[block].Len(); i++ {
			if outs[block].Test(i) {
				out[block].Add(varsNames[vars[i]])
			}
		}
	}
	return LiveVariablesInfo{
		in:  in,
		out: out,
	}
}


// defUseBitsets builds def and use bitsets for the given cfg in the context
// of the given loader.PackageInfo. Each entry in the resulting bitsets maps back to
// the same index in the returned vars slice.
func defUseBitsets(cfg *CFG) (vars []string, varsNames map[string]Variable, def, use map[generic_ast.NormalNode]*bitset.BitSet) {
	blocks := cfg.Blocks()

	def = make(map[generic_ast.NormalNode]*bitset.BitSet, len(blocks))
	use = make(map[generic_ast.NormalNode]*bitset.BitSet, len(blocks))
	varIndices := make(map[string]uint) // map var to its index in vars
	varsNames = make(map[string]Variable)

	for _, block := range blocks {
		// prime the def-uses sets
		def[block] = new(bitset.BitSet)
		use[block] = new(bitset.BitSet)

		d := GetAllDeclaredVariables(block).Copy()
		u := GetAllUsagesVariables(block).Copy()

		// use[Exit] = use(each d in cfg.Defers)
		if block == cfg.Exit {
			for _, dfr := range cfg.Defers {
				u.Insert(GetAllUsagesVariables(dfr))
			}
		}

		for _, d := range d {
			// if we have it already, uses that index
			// if we don't, add it to our slice and save its index
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

		//fmt.Printf("For statement %d: %s %s\n", cfg.GetStatementID(block), def[block].String(), use[block].String())
	}
	return vars, varsNames, def, use
}



