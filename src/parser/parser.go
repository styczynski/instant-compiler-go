package parser

import (
	"io"
	"github.com/alecthomas/participle/v2"
)

type LatteParser struct {
	parserInstance *participle.Parser
}

func CreateLatteParser() *LatteParser {
	paserInstance := participle.MustBuild(&LatteProgram{},
		//participle.Lexer(iniLexer),
		participle.UseLookahead(2),
		//participle.Unquote("String"),
	)
	return &LatteParser{
		parserInstance: paserInstance,
	}
}

func (p *LatteParser) ParseInput(input io.Reader, c *ParsingContext) (*LatteProgram, error) {
	output := &LatteProgram{}
	err := p.parserInstance.Parse("", input, output)
	if err != nil {
		return nil, err
	}
	return output, nil
}