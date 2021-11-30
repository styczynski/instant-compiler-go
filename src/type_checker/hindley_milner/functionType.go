package hindley_milner

import "fmt"

type FunctionType struct {
	a, b    Type
	context CodeContext
}

func NewFnType(ts ...Type) *FunctionType {
	if len(ts) < 2 {
		panic("Expected at least 2 input types")
	}

	retVal := borrowFnType()
	retVal.a = ts[0]

	if len(ts) > 2 {
		retVal.b = NewFnType(ts[1:]...)
	} else {
		retVal.b = ts[1]
	}
	return retVal
}

func (t *FunctionType) Name() string { return "→" }
func (t *FunctionType) Apply(sub Subs) Substitutable {
	a := t.a
	b := t.b

	t.a = t.a.Apply(sub).(Type)
	t.b = t.b.Apply(sub).(Type)

	t.a = CopyContextTo(t.a, a, b)
	t.b = CopyContextTo(t.b, b, a)
	return t
}

func (t *FunctionType) FreeTypeVar() TypeVarSet { return t.a.FreeTypeVar().Union(t.b.FreeTypeVar()) }
func (t *FunctionType) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "%s%v → %v", TypeStringPrefix(t), t.a, t.b)
}
func (t *FunctionType) String() string { return fmt.Sprintf("%s%v", TypeStringPrefix(t), t) }
func (t *FunctionType) Normalize(k, v TypeVarSet) (Type, error) {
	var a, b Type
	var err error
	if a, err = t.a.Normalize(k, v); err != nil {
		return nil, err
	}

	if b, err = t.b.Normalize(k, v); err != nil {
		return nil, err
	}

	return NewFnType(a, b), nil
}
func (t *FunctionType) Types() Types {
	retVal := BorrowTypes(2)
	retVal[0] = t.a
	retVal[1] = t.b
	return retVal
}

func (t *FunctionType) Eq(other Type) bool {
	if ot, ok := other.(*FunctionType); ok {
		return TypeEq(ot.a, t.a) && TypeEq(ot.b, t.b)
	}
	return false
}

func (t *FunctionType) Arg() Type { return t.a }

func (t *FunctionType) Ret(recursive bool) Type {
	if !recursive {
		return t.b
	}

	if fnt, ok := t.b.(*FunctionType); ok {
		return fnt.Ret(recursive)
	}

	return t.b
}

func (t *FunctionType) FlatTypes() Types {
	retVal := BorrowTypes(8)
	retVal = retVal[:0]

	if a, ok := t.a.(*FunctionType); ok {
		ft := a.FlatTypes()
		retVal = append(retVal, ft...)
		ReturnTypes(ft)
	} else {
		retVal = append(retVal, t.a)
	}

	if b, ok := t.b.(*FunctionType); ok {
		ft := b.FlatTypes()
		retVal = append(retVal, ft...)
		ReturnTypes(ft)
	} else {
		retVal = append(retVal, t.b)
	}
	return retVal
}

func (t *FunctionType) Clone() interface{} {
	retVal := new(FunctionType)

	if ac, ok := t.a.(Cloner); ok {
		retVal.a = ac.Clone().(Type)
	} else {
		retVal.a = t.a
	}

	if bc, ok := t.b.(Cloner); ok {
		retVal.b = bc.Clone().(Type)
	} else {
		retVal.b = t.b
	}
	return retVal
}

func (t *FunctionType) MapTypes(mapper TypeMapper) Type {
	return mapper(&FunctionType{
		a:       mapper(t.a),
		b:       mapper(t.b),
		context: t.context,
	})
}

func (t *FunctionType) WithContext(c CodeContext) Type {
	return &FunctionType{
		a:       t.a,
		b:       t.b,
		context: c,
	}
}

func (t *FunctionType) GetContext() CodeContext {
	return t.context
}
