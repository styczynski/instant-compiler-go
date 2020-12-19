package hindley_milner

import "fmt"

type Constraint struct {
	a, b Type
	context CodeContext
}

func (c Constraint) Context() CodeContext {
	return c.context
}

func (c Constraint) Apply(sub Subs) Substitutable {
	c.a = c.a.Apply(sub).(Type)
	c.b = c.b.Apply(sub).(Type)
	return c
}

func (c Constraint) FreeTypeVar() TypeVarSet {
	var retVal TypeVarSet
	retVal = c.a.FreeTypeVar().Union(retVal)
	retVal = c.b.FreeTypeVar().Union(retVal)
	return retVal
}

func (c Constraint) Format(state fmt.State, r rune) {
	fmt.Fprintf(state, "{%v = %v}", c.a, c.b)
}
