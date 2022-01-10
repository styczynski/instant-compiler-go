package cfg

import (
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
	Entry, Exit, CodeEnd int
	OutOfFlowBlocks      []int
	blocks               map[int]*Block
	blocksOrder          []int
	codeMapping          map[int]CFGCodeNode
}

func (c *CFG) GetBlockCode(blockID int) CFGCodeNode {
	return c.codeMapping[blockID]
}

type CFGCodeNode interface {
	generic_ast.NormalNode
}

func getBeginPos(node generic_ast.NormalNode) (int, int) {
	return node.Begin().Line, node.Begin().Column
}

func getEndPos(node generic_ast.NormalNode) (int, int) {
	return node.End().Line, node.End().Column
}

type ControlFlowGraphVisitor func(cfg *CFG, block *Block, next func(blockID int))

func (cfg *CFG) OverrideBlockCode(blockID int, newNode CFGCodeNode) {
	cfg.codeMapping[blockID] = newNode
}

func (cfg *CFG) GetAllEndGateways() []CFGCodeNode {
	gateways := []CFGCodeNode{}
	visitedIDs := map[int]struct{}{}
	cfg.VisitGraph(cfg.Entry, func(cfg *CFG, block *Block, next func(blockID int)) {
		if _, wasVisited := visitedIDs[block.ID]; wasVisited {
			return
		}
		visitedIDs[block.ID] = struct{}{}
		for _, childID := range block.succs {
			if childID == cfg.CodeEnd {
				gateways = append(gateways, cfg.codeMapping[block.ID])
			}
			next(childID)
		}
	})
	return gateways
}

func (cfg *CFG) VisitGraph(blockID int, visitor ControlFlowGraphVisitor) {
	block := cfg.blocks[blockID]
	next := func(blockID int) {
		cfg.VisitGraph(blockID, visitor)
	}
	visitor(cfg, cfg.blocks[block.ID], next)
}

func (cfg *CFG) ListBlockIDs() []int {
	return cfg.blocksOrder
}

func (cfg *CFG) RemoveBlocks(idsToRemove map[int]struct{}) {
	newBlocksOrder := []int{}
	for _, blockID := range cfg.ListBlockIDs() {
		if _, ok := idsToRemove[blockID]; !ok {
			newBlocksOrder = append(newBlocksOrder, blockID)
		}
	}
	cfg.blocksOrder = newBlocksOrder

	newOutOfFlowBlocks := []int{}
	for _, blockID := range cfg.OutOfFlowBlocks {
		if _, ok := idsToRemove[blockID]; !ok {
			newOutOfFlowBlocks = append(newOutOfFlowBlocks, blockID)
		}
	}
	cfg.OutOfFlowBlocks = newOutOfFlowBlocks

	for _, block := range cfg.blocks {
		for index, succ := range block.succs {
			block.succs[index] = cfg.blocks[succ].ID
		}
		for index, pred := range block.preds {
			block.preds[index] = cfg.blocks[pred].ID
		}
	}
}

func (cfg *CFG) ResolveID(blockID int) int {
	return cfg.blocks[blockID].ID
}

func (cfg *CFG) ShadowBlock(blockID int, newBlock *Block) {
	cfg.blocks[blockID] = newBlock
}

func (cfg *CFG) GetBlock(blockID int) *Block {
	return cfg.blocks[blockID]
}

func (cfg *CFG) ReplaceBlockCode(blockID int, newCode CFGCodeNode) {

	cfg.codeMapping[blockID] = newCode

	// for block, blockDescription := range cfg.blocks {
	// 	if blockDescription != nil {
	// 		if blockDescription.stmt == old {
	// 			blockDescription.stmt = new
	// 		}
	// 		for j, nestedNode := range blockDescription.succs {
	// 			if nestedNode == old {
	// 				blockDescription.succs[j] = new
	// 			}
	// 		}
	// 		for j, nestedNode := range blockDescription.preds {
	// 			if nestedNode == old {
	// 				blockDescription.preds[j] = new
	// 			}
	// 		}
	// 		if block == old {
	// 			delete(cfg.blocks, old)
	// 			cfg.blocks[new] = blockDescription
	// 		}
	// 	}
	// }
	// for i, block := range cfg.OutOfFlowBlocks {
	// 	if block == old {
	// 		cfg.OutOfFlowBlocks[i] = new
	// 	}
	// }
	// for i, block := range cfg.blocksOrder {
	// 	if block == old {
	// 		cfg.blocksOrder[i] = new
	// 	}
	// }
	// for block, id := range cfg.blocksIDs {
	// 	if block == old {
	// 		delete(cfg.blocksIDs, old)
	// 		cfg.blocksIDs[new] = id
	// 	}
	// }
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

type Block struct {
	preds []int
	succs []int
	ID    int
}

func (b Block) GetPreds() []int {
	return b.preds
}

func (b Block) GetSuccs() []int {
	return b.succs
}

type BuilderBlock struct {
	stmt  generic_ast.NormalNode
	preds []generic_ast.NormalNode
	succs []generic_ast.NormalNode
	ID    int
}

func CreateCFGFromNodes(s []generic_ast.NormalNode) *CFG {
	return newBuilder().build(s)
}

func (c *CFG) BlockPredecessors(blockID int) []*Block {
	preds := []*Block{}
	for _, blockID := range c.blocks[blockID].preds {
		preds = append(preds, c.blocks[blockID])
	}
	return preds
}

func (c *CFG) BlockSuccessors(blockID int) []*Block {
	succs := []*Block{}
	for _, blockID := range c.blocks[blockID].succs {
		succs = append(succs, c.blocks[blockID])
	}
	return succs
}

func (c *CFG) Blocks() []*Block {
	v := []*Block{}
	for _, blockID := range c.blocksOrder {
		v = append(v, c.blocks[blockID])
	}
	return v
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
	blocks               map[generic_ast.NormalNode]*BuilderBlock
	previousBlocks       []generic_ast.NormalNode
	outOfFlowBlocks      []generic_ast.NormalNode
	codeBranches         []generic_ast.NormalNode
	entry, exit, codeEnd generic_ast.NormalNode
}

func newBuilder() *cfgGraphBuilder {
	return &cfgGraphBuilder{
		blocks:  map[generic_ast.NormalNode]*BuilderBlock{},
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

	codeMapping := map[int]CFGCodeNode{}

	codeMapping[b.blocks[b.entry].ID] = b.entry
	codeMapping[b.blocks[b.exit].ID] = b.exit
	codeMapping[b.blocks[b.codeEnd].ID] = b.codeEnd

	for _, block := range b.blocks {
		codeMapping[block.ID] = block.stmt
	}

	outOfFlowBlocksIDs := []int{}
	for _, block := range b.outOfFlowBlocks {
		outOfFlowBlocksIDs = append(outOfFlowBlocksIDs, b.blocks[block].ID)
	}

	blocksOrder := []int{}
	for _, block := range sortedExprs {
		blocksOrder = append(blocksOrder, b.blocks[block].ID)
	}

	allBlocks := map[int]*Block{}
	for stmt, builderBlock := range b.blocks {
		succs := []int{}
		for _, succ := range builderBlock.succs {
			succs = append(succs, b.blocks[succ].ID)
		}
		preds := []int{}
		for _, pred := range builderBlock.preds {
			preds = append(preds, b.blocks[pred].ID)
		}
		allBlocks[b.blocks[stmt].ID] = &Block{
			ID:    builderBlock.ID,
			preds: preds,
			succs: succs,
		}
	}

	return &CFG{
		Entry:           b.blocks[b.entry].ID,
		Exit:            b.blocks[b.exit].ID,
		CodeEnd:         b.blocks[b.codeEnd].ID,
		OutOfFlowBlocks: outOfFlowBlocksIDs,
		blocks:          allBlocks,
		blocksOrder:     blocksOrder,
		codeMapping:     codeMapping,
	}
}

func (b *cfgGraphBuilder) AddBranch(branch generic_ast.NormalNode) {
	b.codeBranches = append(b.codeBranches, branch)
}

func (b *cfgGraphBuilder) block(s generic_ast.NormalNode) *BuilderBlock {
	bl, ok := b.blocks[s]
	if !ok {
		bl = &BuilderBlock{stmt: s}
		// TODO: Shit
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
