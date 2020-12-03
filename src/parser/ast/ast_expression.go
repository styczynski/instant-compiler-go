package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
)

type Expression struct {
	ComplexASTNode
	LogicalOperation *LogicalOperation `@@`
}

func (ast *Expression) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Expression) End() lexer.Position {
	return ast.EndPos
}

func (ast *Expression) GetNode() interface{} {
	return ast
}

func (ast *Expression) GetChildren() []TraversableNode {
	return []TraversableNode{ ast.LogicalOperation, }
}

func printBinaryOperation(c *context.ParsingContext, ast TraversableNode, arg1 string, operator string, arg2 string) string{
	return printNode(c, ast, "%s %s %s", arg1, operator, arg2)
}

func printUnaryOperation(c *context.ParsingContext, ast TraversableNode, operator string, arg string) string{
	return printNode(c, ast, "%s%s", operator, arg)
}

func (ast *Expression) Print(c *context.ParsingContext) string {
	return ast.LogicalOperation.Print(c)
}

