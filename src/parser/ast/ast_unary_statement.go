package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
)

type UnaryStatement struct {
	BaseASTNode
	TargetName *string `@Ident`
	Operation string `@( "+" "+" | "-" "-" ) ";"`
}

func (ast *UnaryStatement) Begin() lexer.Position {
	return ast.Pos
}

func (ast *UnaryStatement) End() lexer.Position {
	return ast.EndPos
}

func (ast *UnaryStatement) GetNode() interface{} {
	return ast
}

func (ast *UnaryStatement) GetChildren() []TraversableNode {
	return []TraversableNode{
		MakeTraversableNodeToken(*ast.TargetName, ast.Pos, ast.EndPos),
		MakeTraversableNodeToken(ast.Operation, ast.Pos, ast.EndPos),
	}
}

func (ast *UnaryStatement) Print(c *context.ParsingContext) string {
	return printNode(c, ast, "%s%s;", *ast.TargetName, ast.Operation)
}