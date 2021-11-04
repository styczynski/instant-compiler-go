package jasmine

type JasminePush struct {
	Value int64
}

func (p *JasminePush) Type() JasmineInstructionType {
	return Push
}

func (p *JasminePush) ToText(emitter EmitterConfig) string {
	if p.Value >= -128 && p.Value <= 127 {
		return emitter.Emit("bipush %d", p.Value)
	}
	if p.Value > -32768 && p.Value <= 32767 {
		return emitter.Emit("sipush %d", p.Value)
	}
	return emitter.Emit("ldc %d", p.Value)
}

func (p *JasminePush) StackSize(previousStackSize int) int {
	if p.Value >= 0 && p.Value <= 255 {
		return 1 + previousStackSize // bipush
	}
	if p.Value > 255 && p.Value <= 65535 {
		return 2 + previousStackSize // sipush
	}
	return 2 + previousStackSize // ldc
}
