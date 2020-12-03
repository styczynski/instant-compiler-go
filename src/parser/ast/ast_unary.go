package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
)

type Unary struct {
	BaseASTNode
	Op      string   `  ( @( "!" | "-" )`
	Unary   *Unary   `    @@ )`
	UnaryApplication *UnaryApplication `| @@`
}

func (ast *Unary) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Unary) End() lexer.Position {
	return ast.EndPos
}

func (ast *Unary) GetNode() interface{} {
	return ast
}

func (ast *Unary) GetChildren() []TraversableNode {
	if ast.IsOperation() {
		return []TraversableNode{
			MakeTraversableNodeToken(ast.Op, ast.Pos, ast.EndPos),
			ast.Unary,
		}
	} else if ast.IsUnaryApplication() {
		return []TraversableNode{
			ast.UnaryApplication,
		}
	}
	return []TraversableNode{}
}

func (ast *Unary) IsOperation() bool {
	return ast.Unary != nil
}

func (ast *Unary) IsUnaryApplication() bool {
	return ast.UnaryApplication != nil
}

func (ast *Unary) Print(c *context.ParsingContext) string {
	if ast.IsOperation() {
		return printUnaryOperation(c, ast, ast.Op, ast.Unary.Print(c))
	} else if ast.IsUnaryApplication() {
		return ast.UnaryApplication.Print(c)
	}
	return "UNKNOWN"
}
