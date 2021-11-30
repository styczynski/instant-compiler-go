package hindley_milner

import (
	"fmt"
)

type TypeMapper = func(t Type) Type

func TypeEq(a Type, b Type) bool {
	//a, b = b, a
	//fmt.Printf("EQ %v %v\n", a, b)

	if unionA, ok := a.(*Union); ok {
		for _, v := range unionA.types {
			if !TypeEq(b, v) {
				return false
			}
		}
		return true
	}
	if unionB, ok := b.(*Union); ok {
		for _, v := range unionB.types {
			if a.Eq(v) {
				return true
			}
		}
		return false
	}

	return a.Eq(b)
}

type Type interface {
	Substitutable
	Name() string
	Normalize(TypeVarSet, TypeVarSet) (Type, error)
	Types() Types
	Eq(Type) bool

	MapTypes(TypeMapper) Type
	WithContext(CodeContext) Type
	GetContext() CodeContext

	fmt.Formatter
	fmt.Stringer
}

func TypeStringPrefix(t Type) string {
	if !t.GetContext().IsEmpty() {

	}
	return ""
}

type Substitutable interface {
	Apply(Subs) Substitutable
	FreeTypeVar() TypeVarSet
}

type TypeConst struct {
	value   string
	context CodeContext
}

func (t TypeConst) Name() string                            { return t.value }
func (t TypeConst) Apply(Subs) Substitutable                { return t }
func (t TypeConst) FreeTypeVar() TypeVarSet                 { return nil }
func (t TypeConst) Normalize(k, v TypeVarSet) (Type, error) { return t, nil }
func (t TypeConst) Types() Types                            { return nil }
func (t TypeConst) String() string                          { return fmt.Sprintf("%s%s", TypeStringPrefix(t), t.value) }
func (t TypeConst) Format(s fmt.State, c rune)              { fmt.Fprintf(s, "%s%s", TypeStringPrefix(t), t.value) }
func (t TypeConst) Eq(other Type) bool {
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

type Record struct {
	ts      []Type
	name    string
	context CodeContext
}

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
			if !TypeEq(v, ot.ts[i]) {
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
		ts:      []Type{},
		name:    t.name,
		context: t.context,
	}
	for _, v := range t.ts {
		newRecord.ts = append(newRecord.ts, v.MapTypes(mapper))
	}
	return mapper(newRecord)
}

func (t *Record) WithContext(c CodeContext) Type {
	return &Record{
		ts:      t.ts,
		name:    t.name,
		context: c,
	}
}

func (t *Record) GetContext() CodeContext {
	return t.context
}

func (t *Record) String() string { return fmt.Sprintf("%s%v", TypeStringPrefix(t), t) }

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
