package parser

type PrinterConfiguration struct {
	SkipStatementIdent bool
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