package context

import (
	"bufio"
	"bytes"
	"strings"
	"github.com/alecthomas/participle/v2/lexer"
)

type TraversableNode interface {
	GetChildren() []TraversableNode
	GetNode() interface{}
	Begin() lexer.Position
	End() lexer.Position
	Print(c *ParsingContext) string
}

type CodeFormatter interface {
	FormatRaw(input string) (string, error)
}

type PrinterConfiguration struct {
	SkipStatementIdent bool
	MaxPrintPosition *lexer.Position
}

type ParsingContext struct {
	BlockDepth int
	PrinterConfiguration PrinterConfiguration
	ParserInput []byte
	Printer CodeFormatter
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

func TraverseAST(node TraversableNode, visitor func(ast TraversableNode)) {
	children := node.GetChildren()
	for _, child := range children {
		visitor(child)
		TraverseAST(child, visitor)
	}
}

func (c *ParsingContext) GetFileContext(program TraversableNode, line int, column int) (string, int, int) {
	lineOffset := 3

	// There is no program AST context available so we must read raw input
	if program == nil {
		scanner := bufio.NewScanner(bytes.NewReader(c.ParserInput))
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

func (c *ParsingContext) BlockPush() {
	c.BlockDepth += 1
}

func (c *ParsingContext) BlockPop() {
	c.BlockDepth -= 1
}

func NewParsingContext(printer CodeFormatter) *ParsingContext {
	return &ParsingContext{
		Printer: printer,
		BlockDepth: 0,
		PrinterConfiguration: PrinterConfiguration{
			SkipStatementIdent: false,
		},
	}
}