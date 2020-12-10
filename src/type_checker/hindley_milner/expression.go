package hindley_milner

import "fmt"

// A Namer is anything that knows its own name
type Namer interface {
	Name() NameGroup
}

type NameGroup struct {
	names []string
	types map[string]*Scheme
	hasTypesMap bool
}

func (g NameGroup) GetNames() []string {
	return g.names
}

func (g NameGroup) HasTypes() bool {
	return g.hasTypesMap
}

func (g NameGroup) GetTypeOf(name string) *Scheme {
	if !g.hasTypesMap {
		return nil
	}
	if v, ok := g.types[name]; ok {
		return v
	}
	return nil
}

func Name(s string) NameGroup {
	return NameGroup{[]string { s }, nil, false}
}

func Names(s []string) NameGroup {
	return NameGroup{s, nil, false}
}

func NamesWithTypes(names []string, types map[string]*Scheme) NameGroup {
	return NameGroup{names, types, true}
}

func NamesWithTypesFromMap(args map[string]*Scheme) NameGroup {
	names := []string{}
	for name, _ := range args {
		names = append(names, name)
	}
	return NameGroup{names, args, true}
}

// A Typer is an Expression node that knows its own Type
type Typer interface {
	Type() Type
}

// An Inferer is an Expression that can infer its own Type given an Env
type Inferer interface {
	Infer(Env, Fresher) (Type, error)
}

type ExpressionMapper = func (e Expression) Expression

type ExpressionType int

const (
	E_VAR ExpressionType = iota
	E_LITERAL
	E_APPLICATION
	E_LAMBDA
	E_FUNCTION
	E_TYPE
	E_BLOCK
	E_OPAQUE_BLOCK
	E_RETURN
	E_LET
	E_LET_RECURSIVE
	E_DECLARATION
	E_FUNCTION_DECLARATION
	E_CUSTOM
	E_PROXY
	E_NONE
)

// An Expression is basically an AST node. In its simplest form, it's lambda calculus
type Expression interface {
	Body() Expression
	Map(mapper ExpressionMapper) Expression
	Visit(mapper ExpressionMapper)
	ExpressionType() ExpressionType
}

type Batch struct {
	Exp []Expression
}

func (b Batch) ExpressionType() ExpressionType {
	return E_BLOCK
}

func (b Batch) Map(mapper ExpressionMapper) Expression {
	mappedExp := []Expression{}
	for _, exp := range b.Exp {
		mappedExp = append(mappedExp, mapper(exp))
	}
	return mapper(Batch{
		Exp: mappedExp,
	})
}

func (b Batch) Visit(mapper ExpressionMapper) {
	for _, exp := range b.Exp {
		mapper(exp)
	}
	mapper(b)
}

func (b Batch) GetContents() Batch {
	return b
}

func (b Batch) IsBlock() bool {
	return true
}

func (b Batch) Expressions() []Expression {
	return b.Exp
}

func (b Batch) Body() Expression {
	panic(fmt.Errorf("Batch Body() method cannot be called."))
}

func IsBatch(exp Expression) bool {
	_, ok := exp.(Batch)
	return ok
}

func FlattenBatch(exp Expression) []Expression {
	if IsBatch(exp) {
		ret := []Expression{}
		for _, e := range exp.(Batch).Expressions() {
			ret = append(ret, FlattenBatch(e)...)
		}
		return ret
	} else {
		return []Expression{ exp }
	}
}

func ApplyBatch(exp Expression, fn func(e Expression) error) error {
	for _, e := range FlattenBatch(exp) {
		err := fn(e)
		if err != nil {
			return err
		}
	}
	return nil
}

// Var is an expression representing a variable
type Var interface {
	Expression
	Namer
	Typer
}

// Literal is an Expression/AST Node representing a literal
type Literal interface {
	Var
}

// Apply is an Expression/AST node that represents a function application
type Apply interface {
	Expression
	Fn() Expression
}

type LetBase interface {
	// let name = def in body
	Expression
	Var() NameGroup
}

// Let is an Expression/AST node that represents the standard let polymorphism found in functional languages
type Let interface {
	LetBase
	Def() Expression
}

// Lambda is an Expression/AST node that represents a function definiton
type Lambda interface {
	Expression
	Args() NameGroup
}

// EmbeddedType is a type directly embedded into the code
type EmbeddedType interface {
	Expression
	Type() *Scheme
}

// Block is an imperative block of code
type Block interface {
	Expression
	GetContents() Batch
}

// Return is an imperative return statement
type Return interface {
	Expression
}

type DefaultTyper interface {
	DefaultType() *Scheme
}

type CustomExpressionEnv struct {
	Env Env
	InferencedType Type
	LookupEnv func(isLiteral bool, name string) error
	GenerateConstraints func(expr Expression) (error, Env, Type, Constraints)
	FreshTypeVariable func() TypeVariable
}

type CustomExpression interface {
	Expression
	GenerateConstraints(context CustomExpressionEnv) (error, Env, Type, Constraints)
}

//type ExpressionWithRequiredType interface {
//	GetRequiredType() *Scheme
//}

type ExpressionWithIdentifiersDeps interface {
	GetIdentifierDeps() []string
}
