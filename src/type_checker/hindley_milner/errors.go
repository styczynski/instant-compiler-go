package hindley_milner

import (
	"fmt"
	"strings"

	"github.com/styczynski/latte-compiler/src/generic_ast"
)

type UnificationLengthError struct {
	TypeA      Type
	TypeB      Type
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

func (err UnificationLengthError) Source() generic_ast.Expression {
	return *(err.Constraint.context.Source)
}

type UnificationWrongTypeError struct {
	TypeA      Type
	TypeB      Type
	Constraint Constraint
	Details    string
}

func (err UnificationWrongTypeError) IsCausedByBuiltin() bool {
	return err.Constraint.context.IsBuiltin()
}

func (err UnificationWrongTypeError) GetCauseName() string {
	return err.Constraint.context.Name
}

func (err UnificationWrongTypeError) HasSource() bool {
	if err.Constraint.context.Source == nil {
		return false
	}
	return true
}

func (err UnificationWrongTypeError) Source() generic_ast.Expression {
	if err.Constraint.context.Source == nil {
		//logs.Debug("LOLZ: %v %v %v %v\n", err.Constraint.a.GetContext().String(), err.Constraint.b.GetContext().String(), err.Constraint.context.String())
	}
	return *(err.Constraint.context.Source)
}

func (err UnificationWrongTypeError) Error() string {
	details := ""
	if len(err.Details) > 0 {
		details = fmt.Sprintf("\n    Details:\n      %s\n   ", err.Details)
	}
	return fmt.Sprintf("Failed to unify types %s and %s. Mismatched types.%s",
		err.TypeA.String(),
		err.TypeB.String(),
		details,
	)
}

type UnificationRecurrentTypeError struct {
	Type               Type
	Variable           TypeVariable
	VariableTypeSource Type
	Constraint         Constraint
}

func (err UnificationRecurrentTypeError) Source() generic_ast.Expression {
	return *(err.Constraint.context.Source)
}

func (err UnificationRecurrentTypeError) Error() string {
	return fmt.Sprintf("Failed to bind type variable %s from type %s to type %s. The type is recurent.",
		err.VariableTypeSource.String(),
		err.Variable.String(),
		err.Type.String(),
	)
}

type UndefinedSymbol struct {
	Name       string
	Source     generic_ast.Expression
	IsLiteral  bool
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

type InvalidOverloadCandidatesError struct {
	Name       string
	Candidates []*Scheme
	Context    CodeContext
}

func (err InvalidOverloadCandidatesError) Error() string {
	candidatesDescriptions := []string{}
	for i, cand := range err.Candidates {
		candidatesDescriptions = append(candidatesDescriptions, fmt.Sprintf("    %d: %v", i+1, cand))
	}
	return fmt.Sprintf("Failed to find matching definition for %s among all overloaded candidates:\n%s",
		err.Name,
		strings.Join(candidatesDescriptions, "\n"),
	)
}

func (err InvalidOverloadCandidatesError) Source() generic_ast.Expression {
	return *(err.Context.Source)
}

type VariableRedefinedError struct {
	Name               string
	PreviousDefinition CodeContext
	Context            CodeContext
}

func (err VariableRedefinedError) Error() string {
	return fmt.Sprintf("Variable %s has second definition in the current scope.",
		err.Name,
	)
}

func (err VariableRedefinedError) Source() generic_ast.Expression {
	return *(err.Context.Source)
}

type BuiltinRedefinedError struct {
	Name    string
	Context CodeContext
}

func (err BuiltinRedefinedError) Error() string {
	return fmt.Sprintf("%s is a builin variable, you cannot override it. Please use different name.",
		err.Name,
	)
}

func (err BuiltinRedefinedError) Source() generic_ast.Expression {
	return *(err.Context.Source)
}
