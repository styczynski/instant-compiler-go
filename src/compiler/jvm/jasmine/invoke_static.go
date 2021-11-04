package jasmine

type JasmineInvokeStatic struct {
	Target  string
	Special bool
}

func (p *JasmineInvokeStatic) Type() JasmineInstructionType {
	return InvokeStatic
}

func (p *JasmineInvokeStatic) ToText(emitter EmitterConfig) string {
	if p.Special {
		return emitter.Emit("invokespecial %s", p.Target)
	}
	return emitter.Emit("invokestatic %s", p.Target)
}

func (p *JasmineInvokeStatic) StackSize(previousStackSize int) int {
	return previousStackSize
}
