package ast

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
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

func (ast *LatteProgram) Body() hindley_milner.Expression {
	panic(fmt.Errorf("Batch Body() method cannot be called."))
}

/////

func (ast *LatteProgram) Map(mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	mappedDef := []*TopDef{}
	for _, def := range ast.Definitions {
		mappedDef = append(mappedDef, mapper(def).(*TopDef))
	}
	return mapper(&LatteProgram{
		BaseASTNode: ast.BaseASTNode,
		Definitions: mappedDef,
	}).(*LatteProgram)
}

func (ast *LatteProgram) Visit(mapper hindley_milner.ExpressionMapper) {
	for _, def := range ast.Definitions {
		mapper(def)
	}
	mapper(ast)
}

func (ast *LatteProgram) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_OPAQUE_BLOCK
}

func (ast *LatteProgram) GetContents() hindley_milner.Batch {
	exp := []hindley_milner.Expression{}
	for _, def := range ast.Definitions {
		exp = append(exp, def)
	}
	return hindley_milner.Batch{
		Exp: exp,
	}
}