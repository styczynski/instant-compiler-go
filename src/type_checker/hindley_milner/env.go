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
}

type SimpleEnv map[string]*Scheme

func CreateSimpleEnv(env map[string]*Scheme) SimpleEnv {
	for k, v := range env {
		env[k].t = v.t.MapTypes(func (child Type) Type {
			return child.WithContext(CreateBuilinCodeContext(k, env[k]))
		})
	}
	return env
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
		v.Apply(sub) // apply mutates Scheme, so no need to set
	}
	return e
}

func (e SimpleEnv) FreeTypeVar() TypeVarSet {
	var retVal TypeVarSet
	for _, v := range e {
		retVal = v.FreeTypeVar().Union(retVal)
	}
	return retVal
}

func (e SimpleEnv) SchemeOf(name string) (retVal *Scheme, ok bool) { retVal, ok = e[name]; return }
func (e SimpleEnv) Clone() Env {
	retVal := make(SimpleEnv)
	for k, v := range e {
		retVal[k] = v.Clone()
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
	e[name] = s
	fmt.Printf("%s => %v\n", name, s)
	return e
}

func (e SimpleEnv) Remove(name string) Env {
	delete(e, name)
	return e
}
