package llvm_ast

import (
	"fmt"
	"strconv"
)

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

func (p *LLVMVar) GetDeclaredVariables() []string {
	return []string{}
}

func (p *LLVMVar) GetUsedVariables() []string {
	return []string{fmt.Sprintf("%%%d", p.ID)}
}

func (p *LLVMVar) ReplaceVariable(oldName string, newName string) {
	if fmt.Sprintf("%%%d", p.ID) == oldName {
		i, err := strconv.Atoi(newName[:len(newName)-1])
		if err != nil {
			panic(err)
		}
		p.ID = int64(i)
	}
}
