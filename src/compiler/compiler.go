package compiler

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/styczynski/latte-compiler/src/errors"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker"
)

type CompilationError struct {
	message     string
	textMessage string
}

var formatErrorBg = color.New(color.BgRed).SprintFunc()
var formatErrorFg = color.New(color.FgHiWhite).SprintFunc()

var formatErrorMessageFg = color.New(color.FgRed).SprintFunc()
var formatErrorMetaInfoFg = color.New(color.FgHiBlue).SprintFunc()

func CreateCompilationError(errorMessage string, details string) *CompilationError {

	textMessage := fmt.Sprintf("%s: %s\n%s\n",
		formatErrorFg(formatErrorBg(fmt.Sprintf(" %s ", "Code Generation Error"))),
		formatErrorMessageFg(errorMessage), formatErrorMessageFg(details))

	return &CompilationError{
		message:     errorMessage,
		textMessage: textMessage,
	}
}

func (e CompilationError) ErrorName() string {
	return "Code Generation Error"
}

func (e CompilationError) Error() string {
	return e.message
}

func (e CompilationError) CliMessage() string {
	return e.textMessage
}

type LatteCompiler struct {
	backend CompilerBackend
}

type CompilerBackend interface {
	Compile(program type_checker.LatteTypecheckedProgram, c *context.ParsingContext) LatteCompiledProgramPromiseChan
	BackendName() string
}

type LatteCompiledProgram struct {
	Program          type_checker.LatteTypecheckedProgram
	CompiledProgram  CompiledProgram
	CompilationError *CompilationError
}

type LatteCompiledProgramPromise interface {
	Resolve() LatteCompiledProgram
}

func (p LatteCompiledProgram) Filename() string {
	return p.Program.Filename()
}

type LatteCompiledProgramPromiseChan <-chan LatteCompiledProgram

func (p LatteCompiledProgramPromiseChan) Resolve() LatteCompiledProgram {
	return <-p
}

func CreateLatteCompiler(Backend CompilerBackend) *LatteCompiler {
	return &LatteCompiler{
		backend: Backend,
	}
}

func (compiler *LatteCompiler) compileAsync(programPromise type_checker.LatteTypecheckedProgramPromise, c *context.ParsingContext) LatteCompiledProgramPromise {
	ret := make(chan LatteCompiledProgram)
	ctx := c.Copy()
	go func() {
		program := programPromise.Resolve()
		defer errors.GeneralRecovery(ctx, "Code generation", program.Filename(), func(message string, textMessage string) {
			ret <- LatteCompiledProgram{
				Program: program,
				CompilationError: &CompilationError{
					message:     message,
					textMessage: textMessage,
				},
			}
		}, func() {
			close(ret)
		})

		backendProcessDescription := fmt.Sprintf("Generate compiled code using backend: %s", compiler.backend.BackendName())

		c.EventsCollectorStream.Start(backendProcessDescription, c, program)
		defer c.EventsCollectorStream.End(backendProcessDescription, c, program)

		compiled := <-compiler.backend.Compile(program, c)
		fmt.Print(compiled.CompiledProgram.ProgramToText())

		ret <- LatteCompiledProgram{
			Program:          compiled.Program,
			CompiledProgram:  compiled.CompiledProgram,
			CompilationError: compiled.CompilationError,
		}
	}()
	return LatteCompiledProgramPromiseChan(ret)
}

func (compiler *LatteCompiler) Compile(programs []type_checker.LatteTypecheckedProgramPromise, c *context.ParsingContext) []LatteCompiledProgramPromise {
	ret := []LatteCompiledProgramPromise{}
	for _, program := range programs {
		ret = append(ret, compiler.compileAsync(program, c))
	}
	return ret
}
