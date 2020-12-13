package flow_analysis

import (
	"github.com/alecthomas/repr"

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
	go func() {
		defer close(r)
		program := programPromise.Resolve()
		ast := program.AST()

		cfgGraph := cfg.FromStmts([]generic_ast.NormalNode{ ast })
		repr.Print(cfgGraph)

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
	}
	return ret
}