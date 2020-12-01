package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/fatih/color"

	"github.com/styczynski/latte-compiler/src/errors"
)

type LatteParser struct {
	parserInstance *participle.Parser
	parserInput []byte
	printer CodeFormatter
}

func CreateLatteParser(codeFormatter CodeFormatter) *LatteParser {
	paserInstance := participle.MustBuild(&LatteProgram{},
		//participle.Lexer(iniLexer),
		participle.UseLookahead(2),
		//participle.Unquote("String"),
	)
	return &LatteParser{
		parserInstance: paserInstance,
		printer: codeFormatter,
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

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func Max(x int, y int ) int {
	if x < y {
		return y
	}
	return x
}

func indentCodeLines(message string, errorLine int, lineStart int) string {
	lines := strings.Split(message, "\n")
	newLines := []string{}
	curLineNo := lineStart
	startTrimMode := true
	for _, line := range lines {
		add := true
		if len(line) > 0 {
			startTrimMode = false
		} else if startTrimMode {
			add = false
		}
		if add {
			lineMarker := ""
			if curLineNo == errorLine {
				lineMarker = "> "
			}
			newLines = append(newLines, fmt.Sprintf("%6s | %s", fmt.Sprintf("%s%d", lineMarker, curLineNo), line))
		}
		curLineNo++
	}
	return strings.Join(newLines, "\n")
}

func (p *LatteParser) getFileContext(c *ParsingContext, program *LatteProgram, line int, column int) (string, int, int) {
	lineOffset := 3

	// There is no program AST context available so we must read raw input
	if program == nil {
		scanner := bufio.NewScanner(bytes.NewReader(p.parserInput))
		curLineNo := 1
		minLineNo := 10000000
		maxLineNo := 0
		contentLines := []string{}
		for scanner.Scan() {
			if curLineNo >= line-lineOffset && curLineNo <= line+lineOffset {
				if curLineNo > maxLineNo {
					maxLineNo = curLineNo
				}
				if curLineNo < minLineNo {
					minLineNo = curLineNo
				}
				contentLines = append(contentLines, scanner.Text())
			}
			curLineNo++
		}
		return strings.Join(contentLines, "\n"), minLineNo, maxLineNo
	} else {
		minDistStart := 10000000
		minDistEnd := 10000000

		var start TraversableNode = program
		var end TraversableNode = program

		TraverseAST(program, func(ast TraversableNode) {
			distStart := Abs(Abs(ast.Begin().Line-line) - lineOffset)
			if distStart < minDistStart {
				minDistStart = distStart
				start = ast
			}
			distEnd := Abs(Abs(ast.End().Line-line) - lineOffset)
			if distEnd < minDistEnd {
				minDistEnd = distEnd
				end = ast
			}
		})

		endPos := end.Begin()
		c.PrinterConfiguration.MaxPrintPosition = &endPos
		content := start.Print(c)
		c.PrinterConfiguration.MaxPrintPosition = nil

		return content, start.Begin().Line, end.Begin().Line
	}
}

type CodeFormatter interface {
	FormatRaw(input string) (string, error)
}

var formatErrorBg = color.New(color.BgRed).SprintFunc()
var formatErrorFg = color.New(color.FgHiWhite).SprintFunc()

var formatErrorMessageFg = color.New(color.FgRed).SprintFunc()

var formatErrorMetaInfoFg = color.New(color.FgHiBlue).SprintFunc()

func (p *LatteParser) formatParsingError(c *ParsingContext, parsingError participle.Error) errors.LatteError {

	locationMessage := fmt.Sprintf("%s %s %s %d %s %d",
		formatErrorMetaInfoFg("Error in file"),
		parsingError.Position().Filename,
		formatErrorMetaInfoFg(" in line "),
		parsingError.Position().Line,
		formatErrorMetaInfoFg(", column"),
		parsingError.Position().Column)

	codeContext, lineStart, _ := p.getFileContext(c, nil, parsingError.Position().Line, parsingError.Position().Column)
	formattedCode, err := p.printer.FormatRaw(codeContext)
	if err != nil {
		// Ignore error
		formattedCode = codeContext
	}

	textMessage := fmt.Sprintf("%s\n%s\n %s: %s\n",
		locationMessage,
		indentCodeLines(formattedCode, parsingError.Position().Line, lineStart),
		formatErrorFg(formatErrorBg(" ParserError ")),
		formatErrorMessageFg(formatParsingErrorMessage(parsingError.Message())))

	message := parsingError.Error()
	return &ParsingError{
		message: message,
		textMessage: textMessage,
	}
}

func formatParsingErrorMessage(message string) string {
	r := regexp.MustCompile(`unexpected token "([^"]*)" (.*)`)
	matches := r.FindStringSubmatch(message)
	if len(matches) > 0 {
		tokenName := matches[1]
		messageRaw := matches[2]
		message := strings.ReplaceAll(messageRaw[1:len(messageRaw)-1], "expected ", "")
		return fmt.Sprintf("The parser encountered unexpected keyword. Expected %s got %s.", message, tokenName)
	}
	return message
}

func (p *LatteParser) ParseInput(input io.Reader, c *ParsingContext) (*LatteProgram, errors.LatteError) {
	output := &LatteProgram{}
	var err error
	p.parserInput, err = ioutil.ReadAll(input)
	if err != nil {
		return nil, errors.NewLatteSimpleError(err)
	}

	err = p.parserInstance.ParseBytes("<input>", p.parserInput, output)
	if err != nil {
		return nil, p.formatParsingError(c, err.(participle.Error))
	}
	return output, nil
}