package cfg

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/styczynski/latte-compiler/src/generic_ast"
)

type CFGBuilder interface {
	GetPrev() []generic_ast.NormalNode
	AddBlockSuccesor(current generic_ast.NormalNode)
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

type CFG struct {
	Entry, Exit, CodeEnd generic_ast.NormalNode
	OutOfFlowBlocks      []generic_ast.NormalNode
	blocks               map[generic_ast.NormalNode]*block
	blocksOrder          []generic_ast.NormalNode
	blocksIDs            map[generic_ast.NormalNode]int
}

func getBeginPos(node generic_ast.NormalNode) (int, int) {
	return node.Begin().Line, node.Begin().Column
}

func getEndPos(node generic_ast.NormalNode) (int, int) {
	return node.End().Line, node.End().Column
}

type ControlFlowGraphVisitor func(cfg *CFG, block *block, next func(node generic_ast.NormalNode))

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
	for i, block := range cfg.OutOfFlowBlocks {
		if block == old {
			cfg.OutOfFlowBlocks[i] = new
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
	ID    int
}

func CreateCFGFromNodes(s []generic_ast.NormalNode) *CFG {
	return newBuilder().build(s)
}

func (c *CFG) Exists(s generic_ast.NormalNode) bool {
	fmt.Printf("check exists for %v\n", reflect.TypeOf(s))
	block, ok := c.blocks[s]
	return ok && block != nil
}

func (c *CFG) BlockPredecessors(s generic_ast.NormalNode) []generic_ast.NormalNode {
	return c.blocks[s].preds
}

func (c *CFG) BlockSuccessors(s generic_ast.NormalNode) []generic_ast.NormalNode {
	if _, ok := c.blocks[s]; !ok {
		panic("Missing succ")
	}
	return c.blocks[s].succs
}

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

type nodesSlice []generic_ast.NormalNode

func (n nodesSlice) Len() int      { return len(n) }
func (n nodesSlice) Swap(i, j int) { n[i], n[j] = n[j], n[i] }
func (n nodesSlice) Less(i, j int) bool {
	return n[i].Begin().Offset < n[j].Begin().Offset
}

func (c *CFG) Sort(stmts []generic_ast.NormalNode) {
	sort.Sort(nodesSlice(stmts))
}

type cfgGraphBuilder struct {
	blocks               map[generic_ast.NormalNode]*block
	previousBlocks       []generic_ast.NormalNode
	outOfFlowBlocks      []generic_ast.NormalNode
	codeBranches         []generic_ast.NormalNode
	entry, exit, codeEnd generic_ast.NormalNode
}

func newBuilder() *cfgGraphBuilder {
	return &cfgGraphBuilder{
		blocks:  map[generic_ast.NormalNode]*block{},
		entry:   generic_ast.CreateVirtualNode(generic_ast.V_NODE_ENTRY),
		exit:    generic_ast.CreateVirtualNode(generic_ast.V_NODE_EXIT),
		codeEnd: generic_ast.CreateVirtualNode(generic_ast.V_NODE_CODE_END),
	}
}

func (b *cfgGraphBuilder) build(s []generic_ast.NormalNode) *CFG {

	b.previousBlocks = []generic_ast.NormalNode{b.entry}
	b.buildBlock(s)
	b.AddBlockSuccesor(b.codeEnd)
	b.previousBlocks = []generic_ast.NormalNode{b.codeEnd}
	b.AddBlockSuccesor(b.exit)

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
		Entry:           b.entry,
		Exit:            b.exit,
		CodeEnd:         b.codeEnd,
		OutOfFlowBlocks: b.outOfFlowBlocks,
		blocks:          b.blocks,
		blocksOrder:     sortedExprs,
		blocksIDs:       exprIDs,
	}
}

func (b *cfgGraphBuilder) AddBranch(branch generic_ast.NormalNode) {
	b.codeBranches = append(b.codeBranches, branch)
}

func (b *cfgGraphBuilder) block(s generic_ast.NormalNode) *block {
	bl, ok := b.blocks[s]
	if !ok {
		bl = &block{stmt: s}
		b.blocks[s] = bl
	}
	return bl
}

func (b *cfgGraphBuilder) Branches() []generic_ast.NormalNode {
	return b.codeBranches
}

func (b *cfgGraphBuilder) buildBlock(block []generic_ast.NormalNode) {
	for _, stmt := range block {
		b.BuildNode(stmt)
	}
}

func (b *cfgGraphBuilder) BuildNode(node generic_ast.NormalNode) {
	if isNilNode(node) {
		return
	}
	if nodeWithControlInformation, ok := node.(NodeWithControlInformation); ok {
		nodeWithControlInformation.BuildFlowGraph(b)
	} else {
		// Default
		b.AddBlockSuccesor(node)
		b.UpdatePrev([]generic_ast.NormalNode{node})
	}
}

func (b *cfgGraphBuilder) BuildBlock(node generic_ast.TraversableNode) {
	for _, child := range node.GetChildren() {
		b.BuildNode(child.(generic_ast.NormalNode))
	}
}

func (b *cfgGraphBuilder) Exit() generic_ast.NormalNode {
	return b.exit
}

func (b *cfgGraphBuilder) GetPrev() []generic_ast.NormalNode {
	return b.previousBlocks
}

func (b *cfgGraphBuilder) UpdatePrev(nodes []generic_ast.NormalNode) {
	b.previousBlocks = nodes
}

func (b *cfgGraphBuilder) AddBlockSuccesor(current generic_ast.NormalNode) {
	cur := b.block(current)

	for _, p := range b.previousBlocks {
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
