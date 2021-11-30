package hindley_milner

import (
	"fmt"

	"github.com/styczynski/latte-compiler/src/generic_ast"
)

type Union struct {
	types []Type
	c     CodeContext
}

func NewUnionType(types []Type) *Union {
	return &Union{
		types: types,
	}
}

func (t *Union) FindMatchingFunction(source generic_ast.Expression, args []Type) (error, Type) {
	for _, alternative := range t.types {
		if altFn, ok := alternative.(*FunctionType); !ok {
			allArgs := []Type{}
			allArgs = append(allArgs, args...)
			allArgs = append(allArgs, TVar(0))
			return UnificationWrongTypeError{
				TypeA: t,
				TypeB: NewFnType(allArgs...).WithContext(CreateCodeContext(source)),
				Constraint: Constraint{
					a:       t,
					b:       NewFnType(allArgs...).WithContext(CreateCodeContext(source)),
					context: CreateCodeContext(source),
				},
				Details: fmt.Sprintf("Value is not a function. You can call only functions. Got type: %v", t),
			}, nil
		} else {
			allArgs := []Type{}
			allArgs = append(allArgs, args...)
			allArgs = append(allArgs, altFn.Ret(true))
			if TypeEq(alternative, NewFnType(allArgs...)) {
				return nil, alternative
			}
		}
	}
	return nil, nil
}

func (t *Union) Name() string {
	return "Union"
}

func (t *Union) Apply(subs Subs) Substitutable {
	types := []Type{}
	for _, v := range t.types {
		types = append(types, v.Apply(subs).(Type))
	}
	return NewUnionType(types)
}

func (t *Union) FreeTypeVar() TypeVarSet {
	var tvs TypeVarSet
	for _, v := range t.types {
		tvs = tvs.Union(v.FreeTypeVar())
	}
	return tvs
}

func (t *Union) Normalize(k, v TypeVarSet) (Type, error) {
	types := []Type{}
	for _, tt := range t.types {
		if normalizedType, err := tt.Normalize(k, v); err == nil {
			types = append(types, normalizedType)
		} else {
			return nil, err
		}
	}
	return NewUnionType(types), nil
}

func (t *Union) Types() Types {
	ts := BorrowTypes(len(t.types))
	o := []Type{}
	for _, v := range t.types {
		o = append(o, v)
	}
	copy(ts, o)
	return ts
}

func (t *Union) Eq(other Type) bool {
	if ot, ok := other.(*Union); ok {
		for _, ov := range ot.types {
			hasMatch := false
			for _, v := range t.types {
				if TypeEq(v, ov) {
					hasMatch = true
					break
				}
			}
			if !hasMatch {
				return false
			}
		}
	} else {
		for _, v := range t.types {
			if !TypeEq(v, other) {
				return false
			}
		}
	}
	return true
}

func (t *Union) Format(f fmt.State, c rune) {
	f.Write([]byte("("))
	count := 0
	for range t.types {
		count++
	}
	i := 0
	for _, v := range t.types {
		if i < count-1 {
			fmt.Fprintf(f, "%v | ", v)
		} else {
			fmt.Fprintf(f, "%v", v)
		}
		i++
	}
	f.Write([]byte(")"))
}

func (t *Union) MapTypes(mapper TypeMapper) Type {
	newUnion := &Union{
		types: []Type{},
	}
	for _, v := range t.types {
		newUnion.types = append(newUnion.types, v.MapTypes(mapper))
	}
	return mapper(newUnion)
}

func (t *Union) WithContext(c CodeContext) Type {
	return &Union{
		types: t.types,
		c:     c,
	}
}

func (t *Union) GetContext() CodeContext {
	return t.c
}

func (t *Union) String() string { return fmt.Sprintf("%v", t) }

func (t *Union) Clone() interface{} {
	retVal := new(Union)
	types := []Type{}
	for _, tt := range t.types {
		if c, ok := tt.(Cloner); ok {
			types = append(types, c.Clone().(Type))
		} else {
			types = append(types, tt)
		}
	}
	retVal.types = types
	return retVal
}
