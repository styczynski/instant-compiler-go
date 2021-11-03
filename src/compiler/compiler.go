package compiler

import (
	"github.com/styczynski/latte-compiler/src/errors"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker"
)

type CompilationError struct {
	message string
	textMessage string
}

func  (e CompilationError) ErrorName() string {
	return "Code Generation Error"
}

func (e CompilationError) Error() string {
	return e.message
}

func (e CompilationError) CliMessage() string {
	return e.textMessage
}

type LatteCompiler struct {}

type LatteCompiledProgram struct {
	Program          type_checker.LatteTypecheckedProgram
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

func CreateLatteCompiler() *LatteCompiler {
	return &LatteCompiler{}
}

func (compiler *LatteCompiler) compileAsync(programPromise type_checker.LatteTypecheckedProgramPromise, c *context.ParsingContext) LatteCompiledProgramPromise {
	ret := make(chan LatteCompiledProgram)
	ctx := c.Copy()
	go func() {
		program := programPromise.Resolve()
		defer errors.GeneralRecovery(ctx, "Code generation", program.Filename(), func(message string, textMessage string) {
			ret <- LatteCompiledProgram{
				Program:          program,
				CompilationError: &CompilationError{
					message:     message,
					textMessage: textMessage,
				},
			}
		}, func() {
			close(ret)
		})

		c.EventsCollectorStream.Start("Generate compiled code", c, program)
		defer c.EventsCollectorStream.End("Generate compiled code", c, program)

		ret <- LatteCompiledProgram{
			Program:          program,
			CompilationError: nil,
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
