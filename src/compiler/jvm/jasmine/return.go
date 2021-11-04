package jasmine

type JasmineReturn struct {
}

func (p *JasmineReturn) Type() JasmineInstructionType {
	return Return
}

func (p *JasmineReturn) ToText(emitter EmitterConfig) string {
	return emitter.Emit("return")
}

func (p *JasmineReturn) StackSize(previousStackSize int) int {
	return previousStackSize
}
