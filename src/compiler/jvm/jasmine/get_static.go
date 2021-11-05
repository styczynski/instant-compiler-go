package jasmine

type JasmineGetStatic struct {
	Source string
	Object string
}

func (p *JasmineGetStatic) Type() JasmineInstructionType {
	return GetStatic
}

func (p *JasmineGetStatic) ToText(emitter EmitterConfig) string {
	return emitter.Emit("getstatic %s %s", p.Source, p.Object)
}

func (p *JasmineGetStatic) StackSize(previousStackSize int) int {
	return previousStackSize + 1
}
