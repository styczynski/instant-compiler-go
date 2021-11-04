package jasmine

import (
	"strings"
)

type JasmineMethod struct {
	Name string
	Body []JasmineInstruction
}

func (p *JasmineMethod) Type() JasmineInstructionType {
	return Method
}

func (p *JasmineMethod) ToText(emitter EmitterConfig) string {
	methodContents := []string{
		emitter.Emit(".method %s", p.Name),
	}
	for _, ins := range p.Body {
		methodContents = append(methodContents, ins.ToText(emitter.ApplyIdent(1)))
	}
	methodContents = append(methodContents, emitter.Emit(".end"))
	return strings.Join(methodContents, "\n")
}

func (p *JasmineMethod) StackSize(previousStackSize int) int {
	s := 0
	for _, ins := range p.Body {
		s = ins.StackSize(s)
	}
	return s
}
