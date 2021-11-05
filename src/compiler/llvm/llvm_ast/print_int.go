package llvm_ast

type LLVMPrintInt struct {
	Target string
}

func (p *LLVMPrintInt) Type() LLVMInstructionType {
	return PrintInt
}

func (p *LLVMPrintInt) ToText(emitter EmitterConfig) string {
	return emitter.Emit("call void @printInt(%s)", p.Target)
}
