package llvm_ast

import "fmt"

type LLVMVal struct {
	Value     int
	ValueType string
}

func (p *LLVMVal) Type() LLVMInstructionType {
	return Val
}

func (p *LLVMVal) ToText(emitter EmitterConfig) string {
	return emitter.Emit("%s %d", p.ValueType, p.Value)
}

func (p *LLVMVal) GetTarget(withType bool) string {
	if !withType {
		return fmt.Sprintf("%d", p.Value)
	}
	return fmt.Sprintf("%s %d", p.ValueType, p.Value)
}

func (p *LLVMVal) IsMovable() bool {
	return true
}

func (p *LLVMVal) MoveTarget(newTarget string) {
	panic("Operation not supported")
}
