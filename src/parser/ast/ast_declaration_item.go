package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type DeclarationItem struct {
	BaseASTNode
	Name string `@Ident`
	Initializer *Expression `( "=" @@ )?`
	ParentNode TraversableNode
}

func (ast *DeclarationItem) Parent() TraversableNode {
	return ast.ParentNode
}

func (ast *DeclarationItem) OverrideParent(node TraversableNode) {
	ast.ParentNode = node
}

func (ast *DeclarationItem) Begin() lexer.Position {
	return ast.Pos
}

func (ast *DeclarationItem) End() lexer.Position {
	return ast.EndPos
}

func (ast *DeclarationItem) GetNode() interface{} {
	return ast
}

func (ast *DeclarationItem) GetChildren() []TraversableNode {
	return []TraversableNode{
		MakeTraversableNodeToken(ast, ast.Name, ast.Pos, ast.EndPos),
		ast.Initializer,
	}
}

func (ast *DeclarationItem) HasInitializer() bool {
	return ast.Initializer != nil
}

func (ast *DeclarationItem) Print(c *context.ParsingContext) string {
	if ast.HasInitializer() {
		return printNode(c, ast, "%s = %s", ast.Name, ast.Initializer.Print(c))
	}
	return printNode(c, ast, "%s", ast.Name)
}

/////

func (ast *DeclarationItem) Body() hindley_milner.Expression {
	if !ast.HasInitializer() {
		return hindley_milner.Batch{Exp: []hindley_milner.Expression{}}
	}
	return ast.Initializer
}

func (ast *DeclarationItem) Map(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	return mapper(parent, &DeclarationItem{
		BaseASTNode: ast.BaseASTNode,
		Name:        ast.Name,
		Initializer: mapper(ast, ast.Initializer).(*Expression),
		ParentNode: parent.(TraversableNode),
	})
}

func (ast *DeclarationItem) Visit(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) {
	mapper(ast, ast.Initializer)
	mapper(parent, ast)
}

func (ast *DeclarationItem) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_PROXY
}