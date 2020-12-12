package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
)

type Arg struct {
	BaseASTNode
	ArgumentType Type `@@`
	Name string `@Ident`
	ParentNode TraversableNode
}

func (ast *Arg) Parent() TraversableNode {
	return ast.ParentNode
}

func (ast *Arg) OverrideParent(node TraversableNode) {
	ast.ParentNode = node
}

func (ast *Arg) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Arg) End() lexer.Position {
	return ast.EndPos
}

func (ast *Arg) GetNode() interface{} {
	return ast
}

func (ast *Arg) GetChildren() []TraversableNode {
	return []TraversableNode{
		&ast.ArgumentType,
		MakeTraversableNodeToken(ast, ast.Name, ast.Pos, ast.EndPos),
	}
}

func (ast *Arg) Print(c *context.ParsingContext) string {
	return printNode(c, ast, "%s %s", ast.ArgumentType.Print(c), ast.Name)
}
