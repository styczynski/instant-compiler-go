package ast

import (
	"strings"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
)

type LatteProgram struct {
	BaseASTNode
	Definitions []*TopDef `@@*`
}

func (ast *LatteProgram) Begin() lexer.Position {
	return ast.Pos
}

func (ast *LatteProgram) End() lexer.Position {
	return ast.EndPos
}

func (ast *LatteProgram) GetNode() interface{} {
	return ast
}

func (ast *LatteProgram) GetChildren() []TraversableNode {
	nodes := make([]TraversableNode, len(ast.Definitions))
	for _, child := range ast.Definitions {
		nodes = append(nodes, child)
	}
	return nodes
}

func (ast *LatteProgram) Print(c *context.ParsingContext) string {
	defs := []string{}
	for _, def := range ast.Definitions {
		defs = append(defs, def.Print(c))
	}
	return printNode(c, ast, "%s\n", strings.Join(defs, "\n\n"))
}