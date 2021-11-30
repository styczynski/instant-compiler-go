package jasmine

type JasmineConstInt struct {
	Index int64
}

func (p *JasmineConstInt) Type() JasmineInstructionType {
	return ConstInt
}

func (p *JasmineConstInt) ToText(emitter EmitterConfig) string {
	if p.Index == -1 {
		return emitter.Emit("iconst_m1")
	}
	if p.Index < -1 || p.Index > 5 {
		// Fallback to push instruction
		return (&JasminePush{
			Value: p.Index,
		}).ToText(emitter)
	}
	return emitter.Emit("iconst_%d", p.Index)
}

func (p *JasmineConstInt) StackSize(previousStackSize int) int {
	return previousStackSize + 1
}
