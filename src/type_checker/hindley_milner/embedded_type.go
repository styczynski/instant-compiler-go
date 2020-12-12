package hindley_milner

type EmbeddedTypeExpr struct {
	GetType func()*Scheme
}

func (n EmbeddedTypeExpr) Body() Expression {
	return n
}

func (n EmbeddedTypeExpr) Map(parent Expression, mapper ExpressionMapper) Expression {
	return mapper(parent, n)
}

func (n EmbeddedTypeExpr) Visit(parent Expression, mapper ExpressionMapper) {
	mapper(parent, n)
}

func (n EmbeddedTypeExpr) ExpressionType() ExpressionType {
	return E_TYPE
}

func (n EmbeddedTypeExpr) EmbeddedType() *Scheme {
	return n.GetType()
}
