package compiler

// Location of the variable
type Location = int64

// Variable metadata
type Variable = int64

const NOT_EXISTING = int64(0)

type CompilerState struct {
	store        map[Variable]Location
	scope        map[string]Variable
	freeLocation Location
	freeVar      Variable

	varNameGenerator UniqueNameGenerator
}

func (s *CompilerState) nextFreeVar() Variable {
	ret := s.freeVar
	s.freeVar++
	return ret
}

func (s *CompilerState) nextFreeLocation() Location {
	ret := s.freeLocation
	s.freeLocation++
	return ret
}

func (s *CompilerState) Allocate(v Variable) Location {
	loc := s.nextFreeLocation()
	s.store[v] = loc
	return loc
}

func (s *CompilerState) AllocateAt(v Variable, loc Location) Location {
	s.store[v] = loc
	return loc
}

func (s *CompilerState) Define(name string) Variable {
	v := s.nextFreeVar()
	s.scope[name] = v
	return v
}

func (s *CompilerState) DefineAndAlloc(name string) (Variable, Location) {
	v := s.Define(name)
	loc := s.Allocate(v)
	return v, loc
}

func (s *CompilerState) NextUniqueVariableName() string {
	return s.varNameGenerator.Next()
}

func (s *CompilerState) GetVariableLocation(v Variable) Location {
	if loc, ok := s.store[v]; ok {
		return loc
	}
	return NOT_EXISTING
}

func (s *CompilerState) GetVariableFromScope(name string) Location {
	if v, ok := s.scope[name]; ok {
		return v
	}
	return NOT_EXISTING
}

func (s *CompilerState) ScopeSize() int {
	return len(s.scope)
}

func (s *CompilerState) GetLocationFromScope(name string) Location {
	v := s.GetVariableFromScope(name)
	if v == NOT_EXISTING {
		return NOT_EXISTING
	}
	return s.GetVariableLocation(v)
}

func (s *CompilerState) Copy() *CompilerState {
	ret := &CompilerState{
		store:            map[Variable]Location{},
		scope:            map[string]Variable{},
		freeLocation:     s.freeLocation,
		freeVar:          s.freeVar,
		varNameGenerator: s.varNameGenerator.Copy(),
	}

	for k, v := range s.store {
		ret.store[k] = v
	}

	for k, v := range s.scope {
		ret.scope[k] = v
	}

	return ret
}

func CreateCompilerState() *CompilerState {
	return &CompilerState{
		store:        map[Variable]Location{},
		scope:        map[string]Variable{},
		freeLocation: 1,
		freeVar:      1,
		varNameGenerator: &SeqNameGenerator{
			nameID: 1,
		},
	}
}
