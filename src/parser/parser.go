package parser

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/alecthomas/participle/v2"

	"github.com/styczynski/latte-compiler/src/input_reader"
	"github.com/styczynski/latte-compiler/src/parser/ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type LatteParser struct {
	parserInstance *participle.Parser
}

type LatteParsedProgram interface {
	AST() *ast.LatteProgram
	Filename() string
	ParsingError() *ParsingError
	Resolve() LatteParsedProgram
}

type LatteParsedProgramImpl struct {
	ast *ast.LatteProgram
	filename string
	error *ParsingError
}

func (prog *LatteParsedProgramImpl) Resolve() LatteParsedProgram {
	return prog
}

func (prog *LatteParsedProgramImpl) AST() *ast.LatteProgram {
	return prog.ast
}

func (prog *LatteParsedProgramImpl) ParsingError() *ParsingError {
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

func examineParsingErrorMessage(message string, recommendedBracket string) string {
	r := regexp.MustCompile(`unexpected token "([^"]*)" (.*)`)
	matches := r.FindStringSubmatch(message)
	if len(matches) > 0 {
		tokenName := matches[1]
		messageRaw := matches[2]
		message := strings.ReplaceAll(messageRaw[1:len(messageRaw)-1], "expected ", "")
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

		defer close(ret)
		var err error = nil
		reader := input.Read()
		output := &ast.LatteProgram{}
		ctx.ParserInput, err = ioutil.ReadAll(reader)
		if err != nil {
			ret <- &LatteParsedProgramImpl{
				ast:      nil,
				filename: input.Filename(),
				error:      &ParsingError{
					message:     err.Error(),
					textMessage: err.Error(),
				},
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
			}
			return
		}

		var parentSetterVisitor hindley_milner.ExpressionMapper
		parentSetterVisitor = func(parent hindley_milner.Expression, e hindley_milner.Expression) hindley_milner.Expression {
			node := e.(ast.TraversableNode)
			if parent == e {
				return e
			}
			if node.Parent() != nil {
				// Prevent infinite loop
				return e
			}
			node.OverrideParent(parent.(interface{}).(ast.TraversableNode))
			e.Visit(parent, parentSetterVisitor)
			return e
		}

		output.Visit(output, parentSetterVisitor)
		ret <- &LatteParsedProgramImpl{
			ast:      output,
			filename: input.Filename(),
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
			},
		}
	}

	programs := []LatteParsedProgramPromise{}
	for _, input := range inputs {
		programs = append(programs, p.parseAsync(c, input))
	}

	return programs
}