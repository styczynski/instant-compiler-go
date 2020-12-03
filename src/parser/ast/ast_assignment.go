package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
)

type Assignment struct {
	BaseASTNode
	TargetName string `@Ident`
	Value *Expression `"=" @@ ";"`
}

func (ast *Assignment) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Assignment) End() lexer.Position {
	return ast.EndPos
}

func (ast *Assignment) GetNode() interface{} {
	return ast
}

func (ast *Assignment) GetChildren() []TraversableNode {
	return []TraversableNode{
		MakeTraversableNodeToken(ast.TargetName, ast.Pos, ast.EndPos),
		ast.Value,
	}
}

func (ast *Assignment) Print(c *context.ParsingContext) string {
	return printNode(c, ast, "%s = %s;", ast.TargetName, ast.Value.Print(c))
}
