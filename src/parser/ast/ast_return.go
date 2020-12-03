package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
)

type Return struct {
	BaseASTNode
	Expression *Expression `"return" (@@)? ";"`
}

func (ast *Return) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Return) End() lexer.Position {
	return ast.EndPos
}

func (ast *Return) GetNode() interface{} {
	return ast
}

func (ast *Return) GetChildren() []TraversableNode {
	return []TraversableNode{
		ast.Expression,
	}
}

func (ast *Return) Print(c *context.ParsingContext) string {
	return printNode(c, ast, "return %s;", ast.Expression.Print(c))
}