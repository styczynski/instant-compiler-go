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

type CompiledCodeRunContext interface {
	GetOutputFilePathByExtension(extension string) string
	GetCompilerMeta(key string) interface{}
	Call(name string, errorPattern string, args ...interface{}) ([]string, *RunError)
	ReadFileByExt(extension string) ([]byte, error)
}

type RunError struct {
	message     string
	textMessage string
}

func (e RunError) ErrorName() string {
	return "Code Generation Error"
}

func (e RunError) Error() string {
	return e.message
}

func (e RunError) CliMessage() string {
	return e.textMessage
}

func CreateRunError(errorMessage string, details string) *RunError {

	textMessage := fmt.Sprintf("%s: %s\n%s\n",
		formatErrorFg(formatErrorBg(fmt.Sprintf(" %s ", "Run Error"))),
		formatErrorMessageFg(errorMessage), formatErrorMessageFg(details))

	return &RunError{
		message:     errorMessage,
		textMessage: textMessage,
	}
}

type CompilerBackend interface {
	Compile(program type_checker.LatteTypecheckedProgram, c *context.ParsingContext, b *BuildContext) LatteCompiledProgramPromiseChan
	RunCompiledCode(runContext CompiledCodeRunContext, c *context.ParsingContext) ([]string, *RunError)
	BackendName() string
}

type LatteCompiledProgram struct {
	Program          type_checker.LatteTypecheckedProgram
	CompiledProgram  CompiledProgram
	CompilationError *CompilationError
	Backend          CompilerBackend
	OutputFilesByExt map[string]string
	CompilerMeta     map[string]interface{}
}

type LatteCompiledProgramPromise interface {
	Resolve() LatteCompiledProgram
}

func (p LatteCompiledProgram) Filename() string {
	return p.Program.Filename()
}

func (p LatteCompiledProgram) Resolve() LatteCompiledProgram {
	return p
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

		buildContext := CreateBuildContext(program, c)

		defer errors.GeneralRecovery(ctx, "Code generation", program.Filename(), func(message string, textMessage string) {
			ret <- LatteCompiledProgram{
				Program: program,
				Backend: compiler.backend,
				CompilationError: &CompilationError{
					message:     message,
					textMessage: textMessage,
				},
			}
		}, func() {
			close(ret)
		})

		if program.TypeCheckingError != nil {
			ret <- LatteCompiledProgram{
				Program: program,
				Backend: compiler.backend,
			}
			return
		}

		if program.Program.ParsingError() != nil {
			ret <- LatteCompiledProgram{
				Program: program,
				Backend: compiler.backend,
			}
			return
		}

		backendProcessDescription := fmt.Sprintf("Generate compiled code using backend: %s", compiler.backend.BackendName())

		c.EventsCollectorStream.Start(backendProcessDescription, c, program)
		defer c.EventsCollectorStream.End(backendProcessDescription, c, program)

		compiled := <-compiler.backend.Compile(program, c, buildContext)

		c.EventsCollectorStream.EmitOutputFiles(backendProcessDescription, c, buildContext.GetOutputFiles())

		ret <- LatteCompiledProgram{
			Program:          compiled.Program,
			Backend:          compiler.backend,
			CompiledProgram:  compiled.CompiledProgram,
			CompilationError: compiled.CompilationError,
			OutputFilesByExt: buildContext.GetOutputFilesByExt(),
			CompilerMeta:     buildContext.GetCompilerMeta(),
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
