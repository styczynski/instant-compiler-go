package generic_ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
)

type VirtualNode struct {
	parent TraversableNode
	nodeType VirtualNodeType
}

type VirtualNodeType string

const (
	V_NODE_CODE_END VirtualNodeType = "END"
	V_NODE_ENTRY VirtualNodeType = "ENTRY"
	V_NODE_EXIT VirtualNodeType = "EXIT"
)

func CreateVirtualNode(nodeType VirtualNodeType) *VirtualNode {
	return &VirtualNode{
		parent: nil,
		nodeType: nodeType,
	}
}

func (n *VirtualNode) Name() string {
	return string(n.nodeType)
}

func (n *VirtualNode) Print(c *context.ParsingContext) string {
	return string(n.nodeType)
}

func (n *VirtualNode) GetChildren() []TraversableNode {
	return []TraversableNode{}
}

func (n *VirtualNode) GetNode() interface{} {
	return n	
}

func (n *VirtualNode) Parent() TraversableNode {
	return n.parent
}

func (n *VirtualNode) OverrideParent(node TraversableNode) {
	n.parent = node
}

func (n *VirtualNode) Begin() lexer.Position {
	dummyOffset := 0
	if n.nodeType == V_NODE_ENTRY {
		dummyOffset = -1
	} else if n.nodeType == V_NODE_CODE_END {
		dummyOffset = -2
	} else {
		dummyOffset = -3
	}
	return lexer.Position{
		Filename: "",
		Offset:   dummyOffset,
		Line:     dummyOffset,
		Column:   dummyOffset,
	}
}

func (n *VirtualNode) End() lexer.Position {
	return n.Begin()
}
