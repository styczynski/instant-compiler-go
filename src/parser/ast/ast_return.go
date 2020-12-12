package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Return struct {
	BaseASTNode
	Expression *Expression `"return" (@@)? ";"`
	ParentNode TraversableNode
}

func (ast *Return) Parent() TraversableNode {
	return ast.ParentNode
}

func (ast *Return) OverrideParent(node TraversableNode) {
	ast.ParentNode = node
}

func (ast *Return) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Return) End() lexer.Position {
	return ast.EndPos
}

func (ast *Return) GetNode() interface{} {
	return ast
}

func (ast *Return) GetChildren() []TraversableNode {
	return []TraversableNode{
		ast.Expression,
	}
}

func (ast *Return) Print(c *context.ParsingContext) string {
	return printNode(c, ast, "return %s;", ast.Expression.Print(c))
}

///

func (ast *Return) Body() hindley_milner.Expression {
	return ast.Expression
}

func (ast *Return) Map(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	return mapper(parent, &Return{
		BaseASTNode: ast.BaseASTNode,
		Expression:  mapper(ast, ast.Expression).(*Expression),
		ParentNode: parent.(TraversableNode),
	})
}

func (ast *Return) Visit(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) {
	mapper(ast, ast.Expression)
	mapper(parent, ast)
}

func (ast *Return) ExpressionType() hindley_milner.ExpressionType { return hindley_milner.E_RETURN }
