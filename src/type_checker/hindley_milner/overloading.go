package hindley_milner

import (
	"fmt"
)

func MergeDefinitionsWithOverloads(groups ...*NameGroup) (error, *NameGroup) {
	ret := EmptyNameGroup()
	for _, group := range groups {
		for _, name := range group.names {
			if ret.Has(name) {
				err, mergedDef := FunctionUnionMerge(
					ret.GetTypeOf(name).Concrete(),
					group.GetTypeOf(name).Concrete())
				if err != nil {
					return fmt.Errorf("Cannot overload variable %s: %w", name, err), nil
				}
				ret.RemoveAll(name)
				ret.Add(name, Concreate(mergedDef))
			} else {
				ret.Add(name, group.GetTypeOf(name))
			}
		}
	}

	return nil, ret
}

func FunctionUnionDissolve(a Type) (error, []*FunctionType) {
	ret := []*FunctionType{}
	if a == nil {
		return nil, nil
	}
	if unionA, ok := a.(*Union); ok {
		for _, t := range unionA.types {
			err, s := FunctionUnionDissolve(t)
			if err != nil {
				return err, nil
			}
			ret = append(ret, s...)
		}
	} else if fnA, ok := a.(*FunctionType); ok {
		ret = append(ret, fnA)
	} else {
		return fmt.Errorf("Invalid type found that is nor function nor union type: %v", a), nil
	}
	return nil, ret
}

func FunctionUnionCreate(unnormalizedFns ...*FunctionType) (error, Type) {
	fns := []*FunctionType{}
	for _, fn := range unnormalizedFns {
		found := false
		for _, curFn := range fns {
			if curFn.Eq(fn) {
				found = true
				break
			}
		}
		if !found {
			fns = append(fns, fn)
		}
	}

	if len(fns) == 0 {
		return fmt.Errorf("Empty union type"), nil
	} else if len(fns) == 1 {
		return nil, fns[0]
	} else {
		allTypes := []Type{}
		for _, fn := range fns {
			allTypes = append(allTypes, fn)
		}
		return nil, NewUnionType(allTypes)
	}
}

func FunctionUnionMerge(a Type, b Type) (error, Type) {
	err, aFns := FunctionUnionDissolve(a)
	if err != nil {
		return err, nil
	}
	err, bFns := FunctionUnionDissolve(b)
	if err != nil {
		return err, nil
	}
	aFns = append(aFns, bFns...)
	return FunctionUnionCreate(aFns...)
}
