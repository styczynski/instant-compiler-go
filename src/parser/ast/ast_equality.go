package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
)

type Equality struct {
	BaseASTNode
	Comparison *Comparison `@@`
	Op         string      `[ @( "!" "=" | "=" "=" )`
	Next       *Equality   `  @@ ]`
}

func (ast *Equality) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Equality) End() lexer.Position {
	return ast.EndPos
}

func (ast *Equality) GetNode() interface{} {
	return ast
}

func (ast *Equality) GetChildren() []TraversableNode {
	return []TraversableNode{
		ast.Comparison,
		MakeTraversableNodeToken(ast.Op, ast.Pos, ast.EndPos),
		ast.Next,
	}
}

func (ast *Equality) HasNext() bool {
	return ast.Next != nil
}

func (ast *Equality) Print(c *context.ParsingContext) string {
	if ast.HasNext() {
		return printBinaryOperation(c, ast, ast.Comparison.Print(c), ast.Op, ast.Next.Print(c))
	}
	return ast.Comparison.Print(c)
}
