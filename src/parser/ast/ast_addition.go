package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
)

type Addition struct {
	BaseASTNode
	Multiplication *Multiplication `@@`
	Op             string          `[ @( "-" | "+" )`
	Next           *Addition       `  @@ ]`
}

func (ast *Addition) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Addition) End() lexer.Position {
	return ast.EndPos
}

func (ast *Addition) GetNode() interface{} {
	return ast
}

func (ast *Addition) GetChildren() []TraversableNode {
	return []TraversableNode{
		ast.Multiplication,
		MakeTraversableNodeToken(ast.Op, ast.Pos, ast.EndPos),
		ast.Next,
	}
}

func (ast *Addition) HasNext() bool {
	return ast.Next != nil
}

func (ast *Addition) Print(c *context.ParsingContext) string {
	if ast.HasNext() {
		return printBinaryOperation(c, ast, ast.Multiplication.Print(c), ast.Op, ast.Next.Print(c))
	}
	return ast.Multiplication.Print(c)
}
