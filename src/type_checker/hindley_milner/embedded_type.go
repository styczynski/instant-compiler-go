package hindley_milner

type EmbeddedTypeExpr struct {
	GetType func()*Scheme
}

func (n EmbeddedTypeExpr) Body() Expression {
	return n
}

func (n EmbeddedTypeExpr) Map(mapper ExpressionMapper) Expression {
	return mapper(n)
}

func (n EmbeddedTypeExpr) Visit(mapper ExpressionMapper) {
	mapper(n)
}

func (n EmbeddedTypeExpr) ExpressionType() ExpressionType {
	return E_TYPE
}

func (n EmbeddedTypeExpr) EmbeddedType() *Scheme {
	return n.GetType()
}
