package errors

import "github.com/styczynski/latte-compiler/src/generic_ast"

type LatteError interface {
	Error() string
	CliMessage() string
	ErrorName() string
}

type LocalizedError struct {
	message string
	errorName string
	source generic_ast.NormalNode
}

func CreateLocalizedError(errorName string, message string, source generic_ast.NormalNode) *LocalizedError {
	return &LocalizedError{
		message: message,
		errorName: errorName,
		source: source,
	}
}

func (le *LocalizedError) Error() string {
	return le.message
}

func (le *LocalizedError) ErrorName() string {
	return le.errorName
}

func (le *LocalizedError) Source() generic_ast.NormalNode {
	return le.source
}