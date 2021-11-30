package hindley_milner

import (
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/styczynski/latte-compiler/src/generic_ast"
)

type ApplicationExpr struct {
	source generic_ast.TraversableNode
	fn     generic_ast.Expression
	arg    generic_ast.Expression
}

func ExpressionApplication(source generic_ast.TraversableNode, fn generic_ast.Expression, arg generic_ast.Expression) *ApplicationExpr {
	return &ApplicationExpr{
		fn:     fn,
		arg:    arg,
		source: source,
	}
}

func (ast *ApplicationExpr) Begin() lexer.Position {
	return ast.source.Begin()
}

func (ast *ApplicationExpr) End() lexer.Position {
	return ast.source.End()
}

func (ast *ApplicationExpr) GetNode() interface{} {
	return ast.source.GetNode()
}

func (ast *ApplicationExpr) GetChildren() []generic_ast.TraversableNode {
	return ast.source.GetChildren()
}

func (ast *ApplicationExpr) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	return mapper(parent, ast, context, false)
}

func (ast *ApplicationExpr) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(parent, ast, context)
}

func (ast *ApplicationExpr) ExpressionType() ExpressionType {
	return E_APPLICATION
}

func (ast *ApplicationExpr) Fn(c InferContext) generic_ast.Expression {
	return ast.fn
}

func (ast *ApplicationExpr) Body() generic_ast.Expression {
	return ast.arg
}
