package hindley_milner

import "fmt"

type Env interface {
	Substitutable
	SchemeOf(string) (*Scheme, bool)
	Lookup(f Fresher, name string) (Type, error, bool)
	Clone() Env

	Has(string) bool
	AddPrototype(Fresher, string, *Scheme, int) (Env, *Scheme, *Scheme, DeclarationInfo, error)
	Add(Fresher, string, *Scheme, int, bool) (Env, *Scheme, *Scheme, DeclarationInfo, error)
	Remove(string) Env
	VarsNames() []string
	IsOverloaded(string) bool
	OverloadedAlternatives(string) []*Scheme
	IsBuiltin(name string) bool

	RegisterIntrospectionListener(listener IntrospecionListener)
	GetIntrospecionListener() IntrospecionListener
}

type IntrospecionListener interface {
	OnApply(sub Subs)
	OnApplySingle(tv TypeVariable, t Type)
	AddIntrospectionVariable(tv TypeVariable)
	GetIntrospectionVariable(tv TypeVariable) Type
}

type levelInfo struct {
	level      int
	isProt     bool
	hasAnyProt bool
	uid        int
	baseTV     TypeVariable
	baseScheme *Scheme
}

func (i levelInfo) GetUID() int {
	return i.uid
}

func (i levelInfo) GetTV() TypeVariable {
	return i.baseTV
}

type DeclarationInfo interface {
	GetUID() int
	GetTV() TypeVariable
}

type SimpleEnv struct {
	env      map[string][]*Scheme
	builtins map[string]func() []*Scheme
	levels   map[string][]levelInfo
	uid      int

	listener IntrospecionListener
}

func CreateSimpleEnv(env map[string][]*Scheme) *SimpleEnv {
	builtins := map[string]func() []*Scheme{}
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
		env:      newEnv,
		builtins: builtins,
		levels:   map[string][]levelInfo{},
	}
}

func SingleDef(tvs TypeVarSet, t Type) []*Scheme {
	return []*Scheme{NewScheme(tvs, t)}
}

func PrintEnv(env Env) {
	logf("====== Environment ======\n")
	for _, v := range env.VarsNames() {
		scheme, _ := env.SchemeOf(v)
		logf("%s => %v\n", v, scheme)
	}
	logf("=========================\n")
}

func (e *SimpleEnv) RegisterIntrospectionListener(listener IntrospecionListener) {
	e.listener = listener
}

func (e *SimpleEnv) GetIntrospecionListener() IntrospecionListener {
	return e.listener
}

func (e *SimpleEnv) Apply(sub Subs) Substitutable {
	logf("Applying %v to env", sub)

	if sub == nil {
		return e
	}

	e.listener.OnApply(sub)
	for name, v := range e.env {
		if _, ok := e.builtins[name]; ok {
			continue
		}
		v[0].Apply(sub)
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

			continue
		}
		retVal = v[0].FreeTypeVar().Union(retVal)
	}
	return retVal
}

func (e *SimpleEnv) IsBuiltin(name string) bool {
	_, ok := e.builtins[name]
	return ok
}

func (e *SimpleEnv) Lookup(f Fresher, name string) (Type, error, bool) {
	if e.IsBuiltin(name) {

		scheme, _ := e.SchemeOf(name)

		return Instantiate(f, scheme.Clone()), nil, false
	}

	s, ok := e.SchemeOf(name)
	if !ok {
		return nil, fmt.Errorf("Unknwon symbol: %s", name), false
	}

	if oldLevels, ok := e.levels[name]; ok && len(oldLevels) > 0 {
		oldLevel := oldLevels[len(oldLevels)-1]
		if oldLevel.hasAnyProt {
			return oldLevel.baseScheme.Concrete(), nil, true
			//return oldLevel.baseTV, nil, true
		}
	}

	return Instantiate(f, s), nil, false
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

	retVal := &SimpleEnv{
		env:      make(map[string][]*Scheme),
		builtins: make(map[string]func() []*Scheme),
		levels:   map[string][]levelInfo{},
		uid:      e.uid,
		listener: e.listener,
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
		newLevels := []levelInfo{}
		for _, i := range v {
			newLevels = append(newLevels, levelInfo{
				level:      i.level,
				isProt:     i.isProt,
				hasAnyProt: i.hasAnyProt,
				uid:        i.uid,
				baseTV:     i.baseTV,
				baseScheme: i.baseScheme,
			})
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

type envError struct {
	variableName string
	oldDef       CodeContext
	builtin      bool
}

func (err *envError) Error() string {
	return fmt.Sprintf("Variable redefined: %s", err.variableName)
}

func (e *SimpleEnv) Add(f Fresher, name string, s *Scheme, blockScopeLevel int, isRedefinable bool) (Env, *Scheme, *Scheme, DeclarationInfo, error) {
	return e.addVar(f, name, s, blockScopeLevel, false, isRedefinable)
}

func (e *SimpleEnv) AddPrototype(f Fresher, name string, s *Scheme, blockScopeLevel int) (Env, *Scheme, *Scheme, DeclarationInfo, error) {
	return e.addVar(f, name, s, blockScopeLevel, true, false)
}

func (e *SimpleEnv) addVar(f Fresher, name string, s *Scheme, blockScopeLevel int, isPrototype bool, isRedefinable bool) (Env, *Scheme, *Scheme, DeclarationInfo, error) {
	logf("Add %s ==> %v [%d]\n", name, s, blockScopeLevel)
	if _, ok := e.builtins[name]; ok && blockScopeLevel == 0 && !isRedefinable {

		return e, nil, nil, nil, &envError{
			variableName: name,
			oldDef:       s.t.GetContext(),
			builtin:      true,
		}
	}

	if oldLevels, ok := e.levels[name]; ok && len(oldLevels) > 0 && !isRedefinable {
		oldLevel := oldLevels[len(oldLevels)-1]

		if oldLevel.level == blockScopeLevel {

			inf := levelInfo{
				level:      blockScopeLevel,
				isProt:     isPrototype && oldLevel.isProt,
				hasAnyProt: isPrototype || oldLevel.hasAnyProt,
				uid:        e.uid,
				baseTV:     oldLevel.baseTV,
				baseScheme: oldLevel.baseScheme,
			}
			e.uid++
			e.levels[name] = append(e.levels[name], inf)

			//con1 := NewScheme(nil, inf.baseTV)
			con1 := inf.baseScheme
			con2 := s

			logf("  -> ADD MERGE %v ~ %v\n", con1, con2)
			if !oldLevel.isProt && !isPrototype {

				return e, con1, con2, oldLevel, &envError{
					variableName: name,
					oldDef:       e.env[name][0].t.GetContext(),
				}
			}
			return e, con1, con2, inf, nil
		}
	}
	e.env[name] = []*Scheme{s}
	tv := f.Fresh()
	inf := levelInfo{
		level:      blockScopeLevel,
		isProt:     isPrototype,
		hasAnyProt: isPrototype,
		uid:        e.uid,
		baseTV:     tv,
		baseScheme: s,
	}
	e.uid++
	e.levels[name] = append(e.levels[name], inf)
	logf("  -> ADD OVERRIDE %s => %v\n", name, s)
	return e, nil, nil, inf, nil
}

func (e *SimpleEnv) Remove(name string) Env {
	if _, ok := e.builtins[name]; ok {

		return e
	}
	logf("Remove %s\n", name)
	if len(e.levels[name]) > 0 {
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

type IntrospecionSimpleListener struct {
	introspectionVars  []TypeVariable
	introspectionTypes []Type
}

func NewIntrospecionSimpleListener() *IntrospecionSimpleListener {
	return &IntrospecionSimpleListener{
		introspectionVars:  []TypeVariable{},
		introspectionTypes: []Type{},
	}
}

func (e *IntrospecionSimpleListener) OnApply(sub Subs) {
	for introIndex, introTV := range e.introspectionVars {
		if t, ok := sub.Get(introTV); ok {
			e.introspectionTypes[introIndex] = t
		}
	}
}

func (e *IntrospecionSimpleListener) OnApplySingle(tv TypeVariable, t Type) {
	//fmt.Printf("(?) INTRO Set %d (?%v) to %v\n", tv.value, e.introspectionVars, t)
	for introIndex, introTV := range e.introspectionVars {
		if TypeEq(introTV, tv) {
			e.introspectionTypes[introIndex] = t
			//fmt.Printf("INTRO Set %d to %v\n", tv.value, t)
		}
	}
}

func (e *IntrospecionSimpleListener) AddIntrospectionVariable(tv TypeVariable) {
	e.introspectionVars = append(e.introspectionVars, tv)
	e.introspectionTypes = append(e.introspectionTypes, tv)
}

func (e *IntrospecionSimpleListener) GetIntrospectionVariable(tv TypeVariable) Type {
	for ti, tv0 := range e.introspectionVars {
		if TypeEq(tv, tv0) {
			return e.introspectionTypes[ti]
		}
	}
	return nil
}
