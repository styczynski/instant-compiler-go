package ast

import (
	"strings"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
)

type FnDef struct {
	BaseASTNode
	ReturnType Type `@@`
	Name string `@Ident`
	Arg []*Arg `"(" (@@ ( "," @@ )*)? ")"`
	Body Block `@@`
}

func (ast *FnDef) Begin() lexer.Position {
	return ast.Pos
}

func (ast *FnDef) End() lexer.Position {
	return ast.EndPos
}

func (ast *FnDef) GetNode() interface{} {
	return ast
}

func (ast *FnDef) GetChildren() []TraversableNode {
	nodes := make([]TraversableNode, len(ast.Arg) + 3)
	nodes = append(nodes, &ast.ReturnType)
	nodes = append(nodes, MakeTraversableNodeToken(ast.Name, ast.Pos, ast.EndPos))

	for _, child := range ast.Arg {
		nodes = append(nodes, child)
	}
	nodes = append(nodes, &ast.Body)

	return nodes
}

func (ast *FnDef) Print(c *context.ParsingContext) string {
	argsList := []string{}
	for _, arg := range ast.Arg {
		argsList = append(argsList, arg.Print(c))
	}

	return printNode(c, ast, "%s %s(%s) %s",
		ast.ReturnType.Print(c),
		ast.Name,
		strings.Join(argsList, ", "),
		ast.Body.Print(c))
}
