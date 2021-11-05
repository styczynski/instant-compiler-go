package context

import (
	"bufio"
	"bytes"
	"strings"
	"time"

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
	MaxPrintPosition   *lexer.Position
}

func (p PrinterConfiguration) Copy() PrinterConfiguration {
	var pos *lexer.Position = nil
	if p.MaxPrintPosition != nil {
		pos = &(*p.MaxPrintPosition)
	}
	return PrinterConfiguration{
		SkipStatementIdent: p.SkipStatementIdent,
		MaxPrintPosition:   pos,
	}
}

type ProcessingStage struct {
	Start *time.Time
	End   *time.Time
}

func (stage *ProcessingStage) Copy() *ProcessingStage {
	return &ProcessingStage{
		Start: copyTimePtr(stage.Start),
		End:   copyTimePtr(stage.End),
	}
}

type EventCollectorMessageInput interface {
	Filename() string
}

type EventsCollectorStream interface {
	Start(processName string, c *ParsingContext, input EventCollectorMessageInput)
	End(processName string, c *ParsingContext, input EventCollectorMessageInput)
	EmitOutputFiles(processName string, c *ParsingContext, outputFiles map[string]map[string]string)
}

type ParsingContext struct {
	BlockDepth            int
	PrinterConfiguration  PrinterConfiguration
	ParserInput           []byte
	Printer               CodeFormatter
	Start                 *time.Time
	End                   *time.Time
	EventsCollectorStream EventsCollectorStream
}

func copyTimePtr(val *time.Time) *time.Time {
	if val == nil {
		return nil
	}
	newVal := *val
	return &newVal
}

func (c *ParsingContext) Copy() *ParsingContext {
	input := c.ParserInput
	return &ParsingContext{
		BlockDepth:            c.BlockDepth,
		PrinterConfiguration:  c.PrinterConfiguration.Copy(),
		ParserInput:           input,
		Printer:               c.Printer,
		Start:                 copyTimePtr(c.Start),
		End:                   copyTimePtr(c.End),
		EventsCollectorStream: c.EventsCollectorStream,
	}
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func Max(x int, y int) int {
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

func (c *ParsingContext) Close() {
	end := time.Now()
	c.End = &end
}

func NewParsingContext(printer CodeFormatter, eventsColelctorStream EventsCollectorStream) *ParsingContext {
	start := time.Now()
	return &ParsingContext{
		Printer:    printer,
		BlockDepth: 0,
		PrinterConfiguration: PrinterConfiguration{
			SkipStatementIdent: false,
		},
		Start:                 &start,
		EventsCollectorStream: eventsColelctorStream,
	}
}
