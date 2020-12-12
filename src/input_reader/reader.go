package input_reader

import (
	"io"
	"strings"

	"github.com/styczynski/latte-compiler/src/parser/context"
)

type LatteInputReader struct {}

func CreateLatteInputReader() *LatteInputReader {
	return &LatteInputReader{}
}

func (compiler *LatteInputReader) Read(c *context.ParsingContext) (io.Reader, error) {
	c.ProcessingStageStart("Read input")
	defer c.ProcessingStageEnd("Read input")

	// TODO: Implement
	return strings.NewReader(""), nil
}
