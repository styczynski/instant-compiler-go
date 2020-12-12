package parser

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/alecthomas/participle/v2"

	"github.com/styczynski/latte-compiler/src/errors"
	"github.com/styczynski/latte-compiler/src/parser/ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type LatteParser struct {
	parserInstance *participle.Parser
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

func (p *LatteParser) ParseInput(input io.Reader, c *context.ParsingContext) (*ast.LatteProgram, errors.LatteError) {
	output := &ast.LatteProgram{}
	var err error
	c.ParserInput, err = ioutil.ReadAll(input)
	if err != nil {
		return nil, errors.NewLatteSimpleError(err)
	}

	err = p.parserInstance.ParseBytes("<input>", c.ParserInput, output)
	if err != nil {
		parserError := err.(participle.Error)
		bracket := tryInsertingBrackets(p.parserInstance, c.ParserInput, parserError.Position().Line, parserError.Position().Column)
		message, textMessage := c.FormatParsingError(
			"Parsing Error",
			parserError.Error(),
			parserError.Position().Line,
			parserError.Position().Column,
			parserError.Position().Filename,
			bracket,
			examineParsingErrorMessage(parserError.Message(), bracket))
		return nil, &ParsingError{
			message:     message,
			textMessage: textMessage,
		}
	}

	var parentSetterVisitor hindley_milner.ExpressionMapper
	parentSetterVisitor = func(parent hindley_milner.Expression, e hindley_milner.Expression) hindley_milner.Expression {
		node := e.(ast.TraversableNode)
		//fmt.Printf("CUR PAR %v\n", node.Parent())
		if parent == e {
			return e
		}
		if node.Parent() != nil {
			// Prevent infinite loop
			return e
		}
		//fmt.Printf(" VISIT => %v FROM %v\n", c, parent)
		//fmt.Printf(" VISIT => %s FROM %s\n", node.(ast.PrintableNode).Print(c), parent.(ast.PrintableNode).Print(c))
		node.OverrideParent(parent.(interface{}).(ast.TraversableNode))
		e.Visit(parent, parentSetterVisitor)
		return e
	}

	output.Visit(output, parentSetterVisitor)

	return output, nil
}