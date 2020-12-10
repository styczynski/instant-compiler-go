package ast

import (
	"strings"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Declaration struct {
	BaseASTNode
	DeclarationType Type `@@`
	Items []*DeclarationItem `( @@ ( "," @@ )* ) ";"`
}

func (ast *Declaration) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Declaration) End() lexer.Position {
	return ast.EndPos
}

func (ast *Declaration) GetNode() interface{} {
	return ast
}

func (ast *Declaration) GetChildren() []TraversableNode {
	nodes := make([]TraversableNode, len(ast.Items)+1)
	nodes = append(nodes, &ast.DeclarationType)
	for _, child := range ast.Items {
		nodes = append(nodes, child)
	}
	return nodes
}

func (ast *Declaration) Print(c *context.ParsingContext) string {
	declarationItemsList := []string{}
	for _, item := range ast.Items {
		declarationItemsList = append(declarationItemsList, item.Print(c))
	}
	return printNode(c, ast, "%s %s", ast.DeclarationType.Print(c), strings.Join(declarationItemsList, ", "))
}

//////

func (ast *Declaration) Body() hindley_milner.Expression {
	return ast
}

func (ast *Declaration) Map(mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	return mapper(&Declaration{
		BaseASTNode: ast.BaseASTNode,
		DeclarationType: ast.DeclarationType,
		Items: ast.Items,
	}).(*Declaration)
}

func (ast *Declaration) Visit(mapper hindley_milner.ExpressionMapper) {
	mapper(ast)
}

func (ast *Declaration) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_DECLARATION
}

func (ast *Declaration) Var() hindley_milner.NameGroup           {
	names := []string{}
	for _, item := range ast.Items {
		names = append(names, item.Name)
	}
	return hindley_milner.Names(names)
}

func (ast *Declaration) Def() hindley_milner.Expression {
	defs := []hindley_milner.Expression{}
	for _, item := range ast.Items {
		defs = append(defs, item)
	}
	return hindley_milner.Batch{
		Exp: defs,
	}
}