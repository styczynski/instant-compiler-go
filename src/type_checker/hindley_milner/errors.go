package hindley_milner

import (
	"fmt"
)

type UnificationLengthError struct {
	TypeA Type
	TypeB Type
	Constraint Constraint
}

func (err UnificationLengthError) Error() string {
	return fmt.Sprintf("Failed to unify types %s and %s. They have different lengths.",
		err.TypeA.String(),
		err.TypeB.String())
}

func (err UnificationLengthError) IsCausedByBuiltin() bool {
	return err.Constraint.context.IsBuiltin()
}

func (err UnificationLengthError) GetCauseName() string {
	return err.Constraint.context.Name
}

func (err UnificationLengthError) Source() Expression {
	return *(err.Constraint.context.Source)
}

type UnificationWrongTypeError struct {
	TypeA Type
	TypeB Type
	Constraint Constraint
}

func (err UnificationWrongTypeError) IsCausedByBuiltin() bool {
	return err.Constraint.context.IsBuiltin()
}

func (err UnificationWrongTypeError) GetCauseName() string {
	return err.Constraint.context.Name
}

func (err UnificationWrongTypeError) Source() Expression {
	return *(err.Constraint.context.Source)
}

func (err UnificationWrongTypeError) Error() string {
	return fmt.Sprintf("Failed to unify types %s and %s. Mismatched types.",
		err.TypeA.String(),
		err.TypeB.String(),
	)
}

type UnificationRecurrentTypeError struct {
	Type Type
	Variable TypeVariable
	VariableTypeSource Type
	Constraint Constraint
}

func (err UnificationRecurrentTypeError) Error() string {
	return fmt.Sprintf("Failed to bind type variable %s from type %s to type %s. The type is recurent.",
		err.VariableTypeSource.String(),
		err.Variable.String(),
		err.Type.String(),
	)
}

type UndefinedSymbol struct {
	Name string
	Source Expression
	IsLiteral bool
	IsVariable bool
}

func (err UndefinedSymbol) Error() string {
	name := "symbol"
	if err.IsVariable {
		name = "variable"
	} else if err.IsLiteral {
		name = "literal"
	}
	return fmt.Sprintf("Unknown %s was used: \"%s\"",
		name,
		err.Name,
	)
}
