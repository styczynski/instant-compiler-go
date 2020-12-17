package flow_analysis

import (
	"github.com/styczynski/latte-compiler/src/errors"
	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker"
)

type LatteFlowAnalyzer struct {}

func CreateLatteFlowAnalyzer() *LatteFlowAnalyzer {
	return &LatteFlowAnalyzer{}
}

type FlowAnalysisError struct {
	message string
	textMessage string
	errorName string
	source generic_ast.NormalNode
}

func (e *FlowAnalysisError) ErrorName() string {
	return e.errorName
}

func (e *FlowAnalysisError) Error() string {
	return e.message
}

func (e *FlowAnalysisError) CliMessage() string {
	return e.textMessage
}

type LatteAnalyzedProgram struct {
	Program type_checker.LatteTypecheckedProgram
	FlowAnalysisError *FlowAnalysisError
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

type FlowAnalyzableNode interface {
	OnFlowAnalysis(flow cfg.FlowAnalysis) error
}

func wrapFlowAnalysisError(err error, source generic_ast.NormalNode, c *context.ParsingContext) *FlowAnalysisError {
	if flowErr, ok := err.(*errors.LocalizedError); ok {
		src := flowErr.Source()
		errorName := flowErr.ErrorName()
		message, textMessage := c.FormatParsingError(
			errorName,
			flowErr.Error(),
			src.Begin().Line,
			src.Begin().Column,
			src.Begin().Filename,
			"",
			flowErr.Error(),
		)
		return &FlowAnalysisError{
			message:     message,
			textMessage: textMessage,
			errorName:   errorName,
		}
	} else {
		src := source.(interface{}).(generic_ast.NodeWithPosition)
		errorName := "Flow error"
		message, textMessage := c.FormatParsingError(
			errorName,
			err.Error(),
			src.Begin().Line,
			src.Begin().Column,
			src.Begin().Filename,
			"",
			err.Error(),
		)
		return &FlowAnalysisError{
			message:     message,
			textMessage: textMessage,
			errorName:   errorName,
		}
	}
}

func (fa *LatteFlowAnalyzer) analyzerAsync(programPromise type_checker.LatteTypecheckedProgramPromise, c *context.ParsingContext) LatteAnalyzedProgramPromise {
	r := make(chan LatteAnalyzedProgram)
	ctx := c.Copy()
	go func() {
		defer close(r)
		program := programPromise.Resolve()
		if program.TypeCheckingError != nil {
			r <- LatteAnalyzedProgram{
				Program:  program,
				filename: program.Filename(),
			}
			return
		}
		if program.Program.Context() != nil {
			ctx = program.Program.Context()
		}
		ast := program.Program.AST()

		var flowVisitor generic_ast.ExpressionVisitor
		visitedNodes := map[generic_ast.Expression]interface{}{}

		var flowErrGlobal *FlowAnalysisError = nil

		flowVisitor = func (parent generic_ast.Expression, e generic_ast.Expression, context generic_ast.VisitorContext) {
			if _, ok := visitedNodes[e]; ok {
				return
			}
			visitedNodes[e] = true
			if nodeForAnalysis, ok := (e.(FlowAnalyzableNode)); ok {
				ast := nodeForAnalysis.(generic_ast.NormalNode)
				flow := cfg.CreateFlowAnalysis(ast)


				//fmt.Printf("\n\nENTIRE GRAPH:\n\n")
				//fmt.Print(flow.Print(ctx))
				//fmt.Printf("\nPerform fold()\n")
				flow.ConstFold(c)
				flow.Rebuild()
				ast = flow.Output()

				//fmt.Printf("Fold done.\n")
				//fmt.Printf("\n\nENTIRE CODE:\n\n%s", ast.Print(c))
				//fmt.Printf("\n\nENTIRE GRAPH:\n\n")
				//fmt.Print(flow.Print(ctx))
				//fmt.Printf("Yeah.\n")

				customErr := nodeForAnalysis.OnFlowAnalysis(flow)
				if customErr != nil {
					if flowErrGlobal == nil {
						flowErrGlobal = wrapFlowAnalysisError(customErr, nodeForAnalysis.(generic_ast.NormalNode), ctx)
					}
					return
				}

				return
			}
			e.Visit(parent, flowVisitor, context)
			return
		}
		ast.Visit(ast, flowVisitor, generic_ast.NewEmptyVisitorContext())

		if flowErrGlobal != nil {
			r <- LatteAnalyzedProgram{
				Program:  program,
				FlowAnalysisError: flowErrGlobal,
				filename: program.Filename(),
			}
			return
		}

		r <- LatteAnalyzedProgram{
			Program:  program,
			filename: program.Filename(),
		}
	}()
	return LatteAnalyzedProgramPromiseChan(r)
}

func (fa *LatteFlowAnalyzer) Analyze(programs []type_checker.LatteTypecheckedProgramPromise, c *context.ParsingContext) []LatteAnalyzedProgramPromise {
	ret := []LatteAnalyzedProgramPromise{}
	for _, programPromise := range programs {
		ret = append(ret, fa.analyzerAsync(programPromise, c))
	}
	return ret
}