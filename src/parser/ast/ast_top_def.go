package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
)

type TopDef struct {
	BaseASTNode
	Function *FnDef `@@`
}

func (ast *TopDef) Begin() lexer.Position {
	return ast.Pos
}

func (ast *TopDef) End() lexer.Position {
	return ast.EndPos
}

func (ast *TopDef) GetNode() interface{} {
	return ast
}

func (ast *TopDef) GetChildren() []TraversableNode {
	return []TraversableNode{
		ast.Function,
	}
}

func (ast *TopDef) Print(c *context.ParsingContext) string {
	return printNode(c, ast, "%s", ast.Function.Print(c))
}
