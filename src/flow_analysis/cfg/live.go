package cfg

import (
	"github.com/willf/bitset"

	"github.com/styczynski/latte-compiler/src/generic_ast"
)

// Copyright 2015-2018 Auburn University and others. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


// File defines live variables analysis for a statement
// level control flow graph. Defer has quirks, see LiveVars func.
//
// based on algo from ch 9.2, p.610 Dragonbook, v2.2,
// "Iterative algorithm to compute live variables":
//
// IN[EXIT] = use[each D in Defers];
// for(each basic block B other than EXIT) IN[B} = {};
// for(changes to any IN occur)
//    for(each basic block B other than EXIT) {
//      OUT[B] = Union(S a successor of B) IN[S];
//      IN[B] = use[b] Union (OUT[B] - def[b]);
//    }

// NOTE: for extract function: defers in the block to extract can
// (probably?) be extracted if all variables used in the defer statement are
// not live at the beginning and the end of the block to extract

// LiveAt returns the in and out set of live variables for each block in
// a given control flow graph (cfg) in the context of a loader.Program,
// including the cfg.Entry and cfg.Exit nodes.
//
// The traditional approach of holding the live variables at the exit node
// to the empty set has been deviated from in order to handle defers.
// The live variables in set of the cfg.Exit node will be set to the variables used
// in all cfg.Defers. No liveness is analyzed for the cfg.Defers themselves.
//
// More formally:
//  IN[EXIT] = USE(each d in cfg.Defers)
//  OUT[EXIT] = {}

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
	vars, varsNames, def, use := defUseBitsets(cfg)
	ins, outs := liveVarsBitsets(cfg, def, use)
	return mapVarsResultSets(cfg, vars, varsNames, ins, outs)
}

// liveVarsBitsets generatates live variable analysis in and out bitsets from def and use sets
func liveVarsBitsets(cfg *CFG, def, use map[generic_ast.NormalNode]*bitset.BitSet) (in, out map[generic_ast.NormalNode]*bitset.BitSet) {
	blocks := cfg.Blocks()
	in = make(map[generic_ast.NormalNode]*bitset.BitSet, len(blocks))
	out = make(map[generic_ast.NormalNode]*bitset.BitSet, len(blocks))
	// for(each basic block B) IN[B} = {};
	for _, block := range blocks {
		in[block] = new(bitset.BitSet)
		out[block] = new(bitset.BitSet)
	}

	// for(changes to any IN occur)
	for {
		var change bool

		// for(each basic block B) {
		for _, block := range blocks {

			// OUT[B] = Union(S a succ of B) IN[S]
			for _, s := range cfg.Succs(block) {
				out[block].InPlaceUnion(in[s])
			}

			old := in[block].Clone()

			// IN[B] = uses[B] U (OUT[B] - def[B])
			in[block] = use[block].Union(out[block].Difference(def[block]))
			//fmt.Printf("Inspect statement %d: %s\n", cfg.GetStatementID(block), in[block].String())

			change = change || !old.Equal(in[block])
		}

		if !change {
			break
		}
	}
	return in, out
}