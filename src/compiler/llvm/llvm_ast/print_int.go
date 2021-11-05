package llvm_ast

type LLVMPrintInt struct {
	Target string
}

func (p *LLVMPrintInt) Type() LLVMInstructionType {
	return PrintInt
}

func (p *LLVMPrintInt) ToText(emitter EmitterConfig) string {
	return emitter.Emit("call void @printInt(i32 %s)", p.Target)
}

func (p *LLVMPrintInt) GetDeclaredVariables() []string {
	return []string{}
}

func (p *LLVMPrintInt) GetUsedVariables() []string {
	if p.Target[0] == '%' {
		return []string{p.Target}
	}
	return []string{}
}

func (p *LLVMPrintInt) ReplaceVariable(oldName string, newName string) {
	if p.Target == oldName {
		p.Target = newName
	}
}
