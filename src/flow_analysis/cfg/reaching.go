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


// File implements reaching definitions data flow analysis for a
// control flow graph with statement granularity.
//
// based on algo from ch 9.2, p.607 Dragonbook, v2.2,
// "Iterative algorithm to compute reaching definitions":
//
// OUT[ENTRY] = {};
// for(each basic block B other than ENTRY) OUT[B] = {};
// for(changes to any OUT occur)
//    for(each basic block B other than ENTRY) {
//      IN[B] = Union(P a pred of B) OUT[P];
//      OUT[B] = gen[b] Union (IN[B] - kill[b]);
//    }

// DefUse builds reaching definitions for a given control flow graph, returning
// a map that maps each statement that defines a variable (i.e., declares or
// assigns it) to the set of statements that use that variable.
//
// Note: An assignment to a struct field or array element is treated as both a
// use and a definition of that variable, since only part of its value is
// assigned.  For analysis purposes, it's treated as though the entire value
// is read, then part of it is modified, then the entire value is assigned back
// to the variable.  (This is necessary for the analysis to produce correct
// results.)
//
// No nodes from the cfg.Defers list will be returned in the output of
// this function as they are disjoint from a cfg's blocks.
// For analyzing the statements in the cfg.Defers list, each defer
// should be treated as though it has the same in and out sets as the cfg.Exit node.
func DefUse(cfg *CFG) ReachingVariablesInfo {
	blocks, gen, kill := genKillBitsets(cfg)
	ins, _ := reachingDefBitsets(cfg, gen, kill)
	return ReachingVariablesInfo{
		usagesMap: defUseResultSet(blocks, ins),
	}
}

// DefsReaching builds reaching definitions for a given control flow graph,
// returning the set of statements that define a variable (i.e., declare or
// assign it) where that definition reaches the given statement.
func DefsReaching(stmt generic_ast.NormalNode, cfg *CFG) map[generic_ast.NormalNode]struct{} {
	blocks, gen, kill := genKillBitsets(cfg)
	ins, _ := reachingDefBitsets(cfg, gen, kill)
	return defsReachingResultSet(stmt, blocks, ins)
}

// genKillBitsets builds the gen and kill bitsets for each block in a cfg,
// these are used to compute reaching definitions.
func genKillBitsets(cfg *CFG) (blocks []generic_ast.NormalNode, gen, kill map[generic_ast.NormalNode]*bitset.BitSet) {
	okills := make(map[Variable]*bitset.BitSet)
	gen = make(map[generic_ast.NormalNode]*bitset.BitSet)
	kill = make(map[generic_ast.NormalNode]*bitset.BitSet)
	blocks = cfg.Blocks()

	for _, b := range blocks { // prime
		gen[b] = new(bitset.BitSet)
		kill[b] = new(bitset.BitSet)
	}

	// Iterate over all blocks twice, because a block may not know the entirety of what
	// it kills until all blocks have been iterated over.
	for i := 0; i < 2; i++ {
		for j, block := range blocks {
			j := uint(j)

			def := GetAllDeclaredVariables(block)

			for _, d := range def {
				if _, ok := okills[d]; !ok {
					okills[d] = new(bitset.BitSet)
				}
				gen[block].Set(j) // GEN this obj
				okills[d].Set(j)  // KILL this obj for everyone else
				// our kills are KILL[obj] - GEN[B]
				kill[block] = kill[block].Union(okills[d]).Difference(gen[block])
			}
		}
	}
	return blocks, gen, kill
}

// reachingDefBitsets will compute the reaching definitions in and out sets from gen and kill bitsets.
func reachingDefBitsets(cfg *CFG, gen, kill map[generic_ast.NormalNode]*bitset.BitSet) (in, out map[generic_ast.NormalNode]*bitset.BitSet) {
	in = make(map[generic_ast.NormalNode]*bitset.BitSet)
	out = make(map[generic_ast.NormalNode]*bitset.BitSet)
	blocks := cfg.Blocks()

	// OUT[ENTRY] = {};
	// for(each basic block B other than ENTRY) OUT[B} = {};
	for i := 0; i < len(blocks); i++ {
		block := blocks[i]
		in[block] = new(bitset.BitSet)
		out[block] = new(bitset.BitSet)
		if block == cfg.Entry {
			blocks = append(blocks[:i], blocks[i+1:]...)
			i--
		}
	}

	// for(changes to any OUT occur)
	for {
		var changed bool

		// for(each basic block B other than ENTRY) {
		for _, block := range blocks {

			// IN[B] = Union(P a pred of B) OUT[P];
			for _, p := range cfg.Preds(block) {
				in[block].InPlaceUnion(out[p])
			}

			old := out[block].Clone()

			// OUT[B] = gen[b] Union (IN[B] - kill[b]);
			out[block] = gen[block].Union(in[block].Difference(kill[block]))

			changed = changed || !old.Equal(out[block])
		}

		if !changed {
			break
		}
	}
	return in, out
}

// defUseResultSet maps reaching definitions in bitsets back to their corresponding statements, using
// this information to determine use-def and def-use information.
// blocks should be the blocks used to generate the analyses, as their indices are what will be used to map
// bits in each bitset back to the corresponding statement.
func defUseResultSet(blocks []generic_ast.NormalNode, ins map[generic_ast.NormalNode]*bitset.BitSet) map[generic_ast.NormalNode]ReachedBlocks {
	du := make(map[generic_ast.NormalNode]ReachedBlocks)

	// map bits from in and out sets back to corresponding blocks (with cfg.Entry)
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

// defUseResultSet maps reaching definitions in bitsets back to their
// corresponding statements, returning the set of definition statements that
// reach the given stmt.
func defsReachingResultSet(stmt generic_ast.NormalNode, blocks []generic_ast.NormalNode, ins map[generic_ast.NormalNode]*bitset.BitSet) map[generic_ast.NormalNode]struct{} {
	result := make(map[generic_ast.NormalNode]struct{})
	insStmt, found := ins[stmt]
	if !found {
		panic("stmt not in CFG")
	}
	for i, ok := uint(0), true; ok; i++ {
		if i, ok = insStmt.NextSet(i); ok {
			result[blocks[i]] = struct{}{}
		}
	}
	return result
}