package jasmine

type JasmineLoadInt struct {
	Index int64
}

func (p *JasmineLoadInt) Type() JasmineInstructionType {
	return LoadInt
}

func (p *JasmineLoadInt) ToText(emitter EmitterConfig) string {
	if p.Index >= 0 && p.Index <= 3 {
		return emitter.Emit("iload_%d", p.Index)
	}
	return emitter.Emit("iload %d", p.Index)
}

func (p *JasmineLoadInt) StackSize(previousStackSize int) int {
	return previousStackSize + 1
}
