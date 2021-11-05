package llvm_ast

import "fmt"

type LLVMVar struct {
	ID int64
}

func (p *LLVMVar) Type() LLVMInstructionType {
	return Var
}

func (p *LLVMVar) ToText(emitter EmitterConfig) string {
	return emitter.Emit("%d", p.ID)
}

func (p *LLVMVar) GetTarget(withType bool) string {
	return fmt.Sprintf("%%%d", p.ID)
}

func (p *LLVMVar) IsMovable() bool {
	return true
}

func (p *LLVMVar) MoveTarget(newTarget string) {
	panic("Operation not supported")
}
