package hindley_milner

import (
	"fmt"
)

type SignedStruct struct {
	ts      map[string]Type
	name    string
	context CodeContext
}

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

func (t *SignedStruct) IsPrototype() bool {
	return len(t.name) == 0
}

func (t *SignedStruct) CheckIfCanUnionTypes(other interface{}) error {
	if ot, ok := other.(*SignedStruct); ok {
		if t.name != ot.name && !t.IsPrototype() && !ot.IsPrototype() {
			return fmt.Errorf("Two different entities %s and %s are not compatible.", ot.name, t.name)
		}

		for i, _ := range t.ts {
			if _, ok := ot.ts[i]; !ok && !ot.IsPrototype() {
				return fmt.Errorf("Type %s is missing propery %s", ot.name, i)
			}
		}

		for i, _ := range ot.ts {
			if _, ok := t.ts[i]; !ok && !t.IsPrototype() {
				return fmt.Errorf("Type %s is missing propery %s", t.name, i)
			}
		}
	}
	return nil
}

func (t *SignedStruct) FreeTypeVar() TypeVarSet {
	var tvs TypeVarSet
	for _, v := range t.ts {
		tvs = tvs.Union(v.FreeTypeVar())
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
		if ot.name != t.name && !ot.IsPrototype() && !t.IsPrototype() {
			return false
		}

		for i, tv := range t.ts {
			if ov, ok := ot.ts[i]; (!ok && !ot.IsPrototype()) || !tv.Eq(ov) {
				return false
			}
		}

		for i, ov := range ot.ts {
			if tv, ok := t.ts[i]; (!ok && !t.IsPrototype()) || !tv.Eq(ov) {
				return false
			}
		}
		return true
	}
	return false
}

func (t *SignedStruct) Format(f fmt.State, c rune) {
	if t.IsPrototype() {
		f.Write([]byte("*<"))
	} else {
		f.Write([]byte(fmt.Sprintf("%s<", t.name)))
	}
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
			fmt.Fprintf(f, "%s=%v", k, v)
		}
		i++
	}
	f.Write([]byte(">"))
}

func (t *SignedStruct) MapTypes(mapper TypeMapper) Type {
	newSignedStruct := &SignedStruct{
		ts:      map[string]Type{},
		name:    t.name,
		context: t.context,
	}
	for k, v := range t.ts {
		newSignedStruct.ts[k] = v.MapTypes(mapper)
	}
	return mapper(newSignedStruct)
}

func (t *SignedStruct) WithContext(c CodeContext) Type {
	return &SignedStruct{
		ts:      t.ts,
		name:    t.name,
		context: c,
	}
}

func (t *SignedStruct) GetContext() CodeContext {
	return t.context
}

func (t *SignedStruct) String() string { return fmt.Sprintf("%s%v", TypeStringPrefix(t), t) }

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
