package jasmine

type JasmineSwap struct {
}

func (p *JasmineSwap) Type() JasmineInstructionType {
	return Pop
}

func (p *JasmineSwap) ToText(emitter EmitterConfig) string {
	return emitter.Emit("swap")
}

func (p *JasmineSwap) StackSize(previousStackSize int) int {
	return previousStackSize
}
