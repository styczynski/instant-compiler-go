package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
)

type Multiplication struct {
	BaseASTNode
	Unary *Unary          `@@`
	Op    string          `[ @( "/" | "*" )`
	Next  *Multiplication `  @@ ]`
}

func (ast *Multiplication) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Multiplication) End() lexer.Position {
	return ast.EndPos
}

func (ast *Multiplication) GetNode() interface{} {
	return ast
}

func (ast *Multiplication) GetChildren() []TraversableNode {
	return []TraversableNode{
		ast.Unary,
		MakeTraversableNodeToken(ast.Op, ast.Pos, ast.EndPos),
		ast.Next,
	}
}

func (ast *Multiplication) HasNext() bool {
	return ast.Next != nil
}

func (ast *Multiplication) Print(c *context.ParsingContext) string {
	if ast.HasNext() {
		return printBinaryOperation(c, ast, ast.Unary.Print(c), ast.Op, ast.Next.Print(c))
	}
	return ast.Unary.Print(c)
}

