package hindley_milner

import "fmt"

// A Namer is anything that knows its own name
type Namer interface {
	Name() string
}

// A Typer is an Expression node that knows its own Type
type Typer interface {
	Type() Type
}

// An Inferer is an Expression that can infer its own Type given an Env
type Inferer interface {
	Infer(Env, Fresher) (Type, error)
}

// An Expression is basically an AST node. In its simplest form, it's lambda calculus
type Expression interface {
	Body() Expression
}

type Batch struct {
	expressions []Expression
}

func (b Batch) Expressions() []Expression {
	return b.expressions
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
	IsLit() bool
}

// Apply is an Expression/AST node that represents a function application
type Apply interface {
	Expression
	Fn() Expression
}

// LetRec is an Expression/AST node that represents a recursive let
type LetRec interface {
	Let
	IsRecursive() bool
}

// Let is an Expression/AST node that represents the standard let polymorphism found in functional languages
type Let interface {
	// let name = def in body
	Expression
	Namer
	Def() Expression
}

// Lambda is an Expression/AST node that represents a function definiton
type Lambda interface {
	Expression
	Namer
	IsLambda() bool
}
