package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/alecthomas/participle/v2"

	"github.com/styczynski/latte-compiler/src/errors"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/input_reader"
	"github.com/styczynski/latte-compiler/src/parser/ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
)

type LatteParser struct {
	parserInstance *participle.Parser
}

type AbstractParsingError interface {
	ErrorName() string
	Error() string
	CliMessage() string
}

type LatteParsedProgram interface {
	AST() *ast.LatteProgram
	Filename() string
	ParsingError() AbstractParsingError
	Resolve() LatteParsedProgram
	Context() *context.ParsingContext
}

type LatteParsedProgramImpl struct {
	ast *ast.LatteProgram
	filename string
	error AbstractParsingError
	context *context.ParsingContext
}

func (prog *LatteParsedProgramImpl) Context() *context.ParsingContext {
	return prog.context
}

func (prog *LatteParsedProgramImpl) Resolve() LatteParsedProgram {
	return prog
}

func (prog *LatteParsedProgramImpl) AST() *ast.LatteProgram {
	return prog.ast
}

func (prog *LatteParsedProgramImpl) ParsingError() AbstractParsingError {
	return prog.error
}

func (prog *LatteParsedProgramImpl) Filename() string {
	return prog.filename
}

type LatteParsedProgramCollection interface {
	GetAll() []LatteParsedProgram
}

type LatteParsedPrograms []LatteParsedProgram

func (p LatteParsedPrograms) GetAll() []LatteParsedProgram {
	return p
}

type LatteParsedProgramPromise interface {
	Resolve() LatteParsedProgram
}

type LatteParsedProgramPromiseChanel <-chan LatteParsedProgram

func (p LatteParsedProgramPromiseChanel) Resolve() LatteParsedProgram {
	return <- p
}

func CreateLatteParser() *LatteParser {
	paserInstance := participle.MustBuild(&ast.LatteProgram{},
		//participle.Lexer(iniLexer),
		participle.UseLookahead(2),
		//participle.Unquote("String"),
	)
	return &LatteParser{
		parserInstance: paserInstance,
	}
}

type ParsingError struct {
	message string
	textMessage string
}

func (e *ParsingError) ErrorName() string {
	return "Parsing error"
}

func (e *ParsingError) Error() string {
	return e.message
}

func (e *ParsingError) CliMessage() string {
	return e.textMessage
}

type ParsingPanicError struct {
	message string
	textMessage string
}

func (e *ParsingPanicError) ErrorName() string {
	return "PANIC (Parsing)"
}

func (e *ParsingPanicError) Error() string {
	return fmt.Sprintf("%s", e.message)
}

func (e *ParsingPanicError) CliMessage() string {
	return fmt.Sprintf("%s", e.textMessage)
}

func examineParsingErrorMessage(message string, recommendedBracket string) string {
	r := regexp.MustCompile(`unexpected token "([^"]*)" (.*)`)
	matches := r.FindStringSubmatch(message)
	if len(matches) > 0 {
		tokenName := matches[1]
		messageRaw := matches[2]
		message := strings.Replace(messageRaw[1:len(messageRaw)-1], "expected ", "")
		additionalInfo := ""
		if len(recommendedBracket) > 0 {
			if recommendedBracket == "/" {
				additionalInfo = fmt.Sprintf(" You probably forgot a leading \"/\" on the start of the comment?")
			} else {
				additionalInfo = fmt.Sprintf(" You probably forgot a \"%s\" bracket.", recommendedBracket)
			}
		}
		if message == "\")\" Statement" && tokenName == "{" {
			return fmt.Sprintf("Parser encountered invalid bracketing. You misplaced a curly bracket { after/in the same place where ) is.")
		} else if message == "\")\" Statement" {
			return fmt.Sprintf("Parser encountered invalid syntax. The closing bracket \")\" and a statement was expected in place of \"%s\"", tokenName)
		} else if message == "\";\"" {
			return fmt.Sprintf("The parser encountered unexpected keyword. The semicolon was expected in place of \"%s\". Please make sure you have semicolons in right places.", tokenName)
		}

		if message == "\")\" Block" {
			message = "bracket \")\" with instructions block"
		}

		suggestion := searchKeywords(tokenName, ast.SUGGESTED_KEYWORDS)
		suggestionInfo := ""
		if len(suggestion) > 0 {
			suggestionInfo = fmt.Sprintf("\n                Did you mean \"%s\"?", suggestion)
		}
		return fmt.Sprintf("The parser encountered unexpected keyword. Expected %s got \"%s\".%s%s", message, tokenName, additionalInfo, suggestionInfo)
	}

	r = regexp.MustCompile(`unexpected token "([^"]*)"`)
	matches = r.FindStringSubmatch(message)
	if len(matches) > 0 {
		return fmt.Sprintf("The parser encountered unexpected code fragment. \"%s\" was unexpected here.", matches[1])
	}

	return message
}

func (p *LatteParser) parseAsync(c *context.ParsingContext, input input_reader.LatteInput) LatteParsedProgramPromise {

	ret := make(chan LatteParsedProgram)
	ctx := c.Copy()
	go func() {
		c.EventsCollectorStream.Start("Parse input", c, input)
		defer c.EventsCollectorStream.End("Parse input", c, input)

		defer errors.GeneralRecovery(ctx, "Parsing", input.Filename(), func(message string, textMessage string) {
			ret <- &LatteParsedProgramImpl{
				ast:      nil,
				filename: input.Filename(),
				error:      &ParsingPanicError{
					message:     message,
					textMessage: textMessage,
				},
				context: ctx,
			}
		}, func() {
			close(ret)
		})

		var err error = nil
		output := &ast.LatteProgram{}
		q, err := input.Read()
		ctx.ParserInput = q
		if err != nil {
			ret <- &LatteParsedProgramImpl{
				ast:      nil,
				filename: input.Filename(),
				error:      &ParsingError{
					message:     err.Error(),
					textMessage: err.Error(),
				},
				context: ctx,
			}
			return
		}

		err = p.parserInstance.ParseBytes(input.Filename(), ctx.ParserInput, output)
		if err != nil {
			parserError := err.(participle.Error)
			bracket := tryInsertingBrackets(p.parserInstance, ctx.ParserInput, parserError.Position().Line, parserError.Position().Column)
			message, textMessage := ctx.FormatParsingError(
				(&ParsingError{}).ErrorName(),
				parserError.Error(),
				parserError.Position().Line,
				parserError.Position().Column,
				parserError.Position().Filename,
				bracket,
				examineParsingErrorMessage(parserError.Message(), bracket))
			ret <- &LatteParsedProgramImpl{
				ast:      nil,
				filename: input.Filename(),
				error:      &ParsingError{
					message:     message,
					textMessage: textMessage,
				},
				context: ctx,
			}
			return
		}

		var parentSetterVisitor generic_ast.ExpressionVisitor
		visitedNodes := map[generic_ast.Expression]interface{}{}
		parentSetterVisitor = func(parent generic_ast.Expression, e generic_ast.Expression, context generic_ast.VisitorContext) {
			if _, ok := visitedNodes[e]; ok {
				return
			}
			visitedNodes[e] = true
			node, ok := e.(generic_ast.TraversableNode)
			if !ok {
				return
			}
			if parent == e {
				return
			}
			if node.Parent() != nil {
				// Prevent infinite loop
				return
			}
			node.OverrideParent(parent.(interface{}).(generic_ast.TraversableNode))

			e.Visit(parent, parentSetterVisitor, generic_ast.NewEmptyVisitorContext())
			return
		}

		var nodeSyntaxValidationVisitor generic_ast.ExpressionVisitor
		var syntaxValidationError generic_ast.NodeError = nil
		nodeSyntaxValidationVisitor = func(parent generic_ast.Expression, e generic_ast.Expression, context generic_ast.VisitorContext) {
			if _, ok := visitedNodes[e]; ok {
				return
			}
			visitedNodes[e] = true
			node, ok := e.(generic_ast.TraversableNode)
			if !ok {
				return
			}
			if parent == e {
				return
			}
			if nodeWithValidation, ok := node.(generic_ast.NodeWithSyntaxValidation); ok {
				validationError := nodeWithValidation.Validate(c)
				if validationError != nil {
					syntaxValidationError = validationError
					return
				}
			}
			e.Visit(parent, nodeSyntaxValidationVisitor, generic_ast.NewEmptyVisitorContext())
			return
		}

		output.Visit(output, parentSetterVisitor, generic_ast.NewEmptyVisitorContext())
		visitedNodes = map[generic_ast.Expression]interface{}{}
		output.Visit(output, nodeSyntaxValidationVisitor, generic_ast.NewEmptyVisitorContext())

		var out generic_ast.Expression = output
		if nodeWithValidation, ok := out.(generic_ast.NodeWithSyntaxValidation); ok {
			validationError := nodeWithValidation.Validate(c)
			if validationError != nil {
				syntaxValidationError = validationError
			}
		}

		if syntaxValidationError != nil {
			pos := syntaxValidationError.GetSource().(generic_ast.NodeWithPosition).Begin()
			message, textMessage := ctx.FormatParsingError(
				syntaxValidationError.ErrorName(),
				syntaxValidationError.GetMessage(),
				pos.Line,
				pos.Column,
				pos.Filename,
				"",
				examineParsingErrorMessage(syntaxValidationError.GetMessage(), ""))
			ret <- &LatteParsedProgramImpl{
				ast:      nil,
				filename: input.Filename(),
				error:      &ParsingError{
					message:     message,
					textMessage: textMessage,
				},
				context: ctx,
			}
			return
		}

		ret <- &LatteParsedProgramImpl{
			ast:      output,
			filename: input.Filename(),
			context: ctx,
		}
	}()

	return LatteParsedProgramPromiseChanel(ret)
}

func (p *LatteParser) ParseInput(reader *input_reader.LatteInputReader, c *context.ParsingContext) []LatteParsedProgramPromise {
	var err error
	inputs, err := reader.Read(c)
	if err != nil {
		return []LatteParsedProgramPromise{
			&LatteParsedProgramImpl{
				ast:      nil,
				filename: "<unknown input>",
				error:    &ParsingError{
					message:     err.Error(),
					textMessage: err.Error(),
				},
				context: c,
			},
		}
	}

	programs := []LatteParsedProgramPromise{}
	for _, input := range inputs {
		programs = append(programs, p.parseAsync(c, input))
	}

	return programs
}