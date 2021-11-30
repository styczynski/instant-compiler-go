package input_reader

import (
	"github.com/styczynski/latte-compiler/src/events_utils"
	"github.com/styczynski/latte-compiler/src/parser/context"
)

type LatteConstInputReader struct {
	inputs []LatteInput
}

func CreateLatteConstInputReader(inputs []LatteInput) *LatteConstInputReader {
	return &LatteConstInputReader{
		inputs: inputs,
	}
}

func (reader *LatteConstInputReader) Read(c *context.ParsingContext) ([]LatteInput, error) {
	c.EventsCollectorStream.Start("Read input", c, events_utils.GeneralEventSource{})
	defer c.EventsCollectorStream.End("Read input", c, events_utils.GeneralEventSource{})

	return reader.inputs, nil
}
