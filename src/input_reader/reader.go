package input_reader

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/styczynski/latte-compiler/src/events_utils"
	"github.com/styczynski/latte-compiler/src/logs"
	"github.com/styczynski/latte-compiler/src/parser/context"
)

type InputReader interface {
	Read(c *context.ParsingContext) ([]LatteInput, error)
	ResolveInclude(c *context.ParsingContext, includePath string) (LatteInput, error)
}

type LatteInputReader struct {
	input    []string
	includes []string
}

type LatteInput interface {
	Read() ([]byte, error)
	Filename() string
}

type LatteInputImpl struct {
	read     func() ([]byte, error)
	filename func() string
}

func (in *LatteInputImpl) Read() ([]byte, error) {
	return in.read()
}

func (in *LatteInputImpl) Filename() string {
	return in.filename()
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func CreateLatteInputReader(input []string, includes []string) *LatteInputReader {
	return &LatteInputReader{
		input:    input,
		includes: includes,
	}
}

func (reader *LatteInputReader) LogContext(c *context.ParsingContext) map[string]interface{} {
	return map[string]interface{}{}
}

func (reader *LatteInputReader) ResolveInclude(c *context.ParsingContext, includePath string) (LatteInput, error) {
	return localFSResolveInclude(c, includePath)
}

func (reader *LatteInputReader) Read(c *context.ParsingContext) ([]LatteInput, error) {
	c.EventsCollectorStream.Start("Read input", c, events_utils.GeneralEventSource{})
	defer c.EventsCollectorStream.End("Read input", c, events_utils.GeneralEventSource{})

	allInputs := []LatteInput{}
	for _, inp := range reader.input {
		input := inp
		if fileExists(inp) {
			f, err := os.Open(inp)
			if err != nil {
				return nil, err
			}
			allInputs = append(allInputs, &LatteInputImpl{
				read: func() ([]byte, error) {
					defer f.Close()
					return ioutil.ReadAll(f)
				},
				filename: func() string { return input },
			})
			logs.Debug(reader, "Read %s", inp)
		} else if inp == "-" {
			allInputs = append(allInputs, &LatteInputImpl{
				read: func() ([]byte, error) {
					return ioutil.ReadAll(bufio.NewReader(os.Stdin))
				},
				filename: func() string { return "<standard input>" },
			})
			logs.Debug(reader, "Read stdin")
		} else {
			// Use glob
			matches, err := filepath.Glob(inp)
			if err != nil {
				return nil, err
			}
			ret := []LatteInput{}
			for _, path := range matches {
				if path == "." || path == ".." {
					continue
				}
				subreader := CreateLatteInputReader([]string{path}, reader.includes)
				subinputs, err := subreader.Read(c)
				if err != nil {
					return nil, err
				}
				ret = append(ret, subinputs...)
			}
			logs.Debug(reader, "Glob \"%s\" was resolved to %d files", inp, len(ret))
			allInputs = append(allInputs, ret...)
		}
	}

	if len(allInputs) == 0 {
		logs.Warning(reader, "No inputs were detected.")
	}

	inputs, err := mergeAllIncludes(c, allInputs, reader.includes, reader.ResolveInclude)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}
