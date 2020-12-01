package parser

import "github.com/alecthomas/participle/v2/lexer"

type PrinterConfiguration struct {
	SkipStatementIdent bool
	MaxPrintPosition *lexer.Position
}

type ParsingContext struct {
	BlockDepth int
	PrinterConfiguration PrinterConfiguration
}

func (c *ParsingContext) BlockPush() {
	c.BlockDepth += 1
}

func (c *ParsingContext) BlockPop() {
	c.BlockDepth -= 1
}

func NewParsingContext() *ParsingContext {
	return &ParsingContext{
		BlockDepth: 0,
		PrinterConfiguration: PrinterConfiguration{
			SkipStatementIdent: false,
		},
	}
}