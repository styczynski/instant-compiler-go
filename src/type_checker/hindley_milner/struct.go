package hindley_milner

import (
	"fmt"
)

// SignedStruct is a basic record/tuple type. It takes an optional name.
type SignedStruct struct {
	ts map[string]Type
	name string
	context CodeContext
}

// NewSignedStructType creates a new SignedStruct Type
func NewSignedStructType(name string, ts map[string]Type) *SignedStruct {
	return &SignedStruct{
		ts:   ts,
		name: name,
	}
}

func (t *SignedStruct) Apply(subs Subs) Substitutable {
	ts := map[string]Type{}
	for k, v := range t.ts {
		ts[k] = v.Apply(subs).(Type)
	}
	return NewSignedStructType(t.name, ts)
}

func (t *SignedStruct) FreeTypeVar() TypeVarSet {
	var tvs TypeVarSet
	for _, v := range t.ts {
		tvs = v.FreeTypeVar().Union(tvs)
	}
	return tvs
}

func (t *SignedStruct) Name() string {
	return t.name
}

func (t *SignedStruct) Normalize(k, v TypeVarSet) (Type, error) {
	ts := map[string]Type{}
	var err error
	for i, tt := range t.ts {
		if ts[i], err = tt.Normalize(k, v); err != nil {
			return nil, err
		}
	}
	return NewSignedStructType(t.name, ts), nil
}

func (t *SignedStruct) Types() Types {
	ts := BorrowTypes(len(t.ts))
	o := []Type{}
	for _, v := range t.ts {
		o = append(o, v)
	}
	copy(ts, o)
	return ts
}

func (t *SignedStruct) Eq(other Type) bool {
	if ot, ok := other.(*SignedStruct); ok {
		if len(ot.ts) != len(t.ts) {
			return false
		}
		if ot.name != t.name {
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

func (t *SignedStruct) Format(f fmt.State, c rune) {
	f.Write([]byte(fmt.Sprintf("%s<", t.name)))
	f.Write([]byte(TypeStringPrefix(t)))
	count := 0
	for range t.ts {
		count++
	}
	i := 0
	for k, v := range t.ts {
		if i < count-1 {
			fmt.Fprintf(f, "%s=%v, ", k, v)
		} else {
			fmt.Fprintf(f, "%s=%v>", k, v)
		}
		i++
	}
}

func (t *SignedStruct) MapTypes(mapper TypeMapper) Type {
	newSignedStruct := &SignedStruct{
		ts: map[string]Type{},
		name: t.name,
		context: t.context,
	}
	for k, v := range t.ts {
		newSignedStruct.ts[k] = v.MapTypes(mapper)
	}
	return mapper(newSignedStruct)
}

func (t *SignedStruct) WithContext(c CodeContext) Type {
	return &SignedStruct{
		ts:   t.ts,
		name: t.name,
		context: c,
	}
}

func (t *SignedStruct) GetContext() CodeContext {
	return t.context
}

func (t *SignedStruct) String() string { return fmt.Sprintf("%s%v", TypeStringPrefix(t), t) }

// Clone implements Cloner
func (t *SignedStruct) Clone() interface{} {
	retVal := new(SignedStruct)
	ts := map[string]Type{}
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
