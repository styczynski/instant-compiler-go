package jasmine

import "strings"

type JasmineInvokeStatic struct {
	Target  string
	Special bool
	Args    []string
	Return  string
}

func (p *JasmineInvokeStatic) Type() JasmineInstructionType {
	return InvokeStatic
}

func (p *JasmineInvokeStatic) ToText(emitter EmitterConfig) string {
	keyword := "invokestatic"
	if p.Special {
		keyword = "invokespecial"
	}
	return emitter.Emit("%s %s(%s)%s", keyword, p.Target, strings.Join(p.Args, ";"), p.Return)
}

func (p *JasmineInvokeStatic) StackSize(previousStackSize int) int {
	return previousStackSize - len(p.Args)
}
