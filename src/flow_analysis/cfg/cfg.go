// Copyright 2015 Auburn University. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cfg provides intraprocedural control flow graphs (CFGs) with
// statement-level granularity, i.e., CFGs whose nodes correspond 1-1 to the
// Stmt nodes from an abstract syntax tree.
package cfg

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/styczynski/latte-compiler/src/generic_ast"
)

type CFGBuilder interface {
	GetPrev() []generic_ast.NormalNode
	AddSucc(current generic_ast.NormalNode)
	AddBranch(branch generic_ast.NormalNode)
	Branches() []generic_ast.NormalNode
	UpdatePrev(nodes []generic_ast.NormalNode)
	BuildNode(node generic_ast.NormalNode)
	BuildBlock(node generic_ast.TraversableNode)
	Exit() generic_ast.NormalNode
}

type NodeWithControlInformation interface {
	BuildFlowGraph(builder CFGBuilder)
}

// This package can be used to construct a control flow graph from an abstract syntax tree (go/ast).
// This is done by traversing a list of statements (likely from a block)
// depth-first and creating an adjacency list, implemented as a map of blocks.
// Adjacent blocks are stored as predecessors and successors separately for
// control flow information. Any defers encountered while traversing the ast
// will be added to a slice that can be accessed from CFG. Their behavior is such
// that they may or may not be flowed to, potentially multiple times, after Exit.
// This behavior is dependant upon in what control structure they were found,
// i.e. if/for body may never be flowed to.

// TODO(you): defers are lazily done currently. If needed, could likely use a more robust
//  implementation wherein they are represented as a graph after Exit.
// TODO(reed): closures, go func() ?

// CFG defines a control flow graph with statement-level granularity, in which
// there is a 1-1 correspondence between a block in the CFG and an generic_ast.NormalNode.
type CFG struct {
	// Sentinel nodes for single-entry, single-exit CFG. Not in original AST.
	Entry, Exit, CodeEnd generic_ast.NormalNode
	// All defers found in CFG, disjoint from blocks. May be flowed to after Exit.
	Defers []generic_ast.NormalNode
	blocks map[generic_ast.NormalNode]*block
	blocksOrder []generic_ast.NormalNode
	blocksIDs map[generic_ast.NormalNode]int
}

func getBeginPos(node generic_ast.NormalNode) (int, int) {
	return node.Begin().Line, node.Begin().Column
}

func getEndPos(node generic_ast.NormalNode) (int, int) {
	return node.End().Line, node.End().Column
}

type ControlFlowGraphVisitor func (cfg *CFG, block *block, next func(node generic_ast.NormalNode))

func (cfg *CFG) VisitGraph(node generic_ast.NormalNode, visitor ControlFlowGraphVisitor) {
	block := cfg.blocks[node]
	next := func(node generic_ast.NormalNode) {
		cfg.VisitGraph(node, visitor)
	}
	for _, child := range block.succs {
		visitor(cfg, cfg.blocks[child], next)
	}
}

func (cfg *CFG) GetAllEndGateways() []generic_ast.NormalNode {
	gateways := []generic_ast.NormalNode{}
	visitedIDs := map[int]struct{}{}
	cfg.VisitGraph(cfg.Entry, func(cfg *CFG, block *block, next func(node generic_ast.NormalNode)) {
		if _, wasVisited := visitedIDs[block.ID]; wasVisited {
			return
		}
		visitedIDs[block.ID] = struct{}{}
		for _, child := range block.succs {
			if child == cfg.CodeEnd {
				gateways = append(gateways, block.stmt)
			}
		}
		next(block.stmt)
	})
	return gateways
}

func (cfg *CFG) ReplaceBlock(old generic_ast.NormalNode, new generic_ast.NormalNode) {

	for block, blockDescription := range cfg.blocks {
		if blockDescription != nil {
			if blockDescription.stmt == old {
				blockDescription.stmt = new
			}
			for j, nestedNode := range blockDescription.succs {
				if nestedNode == old {
					blockDescription.succs[j] = new
				}
			}
			for j, nestedNode := range blockDescription.preds {
				if nestedNode == old {
					blockDescription.preds[j] = new
				}
			}
			if block == old {
				delete(cfg.blocks, old)
				cfg.blocks[new] = blockDescription
			}
		}
	}
	for i, block := range cfg.Defers {
		if block == old {
			cfg.Defers[i] = new
		}
	}
	for i, block := range cfg.blocksOrder {
		if block == old {
			cfg.blocksOrder[i] = new
		}
	}
	for block, id := range cfg.blocksIDs {
		if block == old {
			delete(cfg.blocksIDs, old)
			cfg.blocksIDs[new] = id
		}
	}
}

type OrderBlocksByPosition []generic_ast.NormalNode

func (s OrderBlocksByPosition) Len() int {
	return len(s)
}

func (s OrderBlocksByPosition) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s OrderBlocksByPosition) Less(i, j int) bool {
	if isNilNode(s[i]) {
		return false
	}
	if isNilNode(s[j]) {
		return true
	}
	if virtualNode, ok := s[i].(*generic_ast.VirtualNode); ok {
		if virtualNode.Name() == "ENTRY" {
			return true
		}
	}
	if virtualNode, ok := s[i].(*generic_ast.VirtualNode); ok {
		if virtualNode.Name() == "EXIT" {
			return false
		}
	}
	b1, e1 := getBeginPos(s[i])
	b2, e2 := getEndPos(s[j])
	if b1 == b2 {
		return e1 < e2
	}
	return b1 < b2
}

type block struct {
	stmt  generic_ast.NormalNode
	preds []generic_ast.NormalNode
	succs []generic_ast.NormalNode
	ID int
}

// FromStmts returns the control-flow graph for the given sequence of statements.
func FromStmts(s []generic_ast.NormalNode) *CFG {
	return newBuilder().build(s)
}

// Preds returns a slice of all immediate predecessors for the given statement.
// May include Entry node.
func (c *CFG) Preds(s generic_ast.NormalNode) []generic_ast.NormalNode {
	return c.blocks[s].preds
}

// Succs returns a slice of all immediate successors to the given statement.
// May include Exit node.
func (c *CFG) Succs(s generic_ast.NormalNode) []generic_ast.NormalNode {
	if _, ok := c.blocks[s]; !ok {
		panic("Missing succ")
	}
	return c.blocks[s].succs
}

// Blocks returns a slice of all blocks in a CFG, including the Entry and Exit nodes.
// The blocks are roughly in the order they appear in the source code.
func (c *CFG) Blocks() []generic_ast.NormalNode {
	v := []generic_ast.NormalNode{}
	for _, block := range c.blocksOrder {
		v = append(v, block)
	}
	return v
}

func (c *CFG) GetStatementID(node generic_ast.NormalNode) int {
	return c.blocksIDs[node]
}

// type for sorting statements by their starting positions in the source code
type stmtSlice []generic_ast.NormalNode

func (n stmtSlice) Len() int      { return len(n) }
func (n stmtSlice) Swap(i, j int) { n[i], n[j] = n[j], n[i] }
func (n stmtSlice) Less(i, j int) bool {
	return n[i].Begin().Offset < n[j].Begin().Offset
}

func (c *CFG) Sort(stmts []generic_ast.NormalNode) {
	sort.Sort(stmtSlice(stmts))
}

func (c *CFG) PrintDot(f io.Writer, addl func(n generic_ast.NormalNode) string) {
	fmt.Fprintf(f, `digraph mgraph {
mode="heir";
splines="ortho";

`)
	blocks := c.Blocks()
	c.Sort(blocks)
	for _, from := range blocks {
		succs := c.Succs(from)
		c.Sort(succs)
		for _, to := range succs {
			fmt.Fprintf(f, "\t\"%s\" -> \"%s\"\n",
				c.printVertex(from, addl(from)),
				c.printVertex(to, addl(to)))
		}
	}
	fmt.Fprintf(f, "}\n")
}

func (c *CFG) printVertex(stmt generic_ast.NormalNode, addl string) string {
	switch stmt {
	case c.Entry:
		return "ENTRY"
	case c.Exit:
		return "EXIT"
	case nil:
		return ""
	}
	addl = strings.Replace(addl, "\n", "\\n", -1)
	if addl != "" {
		addl = "\\n" + addl
	}
	return fmt.Sprintf("%s - line %d%s",
		stmt.Print(nil),
		stmt.Begin().Line,
		addl)
}

type builder struct {
	blocks      map[generic_ast.NormalNode]*block
	prev        []generic_ast.NormalNode        // blocks to hook up to current block
	branches    []generic_ast.NormalNode // accumulated branches from current inner blocks
	entry, exit, codeEnd generic_ast.NormalNode      // single-entry, single-exit nodes
	defers      []generic_ast.NormalNode  // all defers encountered
}

func newBuilder() *builder {
	// The ENTRY and EXIT nodes are given positions -2 and -1 so cfg.Sort
	// will work correct: ENTRY will always be first, followed by EXIT,
	// followed by the other CFG nodes.
	return &builder{
		blocks: map[generic_ast.NormalNode]*block{},
		entry:  generic_ast.CreateVirtualNode(generic_ast.V_NODE_ENTRY),
		exit:   generic_ast.CreateVirtualNode(generic_ast.V_NODE_EXIT),
		codeEnd: generic_ast.CreateVirtualNode(generic_ast.V_NODE_CODE_END),
	}
}

// build runs buildBlock on the given block (traversing nested statements), and
// adds entry and exit nodes.
func (b *builder) build(s []generic_ast.NormalNode) *CFG {

	b.prev = []generic_ast.NormalNode{b.entry}
	b.buildBlock(s)
	b.AddSucc(b.codeEnd)
	b.prev = []generic_ast.NormalNode{b.codeEnd}
	b.AddSucc(b.exit)


	sortedExprs := OrderBlocksByPosition{}
	for e, _ := range b.blocks {
		sortedExprs = append(sortedExprs, e)
	}
	sort.Sort(sortedExprs)

	exprIDs := map[generic_ast.NormalNode]int{}
	freeID := 1
	freeBlockID := 1
	for _, expr := range sortedExprs {
		if expr != nil {
			exprIDs[expr] = freeID
			if block, ok := b.blocks[expr]; ok {
				if block.ID == 0 {
					block.ID = freeBlockID
					freeBlockID++
				}
			}
		}
		freeID++
	}

	return &CFG{
		Entry:       b.entry,
		Exit:        b.exit,
		CodeEnd:     b.codeEnd,
		Defers:      b.defers,
		blocks:      b.blocks,
		blocksOrder: sortedExprs,
		blocksIDs:   exprIDs,
	}
}

func (b *builder) AddBranch(branch generic_ast.NormalNode) {
	b.branches = append(b.branches, branch)
}

func (b *builder) Branches() []generic_ast.NormalNode {
	return b.branches
}

func (b *builder) BuildNode(node generic_ast.NormalNode) {
	if isNilNode(node) {
		return
	}
	if nodeWithControlInformation, ok := node.(NodeWithControlInformation); ok {
		nodeWithControlInformation.BuildFlowGraph(b)
	} else {
		// Default
		b.AddSucc(node)
		b.UpdatePrev([]generic_ast.NormalNode{ node })
	}
}

func (b *builder) BuildBlock(node generic_ast.TraversableNode) {
	for _, child := range node.GetChildren() {
		b.BuildNode(child.(generic_ast.NormalNode))
	}
}

func (b *builder) Exit() generic_ast.NormalNode {
	return b.exit
}

func (b *builder) GetPrev() []generic_ast.NormalNode{
	return b.prev
}

func (b *builder) UpdatePrev(nodes []generic_ast.NormalNode) {
	b.prev = nodes
}

// AddSucc adds a control flow edge from all previous blocks to the block for
// the given statement.
func (b *builder) AddSucc(current generic_ast.NormalNode) {
	cur := b.block(current)

	for _, p := range b.prev {
		p := b.block(p)
		p.succs = appendNoDuplicates(p.succs, cur.stmt)
		cur.preds = appendNoDuplicates(cur.preds, p.stmt)
	}
}

func appendNoDuplicates(list []generic_ast.NormalNode, stmt generic_ast.NormalNode) []generic_ast.NormalNode {
	for _, s := range list {
		if s == stmt {
			return list
		}
	}
	return append(list, stmt)
}

// block returns a block for the given statement, creating one and inserting it
// into the CFG if it doesn't already exist.
func (b *builder) block(s generic_ast.NormalNode) *block {
	bl, ok := b.blocks[s]
	if !ok {
		bl = &block{stmt: s}
		b.blocks[s] = bl
	}
	return bl
}

// buildBlock iterates over a slice of statements (typically the statements
// from an ast.BlockStmt), adding them successively to the CFG.  Upon return,
// b.prev is set to the control exits of the last statement.
func (b *builder) buildBlock(block []generic_ast.NormalNode) {
	for _, stmt := range block {
		b.BuildNode(stmt)
	}
}
