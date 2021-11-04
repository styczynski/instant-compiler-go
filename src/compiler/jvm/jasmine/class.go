package jasmine

import (
	"strings"

	"github.com/styczynski/latte-compiler/src/compiler"
)

type JasmineClass struct {
	Super   string
	Methods []JasmineMethod
}

func (p *JasmineClass) Type() JasmineInstructionType {
	return Class
}

func (p *JasmineClass) ToText(emitter EmitterConfig) string {
	methodContents := []string{
		emitter.Emit(".super %s", p.Super),
	}
	for _, ins := range p.Methods {
		methodContents = append(methodContents, ins.ToText(emitter.ApplyIdent(1)))
	}
	methodContents = append(methodContents, emitter.Emit(".end"))
	return strings.Join(methodContents, "\n")
}

func (p *JasmineClass) StackSize(previousStackSize int) int {
	s := 0
	for _, ins := range p.Methods {
		s = ins.StackSize(s)
	}
	return s
}

func (p *JasmineClass) Validate() *compiler.CompilationError {
	for _, method := range p.Methods {
		err := method.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}
