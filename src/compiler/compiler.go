package compiler

import (
	"github.com/styczynski/latte-compiler/src/parser"
	"github.com/styczynski/latte-compiler/src/parser/context"
)

type LatteCompiler struct {}

func CreateLatteCompiler() *LatteCompiler {
	return &LatteCompiler{}
}

func (compiler *LatteCompiler) Compile(programs parser.LatteParsedProgramCollection, c *context.ParsingContext) (CompiledProgram, error) {
	c.ProcessingStageStart("Generate compiled code")
	defer c.ProcessingStageEnd("Generate compiled code")

	// TODO: Implement
	return CompiledProgram{}, nil
}
