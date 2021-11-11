package ast

import (
	"fmt"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Empty struct {
	generic_ast.BaseASTNode
	content string `@(";")`
}

func (ast *Empty) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_BLOCK
}

func (ast *Empty) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	return mapper(parent, &Empty{
		content: ast.content,
	}, context, true)
}

func (ast *Empty) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(parent, ast, context)
}

func (ast *Empty) GetContents() hindley_milner.Batch {
	return hindley_milner.Batch{}
}

func (ast *Empty) IsBlock() bool {
	return true
}

func (ast *Empty) Expressions() []generic_ast.Expression {
	return []generic_ast.Expression{}
}

func (ast *Empty) Body() generic_ast.Expression {
	panic(fmt.Errorf("Batch Body() method cannot be called."))
}
