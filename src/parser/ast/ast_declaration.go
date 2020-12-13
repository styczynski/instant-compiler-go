package ast

import (
	"strings"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Declaration struct {
	generic_ast.BaseASTNode
	DeclarationType Type `@@`
	Items []*DeclarationItem `( @@ ( "," @@ )* ) ";"`
	ParentNode generic_ast.TraversableNode
}

func (ast *Declaration) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *Declaration) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
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

func (ast *Declaration) GetChildren() []generic_ast.TraversableNode {
	nodes := make([]generic_ast.TraversableNode, len(ast.Items)+1)
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

func (ast *Declaration) Map(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	return mapper(parent, &Declaration{
		BaseASTNode: ast.BaseASTNode,
		DeclarationType: ast.DeclarationType,
		Items: ast.Items,
		ParentNode: parent.(generic_ast.TraversableNode),
	}).(*Declaration)
}

func (ast *Declaration) Visit(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) {
	mapper(parent, ast)
}

func (ast *Declaration) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_DECLARATION
}

func (ast *Declaration) Var() hindley_milner.NameGroup {
	names := []string{}
	types := map[string]*hindley_milner.Scheme{}
	for _, item := range ast.Items {
		names = append(names, item.Name)
		types[item.Name] = ast.DeclarationType.GetType()
	}
	return hindley_milner.NamesWithTypes(names, types)
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