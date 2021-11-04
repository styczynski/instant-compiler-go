package jasmine

type JasmineReferenceLoad struct {
	Index int64
}

func (p *JasmineReferenceLoad) Type() JasmineInstructionType {
	return ReferenceLoad
}

func (p *JasmineReferenceLoad) ToText(emitter EmitterConfig) string {
	return emitter.Emit("aload_%d", p.Index)
}

func (p *JasmineReferenceLoad) StackSize(previousStackSize int) int {
	return previousStackSize
}
