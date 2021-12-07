package hindley_milner

import (
	"fmt"

	"github.com/styczynski/latte-compiler/src/generic_ast"
)

type Namer interface {
	Name() *NameGroup
}

type NameGroup struct {
	names       []string
	types       map[string]*Scheme
	hasTypesMap bool
}

func (g *NameGroup) RemoveAll(name string) {
	newNames := []string{}
	for _, n := range g.names {
		if n != name {
			newNames = append(newNames, n)
		}
	}
	g.names = newNames
	delete(g.types, name)
}

func (g *NameGroup) Add(name string, t *Scheme) bool {
	if g.HasTypes() {
		g.names = append(g.names, name)
		g.types[name] = t
		return true
	}
	return false
}

func (g *NameGroup) GetNames() []string {
	return g.names
}

func (g *NameGroup) HasTypes() bool {
	return g.hasTypesMap
}

func (g *NameGroup) Has(name string) bool {
	if !g.HasTypes() {
		for _, n := range g.names {
			if n == name {
				return true
			}
		}
		return false
	}
	_, ok := g.types[name]
	return ok
}

func (g *NameGroup) GetTypeOf(name string) *Scheme {
	if !g.hasTypesMap {
		return nil
	}
	if v, ok := g.types[name]; ok {
		return v
	}
	return nil
}

func (g *NameGroup) Format(state fmt.State, c rune) {
	state.Write([]byte("{"))
	for i, name := range g.names {
		if g.HasTypes() {
			if i < len(g.names)-1 {
				fmt.Fprintf(state, "%s => %v, ", name, g.types[name])
			} else {
				fmt.Fprintf(state, "%s => %v", name, g.types[name])
			}
		} else {
			if i < len(g.names)-1 {
				fmt.Fprintf(state, "%s, ", name)
			} else {
				fmt.Fprintf(state, "%s", name)
			}
		}
	}
	state.Write([]byte{'}'})
}

func EmptyNameGroup() *NameGroup {
	return NamesWithTypesFromMap(make(map[string]*Scheme))
}

func Name(s string) *NameGroup {
	return &NameGroup{[]string{s}, nil, false}
}

func NameWithType(s string, t *Scheme) *NameGroup {
	return &NameGroup{[]string{s}, map[string]*Scheme{s: t}, true}
}

func Names(s []string) *NameGroup {
	return &NameGroup{s, nil, false}
}

func NamesWithTypes(names []string, types map[string]*Scheme) *NameGroup {
	return &NameGroup{names, types, true}
}

func NamesWithTypesFromMap(args map[string]*Scheme) *NameGroup {
	names := []string{}
	for name, _ := range args {
		names = append(names, name)
	}
	return &NameGroup{names, args, true}
}

type Typer interface {
	Type() Type
}

type Inferer interface {
	Infer(Env, Fresher) (Type, error)
}

type InferContext interface {
	TypeOf(et generic_ast.Expression, contextExpressions ...generic_ast.Expression) (Type, error)
}

type ExpressionType int

const (
	E_VAR ExpressionType = iota
	E_LITERAL
	E_APPLICATION
	E_TYPE_EQUALITY
	E_LAMBDA
	E_FUNCTION
	E_TYPE
	E_BLOCK
	E_OPAQUE_BLOCK
	E_RETURN
	E_LET
	E_LET_RECURSIVE
	E_REDEFINABLE_LET
	E_DECLARATION
	E_FUNCTION_DECLARATION
	E_CUSTOM
	E_PROXY
	E_NONE
	E_INTROSPECTION
)

type HMExpression interface {
	generic_ast.Expression
	ExpressionType() ExpressionType
}

type HMExpressionWithCustomMismatchErrorDescription interface {
	OnTypeMismatch(generic_ast.NodeWithPosition, generic_ast.NodeWithPosition) []string
}

type Batch struct {
	Exp []generic_ast.Expression
}

func (b Batch) ExpressionType() ExpressionType {
	return E_BLOCK
}

func (b Batch) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	mappedExp := []generic_ast.Expression{}
	for _, exp := range b.Exp {
		mappedExp = append(mappedExp, mapper(b, exp, context, false))
	}
	return mapper(parent, Batch{
		Exp: mappedExp,
	}, context, true)
}

func (b Batch) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	for _, exp := range b.Exp {
		mapper(b, exp, context)
	}
	mapper(parent, b, context)
}

func (b Batch) GetContents() Batch {
	return b
}

func (b Batch) IsBlock() bool {
	return true
}

func (b Batch) Expressions() []generic_ast.Expression {
	return b.Exp
}

func (b Batch) Body() generic_ast.Expression {
	panic(fmt.Errorf("Batch Body() method cannot be called."))
}

func IsBatch(exp generic_ast.Expression) bool {
	_, ok := exp.(Batch)
	return ok
}

func FlattenBatch(exp generic_ast.Expression) []generic_ast.Expression {
	if IsBatch(exp) {
		ret := []generic_ast.Expression{}
		for _, e := range exp.(Batch).Expressions() {
			ret = append(ret, FlattenBatch(e)...)
		}
		return ret
	} else {
		return []generic_ast.Expression{exp}
	}
}

func ApplyBatch(exp generic_ast.Expression, fn func(e generic_ast.Expression) error) error {
	for _, e := range FlattenBatch(exp) {
		err := fn(e)
		if err != nil {
			return err
		}
	}
	return nil
}

type Var interface {
	generic_ast.Expression
	Namer
	Typer
}

type Literal interface {
	Var
}

type Apply interface {
	generic_ast.Expression
	Fn(c InferContext) generic_ast.Expression
}

type LetBase interface {
	generic_ast.Expression
	Var(c InferContext) *NameGroup
}

type Let interface {
	LetBase
	Def(c InferContext) generic_ast.Expression
}

type Lambda interface {
	generic_ast.Expression
	Args(c InferContext) *NameGroup
}

type EmbeddedType interface {
	generic_ast.Expression
	EmbeddedType(c InferContext) *Scheme
}

type Block interface {
	generic_ast.Expression
	GetContents() Batch
}

type Return interface {
	generic_ast.Expression
	HasValidReturnType(returnType Type) bool
}

type DefaultTyper interface {
	DefaultType(c InferContext) *Scheme
}

type CustomExpressionEnv struct {
	Env                 Env
	InferencedType      Type
	LookupEnv           func(isLiteral bool, name string) error
	GenerateConstraints func(expr generic_ast.Expression) (error, Env, Type, Constraints)
	FreshTypeVariable   func() TypeVariable
}

type CustomExpression interface {
	generic_ast.Expression
	GenerateConstraints(context CustomExpressionEnv) (error, Env, Type, Constraints)
}

type ExpressionWithIdentifiersDeps interface {
	GetIdentifierDeps(c InferContext, pre bool) (error, *NameGroup)
}

type IntrospectionExpression interface {
	generic_ast.Expression
	OnTypeReturned(t Type)
}
