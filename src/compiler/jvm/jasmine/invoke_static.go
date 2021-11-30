package jasmine

import "strings"

type JasmineInvokeStatic struct {
	Target  string
	Special bool
	Virtual bool
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
	if p.Virtual {
		keyword = "invokevirtual"
	}
	return emitter.Emit("%s %s(%s)%s", keyword, p.Target, strings.Join(p.Args, ""), p.Return)
}

func (p *JasmineInvokeStatic) StackSize(previousStackSize int) int {
	if p.Virtual {
		return previousStackSize - len(p.Args) - 1
	}
	return previousStackSize - len(p.Args)
}
