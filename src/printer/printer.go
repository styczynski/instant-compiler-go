package printer

import (
	"bytes"
	"github.com/styczynski/latte-compiler/src/parser"
	"github.com/styczynski/latte-compiler/src/printer/chroma/styles"
	"github.com/styczynski/latte-compiler/src/printer/chroma/formatters"
	"github.com/styczynski/latte-compiler/src/printer/chroma/lexers"
)

type LattePrinter struct {}

func CreateLattePrinter() *LattePrinter {
	return &LattePrinter{}
}

func (p *LattePrinter) Raw(program *parser.LatteProgram, c *parser.ParsingContext) string {
	return program.Print(c)
}

func (p *LattePrinter) FormatRaw(input string) (string, error) {
	lexer := lexers.Get("latte")
	style := styles.Get("colorful")
	formatter := formatters.Get("terminal16m")
	if formatter == nil {
		formatter = formatters.Fallback
	}
	contents := []byte(input)
	iterator, err := lexer.Tokenise(nil, string(contents))
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	err = formatter.Format(buf, style, iterator)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (p *LattePrinter) Format(program *parser.LatteProgram, c *parser.ParsingContext) (string, error) {
	return p.FormatRaw(p.Raw(program, c))
}