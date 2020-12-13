package generic_ast

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/styczynski/latte-compiler/src/parser/context"
)

type ComplexASTNode struct {
	BaseASTNode
	Tokens []lexer.Token
}

type BaseASTNode struct {
	Pos lexer.Position
	EndPos lexer.Position
}

func (ast *BaseASTNode) End() lexer.Position {
	return ast.EndPos
}

func (ast *BaseASTNode) Begin() lexer.Position {
	return ast.Pos
}

type NodeWithPosition interface {
	Begin() lexer.Position
	End() lexer.Position
}

type PrintableNode interface {
	Print(c *context.ParsingContext) string
}

type TraversableNode interface {
	NodeWithPosition
	GetChildren() []TraversableNode
	GetNode() interface{}
	Parent() TraversableNode
	OverrideParent(node TraversableNode)
}

type NormalNode interface {
	PrintableNode
	TraversableNode
}

type TraversableNodeToken struct {
	Token string
	BeginPos lexer.Position
	EndPos lexer.Position
	ParentNode TraversableNode
}

type TraversableNodeValue struct {
	Value interface{}
	Type string
	BeginPos lexer.Position
	EndPos lexer.Position
	ParentNode TraversableNode
}

func MakeTraversableNodeValue(parent TraversableNode, value interface{}, typeName string, begin lexer.Position, end lexer.Position) TraversableNode {
	return &TraversableNodeValue{
		Value: value,
		Type: typeName,
		BeginPos: begin,
		EndPos: end,
		ParentNode: parent,
	}
}

func MakeTraversableNodeToken(parent TraversableNode, value string, begin lexer.Position, end lexer.Position) TraversableNode {
	return &TraversableNodeToken{
		Token: value,
		BeginPos: begin,
		EndPos: end,
		ParentNode: parent,
	}
}

func (*TraversableNodeValue) GetChildren() []TraversableNode {
	return []TraversableNode{}
}

func (*TraversableNodeToken) GetChildren() []TraversableNode {
	return []TraversableNode{}
}

func (ast *TraversableNodeValue) OverrideParent(node TraversableNode) {
	ast.ParentNode = node
}

func (ast *TraversableNodeToken) OverrideParent(node TraversableNode) {
	ast.ParentNode = node
}

func (ast *TraversableNodeValue) Parent() TraversableNode {
	return ast.ParentNode
}

func (ast *TraversableNodeToken) Parent() TraversableNode {
	return ast.ParentNode
}

func (ast *TraversableNodeValue) GetNode() interface{} {
	return ast
}

func (ast *TraversableNodeToken) GetNode() interface{} {
	return ast
}

func (ast *TraversableNodeValue) End() lexer.Position {
	return ast.EndPos
}

func (ast *TraversableNodeToken) End() lexer.Position {
	return ast.EndPos
}

func (ast *TraversableNodeValue) Begin() lexer.Position {
	return ast.BeginPos
}

func (ast *TraversableNodeToken) Begin() lexer.Position {
	return ast.BeginPos
}

func (ast *TraversableNodeValue) Print(c *context.ParsingContext) string {
	return fmt.Sprintf("%v", ast.Value)
}

func (ast *TraversableNodeToken) Print(c *context.ParsingContext) string {
	return ast.Token
}

func TraverseAST(node TraversableNode, visitor func(ast TraversableNode)) {
	children := node.GetChildren()
	for _, child := range children {
		visitor(child)
		TraverseAST(child, visitor)
	}
}
