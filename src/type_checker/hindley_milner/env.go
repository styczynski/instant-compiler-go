package hindley_milner

import (
	"fmt"
)

// An Env is essentially a map of names to schemes
type Env interface {
	Substitutable
	SchemeOf(string) (*Scheme, bool)
	Clone() Env

	Has(string) bool
	Add(string, *Scheme, int) (Env, *Scheme, *Scheme)
	Remove(string) Env
	VarsNames() []string
	IsOverloaded(string) bool
	OverloadedAlternatives(string) []*Scheme
	IsBuiltin(name string) bool
}

type SimpleEnv struct {
	env map[string][]*Scheme
	builtins map[string]func()[]*Scheme
	levels map[string][]int
}

func CreateSimpleEnv(env map[string][]*Scheme) *SimpleEnv {
	builtins := map[string]func()[]*Scheme{}
	newEnv := map[string][]*Scheme{}
	for k, v := range env {
		name := k
		schemes := v
		for _, s := range schemes {
			s.t = s.t.MapTypes(func(child Type) Type {
				return child.WithContext(CreateBuilinCodeContext(name, s))
			})
		}
		if true {
			builtins[name] = func() []*Scheme {
				//fmt.Printf("Get scheme for %s:\n", name)
				ret := []*Scheme{}
				for _, s := range schemes {
					ret = append(ret, s.DeepClone())
				}
				return ret
			}
		} else {
			newEnv[name] = schemes
		}
	}
	return &SimpleEnv{
		env: newEnv,
		builtins: builtins,
		levels: map[string][]int{},
	}
}

func SingleDef(tvs TypeVarSet, t Type) []*Scheme {
	return []*Scheme{ NewScheme(tvs, t) }
}

func PrintEnv(env Env) {
	fmt.Printf("====== Environment ======\n")
	for _, v := range env.VarsNames() {
		scheme, _ := env.SchemeOf(v)
		fmt.Printf("%s => %v\n", v, scheme)
	}
	fmt.Printf("=========================\n")
}

func (e *SimpleEnv) Apply(sub Subs) Substitutable {
	logf("Applying %v to env", sub)
	if sub == nil {
		return e
	}

	for name, v := range e.env {
		if _, ok := e.builtins[name]; ok {
			// Skip builtins
			continue
		}
		v[0].Apply(sub) // apply mutates Scheme, so no need to set
	}
	return e
}

func (e *SimpleEnv) Has(name string) bool {
	_, ok := e.env[name]
	return ok
}

func (e *SimpleEnv) FreeTypeVar() TypeVarSet {
	var retVal TypeVarSet
	for name, v := range e.env {
		if _, ok := e.builtins[name]; ok {
			// Do not return tv for builtins
			continue
		}
		retVal = v[0].FreeTypeVar().Union(retVal)
	}
	return retVal
}

func (e *SimpleEnv) IsBuiltin(name string) bool {
	_, ok := e.builtins[name];
	return ok
}

func (e *SimpleEnv) SchemeOf(name string) (*Scheme, bool) {
	if b, ok := e.builtins[name]; ok {
		return b()[0], true
	}
	retVal, ok := e.env[name]
	if ok {
		return retVal[0], true
	}
	return nil, false
}
func (e *SimpleEnv) Clone() Env {
	//fmt.Printf("CLONE ENV\n")
	retVal := &SimpleEnv{
		env: make(map[string][]*Scheme),
		builtins: make(map[string]func()[]*Scheme),
		levels: map[string][]int{},
	}
	for k, v := range e.env {
		retVal.env[k] = []*Scheme{}
		for _, s := range v {
			retVal.env[k] = append(retVal.env[k], s.DeepClone())
		}
	}
	for k, v := range e.builtins {
		original := v
		name := k
		retVal.builtins[name] = func() []*Scheme {
			ret := []*Scheme{}
			for _, s := range original() {
				ret = append(ret, s.DeepClone())
			}
			return ret
		}
	}
	for k, v := range e.levels {
		newLevels := []int{}
		for _, i := range v {
			newLevels = append(newLevels, i)
		}
		retVal.levels[k] = newLevels
	}
	return retVal
}

func (e *SimpleEnv) VarsNames() []string {
	names := []string{}
	for name, _ := range e.env {
		names = append(names, name)
	}
	return names
}

func (e *SimpleEnv) Add(name string, s *Scheme, blockScopeLevel int) (Env, *Scheme, *Scheme) {
	if _, ok := e.builtins[name]; ok {
		// Do not override builtins
		return e, nil, nil
	}
	//fmt.Printf("Add %s\n", name)
	if oldLevels, ok := e.levels[name]; ok && len(oldLevels)>0 {
		oldLevel := oldLevels[len(oldLevels)-1]
		//fmt.Printf("Oh noes indentifier is redeclared [%s] <%d, %d>\n", name, oldLevel, blockScopeLevel)
		if oldLevel == blockScopeLevel {
			// Do not redefine!
			e.levels[name] = append(e.levels[name], blockScopeLevel)
			//fmt.Printf("Constrin %v ~ %v\n", e.env[name][0], s)
			return e, e.env[name][0], s
		}
	}
	e.env[name] = []*Scheme{ s }
	e.levels[name] = append(e.levels[name], blockScopeLevel)
	//fmt.Printf("ADD %s => %v\n", name, s)
	return e, nil, nil
}

func (e *SimpleEnv) Remove(name string) Env {
	if _, ok := e.builtins[name]; ok {
		// Do not delete builtins
		return e
	}
	fmt.Printf("Remove %s\n", name)
	if len(e.levels[name])>0 {
		e.levels[name] = e.levels[name][:len(e.levels[name])-1]
	}
	delete(e.env, name)
	return e
}

func (e *SimpleEnv) IsOverloaded(name string) bool {
	if b, ok := e.builtins[name]; ok {
		return len(b()) > 1
	}
	return len(e.env[name]) > 1
}

func (e *SimpleEnv) OverloadedAlternatives(name string) []*Scheme {
	if b, ok := e.builtins[name]; ok {
		return b()
	}
	return e.env[name]
}