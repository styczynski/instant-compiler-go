package compiler

import (
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker"
)

type CompilationError struct {
	message string
	textMessage string
}

func (e CompilationError) Error() string {
	return e.message
}

func (e CompilationError) CliMessage() string {
	return e.textMessage
}

type LatteCompiler struct {}

type LatteCompiledProgram struct {
	TypecheckedProgram type_checker.LatteTypecheckedProgram
	CompilationError *CompilationError
}

type LatteCompiledProgramPromise interface {
	Resolve() LatteCompiledProgram
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
	//ctx := c.Copy()
	go func() {
		defer close(ret)
		program := programPromise.Resolve()
		ret <- LatteCompiledProgram{
			TypecheckedProgram: program,
			CompilationError:   nil,
		}
	}()
	return LatteCompiledProgramPromiseChan(ret)
}

func (compiler *LatteCompiler) Compile(programs []type_checker.LatteTypecheckedProgramPromise, c *context.ParsingContext) []LatteCompiledProgramPromise {
	c.ProcessingStageStart("Generate compiled code")
	defer c.ProcessingStageEnd("Generate compiled code")

	ret := []LatteCompiledProgramPromise{}
	for _, program := range programs {
		ret = append(ret, compiler.compileAsync(program, c))
	}
	return ret
}
