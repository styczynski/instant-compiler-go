package llvm_ast

import "fmt"

type IntOpType int64

const (
	Add IntOpType = iota
	Mul
	Sub
	Div
)

type LLVMIntOp struct {
	Operation IntOpType
	Target    string
	Arg1      string
	Arg2      string
	ArgType   string
}

func CreateLLVMIntOp(Target string, ArgType string, Arg1 string, Operation string, Arg2 string) *LLVMIntOp {
	if Operation == "+" {
		return &LLVMIntOp{Operation: Add, Target: Target, ArgType: ArgType, Arg1: Arg1, Arg2: Arg2}
	}
	if Operation == "-" {
		return &LLVMIntOp{Operation: Sub, Target: Target, ArgType: ArgType, Arg1: Arg1, Arg2: Arg2}
	}
	if Operation == "*" {
		return &LLVMIntOp{Operation: Mul, Target: Target, ArgType: ArgType, Arg1: Arg1, Arg2: Arg2}
	}
	if Operation == "/" {
		return &LLVMIntOp{Operation: Div, Target: Target, ArgType: ArgType, Arg1: Arg1, Arg2: Arg2}
	}
	panic(fmt.Sprintf("Unknown operation symbol in CreateLLVMIntOp(): %s", Operation))
}

func (p *LLVMIntOp) Type() LLVMInstructionType {
	return IntOp
}

func (p *LLVMIntOp) ToText(emitter EmitterConfig) string {
	opText := ""
	if p.Operation == Add {
		opText = "add"
	}
	if p.Operation == Sub {
		opText = "sub"
	}
	if p.Operation == Mul {
		opText = "mul"
	}
	if p.Operation == Div {
		opText = "sdiv"
	}
	if opText == "" {
		panic(fmt.Sprintf("Invalid operation was used in context of LLVMIntOp: %d", p.Operation))
	}
	return emitter.Emit("%%%s = %s %s %s, %s", p.Target, opText, p.ArgType, p.Arg1, p.Arg2)
}

func (p *LLVMIntOp) GetTarget(withType bool) string {
	return fmt.Sprintf("%%%s", p.Target)
}

func (p *LLVMIntOp) IsMovable() bool {
	return false
}

func (p *LLVMIntOp) MoveTarget(newTarget string) {
	p.Target = newTarget
}

func (p *LLVMIntOp) GetDeclaredVariables() []string {
	return []string{fmt.Sprintf("%%%s", p.Target)}
}

func (p *LLVMIntOp) GetUsedVariables() []string {
	used := []string{}
	if p.Arg1[0] == '%' {
		used = append(used, p.Arg1)
	}
	if p.Arg2[0] == '%' {
		used = append(used, p.Arg2)
	}
	return used
}

func (p *LLVMIntOp) ReplaceVariable(oldName string, newName string) {
	if p.Arg1 == oldName {
		p.Arg1 = newName
	}
	if p.Arg2 == oldName {
		p.Arg2 = newName
	}
	if fmt.Sprintf("%%%s", p.Target) == oldName {
		p.Target = newName[1:]
	}
}
