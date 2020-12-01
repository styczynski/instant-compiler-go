package parser

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"
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

type TraversableNode interface {
	GetChildren() []TraversableNode
	GetNode() interface{}
	Begin() lexer.Position
	End() lexer.Position
	Print(c *ParsingContext) string
}

type TraversableNodeToken struct {
	Token string
	BeginPos lexer.Position
	EndPos lexer.Position
}

type TraversableNodeValue struct {
	Value interface{}
	Type string
	BeginPos lexer.Position
	EndPos lexer.Position
}

func MakeTraversableNodeValue(value interface{}, typeName string, begin lexer.Position, end lexer.Position) TraversableNode {
	return &TraversableNodeValue{
		Value: value,
		Type: typeName,
		BeginPos: begin,
		EndPos: end,
	}
}

func MakeTraversableNodeToken(value string, begin lexer.Position, end lexer.Position) TraversableNode {
	return &TraversableNodeToken{
		Token: value,
		BeginPos: begin,
		EndPos: end,
	}
}

func (*TraversableNodeValue) GetChildren() []TraversableNode {
	return []TraversableNode{}
}

func (*TraversableNodeToken) GetChildren() []TraversableNode {
	return []TraversableNode{}
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

func (ast *TraversableNodeValue) Print(c *ParsingContext) string {
	return fmt.Sprintf("%v", ast.Value)
}

func (ast *TraversableNodeToken) Print(c *ParsingContext) string {
	return ast.Token
}