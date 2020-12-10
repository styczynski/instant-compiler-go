package ast

import (
	"fmt"

	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Primitive int

const (
	T_BOOL Primitive = iota
	T_INT
	T_VOID
	T_STRING
)

func CreatePrimitive(p Primitive) PrimitiveType {
	name := "void"
	if p == T_STRING {
		name = "string"
	} else if p == T_BOOL {
		name = "bool"
	} else if p == T_INT {
		name = "int"
	} else if p == T_VOID {
		name = "void"
	}
	return PrimitiveType{
		name:    name,
	}
}

type PrimitiveType struct {
	name string
	context hindley_milner.CodeContext
}

func (t PrimitiveType) Name() string                                   { return t.name }
func (t PrimitiveType) Apply(hindley_milner.Subs) hindley_milner.Substitutable                       { return t }
func (t PrimitiveType) FreeTypeVar() hindley_milner.TypeVarSet                        { return nil }
func (t PrimitiveType) Normalize(hindley_milner.TypeVarSet, hindley_milner.TypeVarSet) (hindley_milner.Type, error) { return t, nil }
func (t PrimitiveType) Types() hindley_milner.Types                                   { return nil }
func (t PrimitiveType) Eq(other hindley_milner.Type) bool {
	if ot, ok := other.(PrimitiveType); ok {
		return ot.name == t.name
	}
	return false
}

func (t PrimitiveType) Format(s fmt.State, c rune) { fmt.Fprintf(s, "%s", t.name) }
func (t PrimitiveType) String() string {
	return fmt.Sprintf("%s", t.name)
}

func (t PrimitiveType) MapTypes(mapper hindley_milner.TypeMapper) hindley_milner.Type {
	return mapper(t)
}

func (t PrimitiveType) WithContext(c hindley_milner.CodeContext) hindley_milner.Type {
	return PrimitiveType{
		name: t.name,
		context: c,
	}
}

func (t PrimitiveType) GetContext() hindley_milner.CodeContext {
	return t.context
}
