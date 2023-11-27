package input_reader

import (
	"github.com/styczynski/latte-compiler/src/events_utils"
	"github.com/styczynski/latte-compiler/src/parser/context"
)

type LatteConstInputReader struct {
	inputs   []LatteInput
	includes []string
}

func CreateLatteConstInputReader(inputs []LatteInput, includes []string) *LatteConstInputReader {
	return &LatteConstInputReader{
		inputs:   inputs,
		includes: includes,
	}
}

func (reader *LatteConstInputReader) ResolveInclude(c *context.ParsingContext, includePath string) (LatteInput, error) {
	return localFSResolveInclude(c, includePath)
}

func (reader *LatteConstInputReader) Read(c *context.ParsingContext) ([]LatteInput, error) {
	c.EventsCollectorStream.Start("Read input", c, events_utils.GeneralEventSource{})
	defer c.EventsCollectorStream.End("Read input", c, events_utils.GeneralEventSource{})

	inputs, err := mergeAllIncludes(c, reader.inputs, reader.includes, reader.ResolveInclude)
	if err != nil {
		return nil, err
	}

	return inputs, nil
}
