package hindley_milner

import (
	"fmt"

	"github.com/styczynski/latte-compiler/src/logs"
	"github.com/styczynski/latte-compiler/src/parser/context"
)

type UnionableType interface {
	Union(t2 interface{}, context Constraint, infer InferenceBackend, listener IntrospecionListener) (Subs, error)
}

type Constraints []Constraint

func (cs Constraints) LogContext(c *context.ParsingContext) map[string]interface{} {
	return map[string]interface{}{}
}

func (cs Constraints) Apply(sub Subs) Substitutable {

	if sub == nil {
		return cs
	}

	if len(cs) == 0 {
		return cs
	}

	logs.Debug(cs, "Number of constraints: %d", len(cs))
	logs.Debug(cs, "Applying %v to %v", sub, cs)
	for i, c := range cs {
		cs[i] = c.Apply(sub).(Constraint)
	}
	logs.Debug(cs, "Output constraints %v", cs)
	return cs
}

func (cs Constraints) FreeTypeVar() TypeVarSet {
	var retVal TypeVarSet
	for _, v := range cs {
		retVal = v.FreeTypeVar().Union(retVal)
	}
	return retVal
}

func (cs Constraints) Format(state fmt.State, c rune) {
	state.Write([]byte("Constraints["))
	for i, c := range cs {
		if i < len(cs)-1 {
			fmt.Fprintf(state, "%v, ", c)
		} else {
			fmt.Fprintf(state, "%v", c)
		}
	}
	state.Write([]byte{']'})
}

type Types []Type

func (ts Types) Contains(t Type) bool {
	for _, T := range ts {
		if TypeEq(t, T) {
			return true
		}
	}
	return false
}
