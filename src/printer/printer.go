package printer

import (
	"github.com/alecthomas/repr"

	"github.com/styczynski/latte-compiler/src/parser"
)

type LattePrinter struct {}

func CreateLattePrinter() *LattePrinter {
	return &LattePrinter{}
}

func (p *LattePrinter) StructRepr(program *parser.LatteProgram) string {
	return repr.String(program)
}
