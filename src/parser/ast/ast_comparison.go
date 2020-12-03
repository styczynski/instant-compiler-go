package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
)

type Comparison struct {
	BaseASTNode
	Addition *Addition   `@@`
	Op       string      `[ @( ">" | ">" "=" | "<" | "<" "=" )`
	Next     *Comparison `  @@ ]`
}

func (ast *Comparison) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Comparison) End() lexer.Position {
	return ast.EndPos
}

func (ast *Comparison) GetNode() interface{} {
	return ast
}

func (ast *Comparison) GetChildren() []TraversableNode {
	return []TraversableNode{
		ast.Addition,
		MakeTraversableNodeToken(ast.Op, ast.Pos, ast.EndPos),
		ast.Next,
	}
}

func (ast *Comparison) HasNext() bool {
	return ast.Next != nil
}

func (ast *Comparison) Print(c *context.ParsingContext) string {
	if ast.HasNext() {
		return printBinaryOperation(c, ast, ast.Addition.Print(c), ast.Op, ast.Next.Print(c))
	}
	return ast.Addition.Print(c)
}
