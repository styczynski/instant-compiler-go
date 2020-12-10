package hindley_milner

import (
	"fmt"

	"github.com/pkg/errors"
)

// TypeVariable is a variable that ranges over the types - that is to say it can take any type.
type TypeVariable struct {
	value int16
	context CodeContext
}

func TVar(name int16) TypeVariable {
	return TypeVariable{
		value: name,
	}
}

func (t TypeVariable) Name() string { return string(t.value) }
func (t TypeVariable) Apply(sub Subs) Substitutable {
	if sub == nil {
		return t
	}

	if retVal, ok := sub.Get(t); ok {
		return retVal
	}

	return t
}

func (t TypeVariable) FreeTypeVar() TypeVarSet { tvs := BorrowTypeVarSet(1); tvs[0] = t; return tvs }
func (t TypeVariable) Normalize(k, v TypeVarSet) (Type, error) {
	if i := k.Index(t); i >= 0 {
		return v[i], nil
	}
	return nil, errors.Errorf("Type Variable %v not in signature", t)
}

func (t TypeVariable) Types() Types               { return nil }
func (t TypeVariable) String() string             { return fmt.Sprintf("%s%s", TypeStringPrefix(t), string(t.value)) }
func (t TypeVariable) Format(s fmt.State, c rune) { fmt.Fprintf(s, "%s%d", TypeStringPrefix(t), t.value) }

func (t TypeVariable) Eq(other Type) bool                      {
	if otherV, ok := other.(TypeVariable); ok {
		return otherV.value == t.value
	}
	return false
}


func (t TypeVariable) MapTypes(mapper TypeMapper) Type {
	return mapper(t)
}

func (t TypeVariable) WithContext(c CodeContext) Type {
	return TypeVariable{
		value:   t.value,
		context: c,
	}
}

func (t TypeVariable) GetContext() CodeContext {
	return t.context
}