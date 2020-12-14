package flow_analysis

import (
	"fmt"
	"os"

	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser"
	"github.com/styczynski/latte-compiler/src/parser/context"
)

type LatteFlowAnalyzer struct {}

func CreateLatteFlowAnalyzer() *LatteFlowAnalyzer {
	return &LatteFlowAnalyzer{}
}

type LatteAnalyzedProgram struct {
	Program parser.LatteParsedProgram
	filename string
}

func (p LatteAnalyzedProgram) Filename() string {
	return p.filename
}

func (p LatteAnalyzedProgram) Resolve() LatteAnalyzedProgram {
	return p
}

type LatteAnalyzedProgramPromise interface {
	Resolve() LatteAnalyzedProgram
}

type LatteAnalyzedProgramPromiseChan <-chan LatteAnalyzedProgram

func (p LatteAnalyzedProgramPromiseChan) Resolve() LatteAnalyzedProgram {
	return <-p
}

func (fa *LatteFlowAnalyzer) analyzerAsync(programPromise parser.LatteParsedProgramPromise, c *context.ParsingContext) LatteAnalyzedProgramPromise {
	r := make(chan LatteAnalyzedProgram)
	ctx := c.Copy()
	go func() {
		defer close(r)
		program := programPromise.Resolve()
		if program.Context() != nil {
			ctx = program.Context()
		}
		if program.ParsingError() != nil {
			fmt.Print(program.ParsingError().CliMessage())
			os.Exit(1)
		}
		ast := program.AST()

		cfgGraph := cfg.FromStmts([]generic_ast.NormalNode{ ast })
		fmt.Printf("\n\nENTIRE GRAPH:\n\n")
		fmt.Print(cfgGraph.PrintCFG(ctx))
		os.Exit(0)

		r <- LatteAnalyzedProgram{
			Program:  program,
			filename: program.Filename(),
		}
	}()
	return LatteAnalyzedProgramPromiseChan(r)
}

func (fa *LatteFlowAnalyzer) Analyze(programs []parser.LatteParsedProgramPromise, c *context.ParsingContext) []LatteAnalyzedProgramPromise {
	ret := []LatteAnalyzedProgramPromise{}
	for _, programPromise := range programs {
		ret = append(ret, fa.analyzerAsync(programPromise, c))
		// TODO: Remove
		ret[len(ret)-1].Resolve()
	}
	return ret
}