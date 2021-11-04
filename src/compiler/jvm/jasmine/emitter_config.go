package jasmine

import (
	"fmt"
	"strings"
)

type EmitterConfig struct {
	Ident int
}

func (c EmitterConfig) Emit(format string, a ...interface{}) string {
	return fmt.Sprintf("%s%s", strings.Repeat("   ", c.Ident), fmt.Sprintf(format, a...))
}

func (c EmitterConfig) ApplyIdent(step int) EmitterConfig{
	return EmitterConfig{
		Ident: c.Ident + step,
	}
}