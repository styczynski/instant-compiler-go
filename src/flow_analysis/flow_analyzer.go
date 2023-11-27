package flow_analysis

import (
	"github.com/styczynski/latte-compiler/src/errors"
	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/ir"
	"github.com/styczynski/latte-compiler/src/parser/ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker"
)

type LatteFlowAnalyzer struct{}

func CreateLatteFlowAnalyzer() *LatteFlowAnalyzer {
	return &LatteFlowAnalyzer{}
}

type FlowAnalysisError struct {
	message     string
	textMessage string
	errorName   string
	source      generic_ast.NormalNode
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
	IR                *ir.IRProgram
	Program           type_checker.LatteTypecheckedProgram
	FlowAnalysisError *FlowAnalysisError
	filename          string
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
	AfterFlowAnalysis(flow cfg.FlowAnalysis)
}

func wrapFlowAnalysisError(err error, source generic_ast.NormalNode, c *context.ParsingContext) *FlowAnalysisError {
	if foldingErr, ok := err.(cfg.ConstFoldingError); ok {
		src := foldingErr.GetSource()
		errorName := "Constant folding error"
		message, textMessage := c.FormatParsingError(
			errorName,
			foldingErr.Error(),
			src.Begin().Line,
			src.Begin().Column,
			src.Begin().Filename,
			"",
			foldingErr.Error(),
		)
		return &FlowAnalysisError{
			message:     message,
			textMessage: textMessage,
			errorName:   errorName,
		}
	} else if flowErr, ok := err.(*errors.LocalizedError); ok {
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
		program := programPromise.Resolve()
		defer errors.GeneralRecovery(ctx, "Flow analysis", program.Filename(), func(message string, textMessage string) {
			r <- LatteAnalyzedProgram{
				Program: program,
				FlowAnalysisError: &FlowAnalysisError{
					message:     message,
					textMessage: textMessage,
					errorName:   "PANIC (Flow analysis)",
					source:      nil,
				},
				filename: program.Filename(),
			}
		}, func() {
			close(r)
		})
		if program.TypeCheckingError != nil {
			r <- LatteAnalyzedProgram{
				Program:  program,
				filename: program.Filename(),
			}
			return
		}
		if program.Program.ParsingError() != nil {
			r <- LatteAnalyzedProgram{
				Program:  program,
				filename: program.Filename(),
			}
			return
		}
		if program.Program.Context() != nil {
			ctx = program.Program.Context()
		}
		astTree := program.Program.AST()
		visitedNodes := map[generic_ast.Expression]interface{}{}
		var mappingVisitor func(parent generic_ast.Expression, e generic_ast.Expression, context generic_ast.VisitorContext)
		mappingVisitor = func(parent generic_ast.Expression, e generic_ast.Expression, context generic_ast.VisitorContext) {
			if _, ok := visitedNodes[e]; ok {
				return
			}
			visitedNodes[e] = true
			if stmt, ok := e.(*ast.Statement); ok {
				if stmt.IsDeclaration() {
					for _, item := range stmt.Declaration.Items {
						val := int64(0)
						if !item.HasInitializer() {
							item.Initializer = &ast.Expression{
								LogicalOperation: &ast.LogicalOperation{
									Equality: &ast.Equality{
										Comparison: &ast.Comparison{
											Addition: &ast.Addition{
												Multiplication: &ast.Multiplication{
													Unary: &ast.Unary{
														UnaryApplication: &ast.UnaryApplication{
															Index: &ast.Index{
																Primary: &ast.Primary{
																	Int: &val,
																},
															},
														},
													},
												},
											},
										},
									},
								},
							}
						}
					}
				} else if stmt.IsUnaryStatement() {
					unaryStmt := stmt.UnaryStatement
					oldStmt := stmt
					val := int64(1)
					opStr := "+"
					if unaryStmt.Operation == "++" {
						opStr = "+"
					} else if unaryStmt.Operation == "--" {
						opStr = "-"
					} else {
						panic("Unsupported unary statement operation")
					}

					stmt.Assignment = &ast.Assignment{
						BaseASTNode: oldStmt.BaseASTNode,
						TargetName:  *unaryStmt.TargetName,
						Value: &ast.Expression{
							LogicalOperation: &ast.LogicalOperation{
								Equality: &ast.Equality{
									Comparison: &ast.Comparison{
										Addition: &ast.Addition{
											Op: opStr,
											Next: &ast.Addition{
												Multiplication: &ast.Multiplication{
													Unary: &ast.Unary{
														UnaryApplication: &ast.UnaryApplication{
															Index: &ast.Index{
																Primary: &ast.Primary{
																	BaseASTNode: unaryStmt.BaseASTNode,
																	Int:         &val,
																},
															},
														},
													},
												},
											},
											Multiplication: &ast.Multiplication{
												Unary: &ast.Unary{
													UnaryApplication: &ast.UnaryApplication{
														Index: &ast.Index{
															Primary: &ast.Primary{
																BaseASTNode: unaryStmt.BaseASTNode,
																Variable:    unaryStmt.TargetName,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						ParentNode: oldStmt.ParentNode,
					}
					stmt.UnaryStatement = nil
				}
			}
			e.Visit(parent, mappingVisitor, context)
		}
		astTree.Visit(astTree, mappingVisitor, generic_ast.NewEmptyVisitorContext())

		var flowVisitor generic_ast.ExpressionVisitor
		visitedNodes = map[generic_ast.Expression]interface{}{}

		var flowErrGlobal *FlowAnalysisError = nil

		irProgram := &ir.IRProgram{
			Statements: []*ir.IRFunction{},
		}
		flowVisitor = func(parent generic_ast.Expression, e generic_ast.Expression, context generic_ast.VisitorContext) {
			if _, ok := visitedNodes[e]; ok {
				return
			}
			visitedNodes[e] = true
			if nodeForAnalysis, ok := (e.(FlowAnalyzableNode)); ok {
				ast := nodeForAnalysis.(generic_ast.NormalNode)
				flow := cfg.CreateFlowAnalysis(ast)

				

				
				
				
				err := flow.ConstFold(c)
				if err != nil {
					if flowErrGlobal == nil {
						flowErrGlobal = wrapFlowAnalysisError(err, nodeForAnalysis.(generic_ast.NormalNode), ctx)
					}
					return
				}
				flow.Rebuild()
				//ast = flow.Output()

				//flow.Optimize(c)
				//nodeForAnalysis.AfterFlowAnalysis(flow)

				//flow.Rebuild()
				//flow.Optimize(c)
				//nodeForAnalysis.AfterFlowAnalysis(flow)

				

				irCode := ir.CreateIR(e, flow, c)
				irProgram.Statements = append(irProgram.Statements, irCode)

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
		astTree.Visit(astTree, flowVisitor, generic_ast.NewEmptyVisitorContext())

		//os.Exit(42)
		if flowErrGlobal != nil {
			r <- LatteAnalyzedProgram{
				Program:           program,
				FlowAnalysisError: flowErrGlobal,
				filename:          program.Filename(),
			}
			return
		}

		r <- LatteAnalyzedProgram{
			Program:  program,
			IR:       irProgram,
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
