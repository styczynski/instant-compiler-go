package hindley_milner

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/generic_ast"
)

type EmbeddedTypeExpr struct {
	GetType func()*Scheme
	Source generic_ast.NormalNode
}

func (n EmbeddedTypeExpr) Begin() lexer.Position {
	return n.Source.Begin()
}

func (n EmbeddedTypeExpr) End() lexer.Position {
	return n.Source.End()
}

func (n EmbeddedTypeExpr) Body() generic_ast.Expression {
	return n
}

func (n EmbeddedTypeExpr) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	return mapper(parent, n, context, false)
}

func (n EmbeddedTypeExpr) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(parent, n, context)
}

func (n EmbeddedTypeExpr) ExpressionType() ExpressionType {
	return E_TYPE
}

func (n EmbeddedTypeExpr) EmbeddedType() *Scheme {
	return n.GetType()
}
