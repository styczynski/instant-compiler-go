package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
)

type ClassField struct {
	generic_ast.BaseASTNode
	ClassFieldType Type `@@`
	Name string `@Ident`
	ParentNode generic_ast.TraversableNode
}

func (ast *ClassField) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *ClassField) OverrideParent(node generic_ast.TraversableNode) {
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

func (ast *ClassField) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{
		&ast.ClassFieldType,
		generic_ast.MakeTraversableNodeToken(ast, ast.Name, ast.Pos, ast.EndPos),
	}
}

func (ast *ClassField) Print(c *context.ParsingContext) string {
	return printNode(c, ast, "%s %s", ast.ClassFieldType.Print(c), ast.Name)
}
