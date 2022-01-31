package input_reader

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/styczynski/latte-compiler/src/events_utils"
	"github.com/styczynski/latte-compiler/src/parser/context"
)

var DEFAULT_RUNTIME_INCLUDES []string = []string{
	"io.lat",
	"string.lat",
	"error.lat",
	"assert.lat",
}

func localFSResolveInclude(c *context.ParsingContext, includePath string) (LatteInput, error) {
	c.EventsCollectorStream.Start("Resolve include", c, events_utils.GeneralEventSource{})
	defer c.EventsCollectorStream.End("Resolve include", c, events_utils.GeneralEventSource{})

	basePath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	libPath, err := filepath.Abs(filepath.Join(basePath, "lib"))
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(libPath); os.IsNotExist(err) {
		return nil, err
	}

	includeFilePath := filepath.Join(libPath, includePath)
	if _, err := os.Stat(includeFilePath); os.IsNotExist(err) {
		return nil, err
	}

	f, err := os.Open(includeFilePath)
	if err != nil {
		return nil, err
	}
	ret := &LatteInputImpl{
		read: func() ([]byte, error) {
			defer f.Close()
			return ioutil.ReadAll(f)
		},
		filename: func() string { return includeFilePath },
	}

	return ret, nil
}

func mergeAllIncludes(c *context.ParsingContext, inputs []LatteInput, includes []string, resolver func(c *context.ParsingContext, includePath string) (LatteInput, error)) ([]LatteInput, error) {
	includesContents := [][]byte{}

	for _, inc := range includes {
		resolvedInclude, err := resolver(c, inc)
		if err != nil {
			return nil, err
		}
		includeInput, err := resolvedInclude.Read()
		if err != nil {
			return nil, err
		}
		includesContents = append(includesContents, includeInput)
	}

	newInputs := []LatteInput{}
	for _, input := range inputs {
		oldInput := input
		newInputs = append(newInputs, &LatteInputImpl{
			read: func() ([]byte, error) {
				inputBytes, err := oldInput.Read()
				if err != nil {
					return nil, err
				}
				newBytes := []byte{}
				for _, con := range includesContents {
					newBytes = append(newBytes, con...)
					newBytes = append(newBytes, []byte("\n\n\n")...)
				}
				newBytes = append(newBytes, inputBytes...)
				return newBytes, nil
			},
			filename: func() string {
				return oldInput.Filename()
			},
		})
	}
	return newInputs, nil
}
