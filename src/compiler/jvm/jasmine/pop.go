package jasmine

type JasminePop struct {
}

func (p *JasminePop) Type() JasmineInstructionType {
	return Pop
}

func (p *JasminePop) ToText(emitter EmitterConfig) string {
	return emitter.Emit("pop")
}

func (p *JasminePop) StackSize(previousStackSize int) int {
	return previousStackSize - 1
}
