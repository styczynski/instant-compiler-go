package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
)

type ClassField struct {
	BaseASTNode
	ClassFieldType Type `@@`
	Name string `@Ident`
	ParentNode TraversableNode
}

func (ast *ClassField) Parent() TraversableNode {
	return ast.ParentNode
}

func (ast *ClassField) OverrideParent(node TraversableNode) {
	ast.ParentNode = node
}

func (ast *ClassField) Begin() lexer.Position {
	return ast.Pos
}

func (ast *ClassField) End() lexer.Position {
	return ast.EndPos
}

func (ast *ClassField) GetNode() interface{} {
	return ast
}

func (ast *ClassField) GetChildren() []TraversableNode {
	return []TraversableNode{
		&ast.ClassFieldType,
		MakeTraversableNodeToken(ast, ast.Name, ast.Pos, ast.EndPos),
	}
}

func (ast *ClassField) Print(c *context.ParsingContext) string {
	return printNode(c, ast, "%s %s", ast.ClassFieldType.Print(c), ast.Name)
}
