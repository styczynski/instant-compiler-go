package error_collector

import (
	"github.com/styczynski/latte-compiler/src/compiler"
	"github.com/styczynski/latte-compiler/src/parser"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker"
)

type InternalError interface {
	CliMessage() string
	Error() string
}

type CollectedError struct {
	err InternalError
}

func (e CollectedError) Error() string {
	return e.err.Error()
}

func (e CollectedError) CliMessage() string {
	return e.err.CliMessage()
}

type ErrorCollector struct {

}

type CollectedErrors interface {
	GetAll() []CollectedError
}

type CollectedErrorsImpl struct {
	errs []CollectedError
}

func (c CollectedErrorsImpl) GetAll() []CollectedError {
	return c.errs
}

func (c CollectedErrorsImpl) Resolve() CollectedErrors {
	return c
}

type CollectedErrorsPromise interface {
	Resolve() CollectedErrors
}

type CollectedErrorsPromiseChan <-chan CollectedErrors

func (p CollectedErrorsPromiseChan) Resolve() CollectedErrors {
	return <-p
}

func CreateErrorCollector() *ErrorCollector {
	return &ErrorCollector{}
}

func (ec *ErrorCollector) collectSyncParsingError(program parser.LatteParsedProgramPromise, c *context.ParsingContext, out []CollectedError) []CollectedError {
	result := program.Resolve()
	if result.ParsingError() != nil {
		return append(out, CollectedError{
			err: result.ParsingError(),
		})
	} else {
		// Do nothing
	}
	return out
}

func (ec *ErrorCollector) collectSyncParsingErrors(programs []parser.LatteParsedProgramPromise, c *context.ParsingContext, out []CollectedError) []CollectedError {
	for _, program := range programs {
		out = ec.collectSyncParsingError(program, c, out)
	}
	return out
}

func (ec *ErrorCollector) collectSyncTypecheckingError(program type_checker.LatteTypecheckedProgramPromise, c *context.ParsingContext, out []CollectedError) []CollectedError {
	result := program.Resolve()
	if result.TypeCheckingError != nil {
		return append(out, CollectedError{
			err: result.TypeCheckingError,
		})
	} else {
		return ec.collectSyncParsingError(result.Program, c, out)
	}
}

func (ec *ErrorCollector) collectSyncTypecheckingErrors(programs []type_checker.LatteTypecheckedProgramPromise, c *context.ParsingContext, out []CollectedError) []CollectedError {
	for _, program := range programs {
		out = ec.collectSyncTypecheckingError(program, c, out)
	}
	return out
}

func (ec *ErrorCollector) collectSyncCompilationError(program compiler.LatteCompiledProgramPromise, c *context.ParsingContext, out []CollectedError) []CollectedError {
	result := program.Resolve()
	if result.CompilationError != nil {
		return append(out, CollectedError{
			err: result.CompilationError,
		})
	} else {
		return ec.collectSyncTypecheckingError(result.TypecheckedProgram, c, out)
	}
}

func (ec *ErrorCollector) collectSyncCompilationErrors(programs []compiler.LatteCompiledProgramPromise, c *context.ParsingContext, out []CollectedError) []CollectedError {
	for _, program := range programs {
		out = ec.collectSyncCompilationError(program, c, out)
	}
	return out
}

func (ec *ErrorCollector) HandleCompilation(programs []compiler.LatteCompiledProgramPromise, c *context.ParsingContext) CollectedErrorsPromise {
	ret := make(chan CollectedErrors)
	go func() {
		defer close(ret)
		out := []CollectedError{}
		out = ec.collectSyncCompilationErrors(programs, c, out)
		ret <- CollectedErrorsImpl{
			errs: out,
		}
	}()
	return CollectedErrorsPromiseChan(ret)
}

func (ec *ErrorCollector) HandleTypechecking(programs []type_checker.LatteTypecheckedProgramPromise, c *context.ParsingContext) CollectedErrorsPromise {
	ret := make(chan CollectedErrors)
	go func() {
		defer close(ret)
		out := []CollectedError{}
		out = ec.collectSyncTypecheckingErrors(programs, c, out)
		ret <- CollectedErrorsImpl{
			errs: out,
		}
	}()
	return CollectedErrorsPromiseChan(ret)
}

func (ec *ErrorCollector) HandleParsing(programs []parser.LatteParsedProgramPromise, c *context.ParsingContext) CollectedErrorsPromise {
	ret := make(chan CollectedErrors)
	go func() {
		defer close(ret)
		out := []CollectedError{}
		out = ec.collectSyncParsingErrors(programs, c, out)
		ret <- CollectedErrorsImpl{
			errs: out,
		}
	}()
	return CollectedErrorsPromiseChan(ret)
}