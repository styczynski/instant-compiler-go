package printer

import (
	"github.com/styczynski/latte-compiler/src/parser"
)

type LattePrinter struct {}

func CreateLattePrinter() *LattePrinter {
	return &LattePrinter{}
}

func (p *LattePrinter) StructRepr(program *parser.LatteProgram, c *parser.ParsingContext) string {
	return program.Print(c)
}
