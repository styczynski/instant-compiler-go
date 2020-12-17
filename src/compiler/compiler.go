package compiler

import (
	"github.com/styczynski/latte-compiler/src/flow_analysis"
	"github.com/styczynski/latte-compiler/src/parser/context"
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
	Program          flow_analysis.LatteAnalyzedProgram
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

func (compiler *LatteCompiler) compileAsync(programPromise flow_analysis.LatteAnalyzedProgramPromise, c *context.ParsingContext) LatteCompiledProgramPromise {
	ret := make(chan LatteCompiledProgram)
	//ctx := c.Copy()
	go func() {
		defer close(ret)
		program := programPromise.Resolve()
		c.EventsCollectorStream.Start("Generate compiled code", c, program)
		defer c.EventsCollectorStream.End("Generate compiled code", c, program)

		ret <- LatteCompiledProgram{
			Program:          program,
			CompilationError: nil,
		}
	}()
	return LatteCompiledProgramPromiseChan(ret)
}

func (compiler *LatteCompiler) Compile(programs []flow_analysis.LatteAnalyzedProgramPromise, c *context.ParsingContext) []LatteCompiledProgramPromise {
	ret := []LatteCompiledProgramPromise{}
	for _, program := range programs {
		ret = append(ret, compiler.compileAsync(program, c))
	}
	return ret
}
