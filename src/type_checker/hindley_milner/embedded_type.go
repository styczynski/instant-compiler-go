package hindley_milner

import "github.com/styczynski/latte-compiler/src/generic_ast"

type EmbeddedTypeExpr struct {
	GetType func()*Scheme
}

func (n EmbeddedTypeExpr) Body() generic_ast.Expression {
	return n
}

func (n EmbeddedTypeExpr) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	return mapper(parent, n, context)
}

func (n EmbeddedTypeExpr) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) {
	mapper(parent, n, context)
}

func (n EmbeddedTypeExpr) ExpressionType() ExpressionType {
	return E_TYPE
}

func (n EmbeddedTypeExpr) EmbeddedType() *Scheme {
	return n.GetType()
}
