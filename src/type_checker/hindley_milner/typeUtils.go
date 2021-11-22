package hindley_milner

func TypeHelperAny() *Scheme {
	return NewScheme(TypeVarSet{
		TVar(0),
	}, TVar(0))
}