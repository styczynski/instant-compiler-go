package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Return struct {
	BaseASTNode
	Expression *Expression `"return" (@@)? ";"`
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

func (ast *Return) Map(mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	return mapper(&Return{
		BaseASTNode: ast.BaseASTNode,
		Expression:  mapper(ast.Expression).(*Expression),
	})
}

func (ast *Return) Visit(mapper hindley_milner.ExpressionMapper) {
	mapper(ast.Expression)
	mapper(ast)
}

func (ast *Return) ExpressionType() hindley_milner.ExpressionType { return hindley_milner.E_RETURN }
