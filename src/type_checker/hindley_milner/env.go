package hindley_milner

import "fmt"

// An Env is essentially a map of names to schemes
type Env interface {
	Substitutable
	SchemeOf(string) (*Scheme, bool)
	Clone() Env

	Add(string, *Scheme) Env
	Remove(string) Env
	VarsNames() []string
	IsOverloaded(string) bool
	OverloadedAlternatives(string) []*Scheme
}

type SimpleEnv map[string][]*Scheme

func CreateSimpleEnv(env map[string][]*Scheme) SimpleEnv {
	for k, v := range env {
		for _, s := range v {
			s.t = s.t.MapTypes(func(child Type) Type {
				return child.WithContext(CreateBuilinCodeContext(k, s))
			})
		}
	}
	return env
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

func (e SimpleEnv) Apply(sub Subs) Substitutable {
	logf("Applying %v to env", sub)
	if sub == nil {
		return e
	}

	for _, v := range e {
		v[0].Apply(sub) // apply mutates Scheme, so no need to set
	}
	return e
}

func (e SimpleEnv) FreeTypeVar() TypeVarSet {
	var retVal TypeVarSet
	for _, v := range e {
		retVal = v[0].FreeTypeVar().Union(retVal)
	}
	return retVal
}

func (e SimpleEnv) SchemeOf(name string) (*Scheme, bool) {
	retVal, ok := e[name]
	if ok {
		return retVal[0], true
	}
	return nil, false
}
func (e SimpleEnv) Clone() Env {
	retVal := make(SimpleEnv)
	for k, v := range e {
		retVal[k] = []*Scheme{}
		for _, s := range v {
			retVal[k] = append(retVal[k], s.Clone())
		}
	}
	return retVal
}

func (e SimpleEnv) VarsNames() []string {
	names := []string{}
	for name, _ := range e {
		names = append(names, name)
	}
	return names
}

func (e SimpleEnv) Add(name string, s *Scheme) Env {
	e[name] = []*Scheme{ s }
	//fmt.Printf("ADD %s => %v\n", name, s)
	return e
}

func (e SimpleEnv) Remove(name string) Env {
	delete(e, name)
	return e
}

func (e SimpleEnv) IsOverloaded(name string) bool {
	return len(e[name]) > 1
}

func (e SimpleEnv) OverloadedAlternatives(name string) []*Scheme {
	return e[name]
}