package hindley_milner

import (
	"fmt"
)

type TypeMapper = func(t Type) Type

// Type represents all the possible type constructors.
type Type interface {
	Substitutable
	Name() string                                   // Name is the name of the constructor
	Normalize(TypeVarSet, TypeVarSet) (Type, error) // Normalize normalizes all the type variable names in the type
	Types() Types                                   // If the type is made up of smaller types, then it will return them
	Eq(Type) bool                                   // equality operation

	MapTypes(TypeMapper) Type
	WithContext(CodeContext) Type
	GetContext() CodeContext

	fmt.Formatter
	fmt.Stringer
}

func TypeStringPrefix(t Type) string {
	if !t.GetContext().IsEmpty() {
		//return "@@"
	}
	return ""
}

// Substitutable is any type that can have a set of substitutions applied on it, as well as being able to know what its free type variables are
type Substitutable interface {
	Apply(Subs) Substitutable
	FreeTypeVar() TypeVarSet
}

// TypeConst are the default implementation of a constant type. Feel free to implement your own. TypeConsts should be immutable (so no pointer types plz)
type TypeConst struct {
	value string
	context CodeContext
}

func (t TypeConst) Name() string                            { return t.value }
func (t TypeConst) Apply(Subs) Substitutable                { return t }
func (t TypeConst) FreeTypeVar() TypeVarSet                 { return nil }
func (t TypeConst) Normalize(k, v TypeVarSet) (Type, error) { return t, nil }
func (t TypeConst) Types() Types                            { return nil }
func (t TypeConst) String() string                          { return fmt.Sprintf("%s%s", TypeStringPrefix(t), t.value) }
func (t TypeConst) Format(s fmt.State, c rune)              { fmt.Fprintf(s, "%s%s", TypeStringPrefix(t), t.value) }
func (t TypeConst) Eq(other Type) bool                      {
	if otherV, ok := other.(TypeConst); ok {
		return otherV.value == t.value
	}
	return false
}
func (t TypeConst) MapTypes(mapper TypeMapper) Type {
	return mapper(t)
}
func (t TypeConst) WithContext(c CodeContext) Type {
	return TypeConst{
		value:   t.value,
		context: c,
	}
}
func (t TypeConst) GetContext() CodeContext {
	return t.context
}

// Record is a basic record/tuple type. It takes an optional name.
type Record struct {
	ts   []Type
	name string
	context CodeContext
}

// NewRecordType creates a new Record Type
func NewRecordType(name string, ts ...Type) *Record {
	return &Record{
		ts:   ts,
		name: name,
	}
}

func (t *Record) Apply(subs Subs) Substitutable {
	ts := make([]Type, len(t.ts))
	for i, v := range t.ts {
		ts[i] = v.Apply(subs).(Type)
	}
	return NewRecordType(t.name, ts...)
}

func (t *Record) FreeTypeVar() TypeVarSet {
	var tvs TypeVarSet
	for _, v := range t.ts {
		tvs = v.FreeTypeVar().Union(tvs)
	}
	return tvs
}

func (t *Record) Name() string {
	if t.name != "" {
		return t.name
	}
	return t.String()
}

func (t *Record) Normalize(k, v TypeVarSet) (Type, error) {
	ts := make([]Type, len(t.ts))
	var err error
	for i, tt := range t.ts {
		if ts[i], err = tt.Normalize(k, v); err != nil {
			return nil, err
		}
	}
	return NewRecordType(t.name, ts...), nil
}

func (t *Record) Types() Types {
	ts := BorrowTypes(len(t.ts))
	copy(ts, t.ts)
	return ts
}

func (t *Record) Eq(other Type) bool {
	if ot, ok := other.(*Record); ok {
		if len(ot.ts) != len(t.ts) {
			return false
		}
		for i, v := range t.ts {
			if !v.Eq(ot.ts[i]) {
				return false
			}
		}
		return true
	}
	return false
}

func (t *Record) Format(f fmt.State, c rune) {
	f.Write([]byte("("))
	f.Write([]byte(TypeStringPrefix(t)))
	for i, v := range t.ts {
		if i < len(t.ts)-1 {
			fmt.Fprintf(f, "%v, ", v)
		} else {
			fmt.Fprintf(f, "%v)", v)
		}
	}

}

func (t *Record) MapTypes(mapper TypeMapper) Type {
	newRecord := &Record{
		ts:   []Type{},
		name: t.name,
		context: t.context,
	}
	for _, v := range t.ts {
		newRecord.ts = append(newRecord.ts, v.MapTypes(mapper))
	}
	return mapper(newRecord)
}

func (t *Record) WithContext(c CodeContext) Type {
	return &Record{
		ts:   t.ts,
		name: t.name,
		context: c,
	}
}

func (t *Record) GetContext() CodeContext {
	return t.context
}

func (t *Record) String() string { return fmt.Sprintf("%s%v", TypeStringPrefix(t), t) }

// Clone implements Cloner
func (t *Record) Clone() interface{} {
	retVal := new(Record)
	ts := BorrowTypes(len(t.ts))
	for i, tt := range t.ts {
		if c, ok := tt.(Cloner); ok {
			ts[i] = c.Clone().(Type)
		} else {
			ts[i] = tt
		}
	}
	retVal.ts = ts
	retVal.name = t.name

	return retVal
}
