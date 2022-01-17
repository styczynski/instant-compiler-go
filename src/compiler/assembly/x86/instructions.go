package x86

import (
	"fmt"

	"github.com/styczynski/latte-compiler/src/ir"
)

// A Label is a label reference
type RelLabel struct {
	label string
}

func CreateRelLabel(name string) *RelLabel {
	return &RelLabel{
		label: name,
	}
}

func (*RelLabel) isArg() {}

func (l *RelLabel) String() string {
	return l.label
}

func DoRet() []*Instruction {

	instRet := Inst{}
	instRet.Op = RET

	instPreserveStack := Inst{}
	instPreserveStack.Op = POP
	instPreserveStack.MemBytes = 16
	instPreserveStack.Args = Args{
		RBP,
	}

	return []*Instruction{
		{
			Inst: instPreserveStack,
		},
		{
			Inst: instRet,
		},
	}
}

func DoRetMain() []*Instruction {
	// inst := Inst{}
	// inst.Op = RET
	// return &Instruction{
	// 	Inst: inst,
	// }
	instSyscallNo := Inst{}
	instSyscallNo.MemBytes = 4
	instSyscallNo.Op = MOV
	instSyscallNo.Args = Args{
		EBX,
		Imm(1),
	}

	instSwap := Inst{}
	instSwap.Op = XCHG
	instSwap.MemBytes = 4
	instSwap.Args = Args{
		EBX,
		EAX,
	}

	instInterrupt := Inst{}
	instInterrupt.Op = INT
	instInterrupt.Args = Args{
		Imm(128),
	}

	instRet := Inst{}
	instRet.Op = RET

	return []*Instruction{
		{
			Inst: instSyscallNo,
		},
		{
			Inst: instSwap,
		},
		{
			Inst: instInterrupt,
		},
		{
			Inst: instRet,
		},
	}
}

func DoPush(reg Reg, size int) *Instruction {
	inst := Inst{}
	inst.Op = PUSH
	inst.MemBytes = size

	if reg == EAX {
		reg = RAX
	} else if reg == EBX {
		reg = RBX
	} else if reg == ECX {
		reg = RCX
	} else if reg == EDX {
		reg = RDX
	}

	inst.Args = Args{
		reg,
	}
	return &Instruction{
		Inst: inst,
	}
}

func DoCall(label string) *Instruction {
	inst := Inst{}
	inst.Op = CALL
	inst.Args = Args{
		&RelLabel{
			label: label,
		},
	}
	return &Instruction{
		Inst: inst,
	}
}

func DoPop(reg Reg, size int) *Instruction {
	inst := Inst{}
	inst.Op = POP
	inst.MemBytes = size

	if reg == EAX {
		reg = RAX
	} else if reg == EBX {
		reg = RBX
	} else if reg == ECX {
		reg = RCX
	} else if reg == EDX {
		reg = RDX
	}

	inst.Args = Args{
		reg,
	}
	return &Instruction{
		Inst: inst,
	}
}

func DoSwap(val1 Arg, val2 Arg, size int) *Instruction {
	inst := Inst{}
	inst.Op = XCHG
	inst.MemBytes = size
	inst.Args = Args{
		val1,
		val2,
	}
	return &Instruction{
		Inst: inst,
	}
}

func DoUnaryOp(self Arg, val Arg, operation ir.IROperator, size int, argType ir.IRType) []*Instruction {
	inst := Inst{}
	inst.MemBytes = size

	if operation == ir.IR_OP_SELF_ADD {
		inst.Op = ADD

		if argType == ir.IR_STRING {
			// String addition
			instFirst := Inst{}
			instFirst.MemBytes = size
			instFirst.Op = MOV
			instFirst.Args = Args{
				EDI,
				self,
			}

			instSecond := Inst{}
			instSecond.MemBytes = size
			instSecond.Op = MOV
			instSecond.Args = Args{
				ESI,
				val,
			}

			instCall := Inst{}
			instCall.Op = CALL
			instCall.Args = Args{
				&RelLabel{
					label: "AddStrings",
				},
			}

			instResult := Inst{}
			instResult.Op = MOV
			instSecond.MemBytes = size
			instResult.Args = Args{
				EAX,
				self,
			}

			return []*Instruction{
				{
					Inst: instFirst,
				},
				{
					Inst: instSecond,
				},
				{
					Inst: instCall,
				},
				{
					Inst: instResult,
				},
			}
		}

	} else if operation == ir.IR_OP_SELF_DIV {
		inst.Op = IDIV
		inst.Args = Args{}

		instClear := Inst{}
		instClear.Op = XOR
		instClear.Args = Args{
			RDX, RDX,
		}

		return []*Instruction{
			// DoPush(RDX, 8),
			// DoPush(RAX, 8),
			// {
			// 	Inst: instClear,
			// },
			// {
			// 	Inst: inst,
			// },
			// DoPop(RDX, 8),
			// DoPop(RAX, 8),
		}
	} else if operation == ir.IR_OP_SELF_SUB {
		inst.Op = SUB
	} else if operation == ir.IR_OP_SELF_MUL {
		inst.Op = IMUL
	} else {
		panic(fmt.Sprintf("Unsuported operation for DoArithmeticSelfOp: %v", operation))
	}
	inst.Args = Args{
		self,
		val,
	}
	return []*Instruction{
		{
			Inst: inst,
		},
	}
}

func DoSub(target Arg, val Arg, size int) *Instruction {
	inst := Inst{}
	inst.Op = SUB
	inst.MemBytes = size
	inst.Args = Args{
		target,
		val,
	}
	return &Instruction{
		Inst: inst,
	}
}

func DoMov(from Arg, to Arg, size int) *Instruction {
	inst := Inst{}
	inst.Op = MOV
	inst.MemBytes = size
	inst.Args = Args{
		from,
		to,
	}
	return &Instruction{
		Inst: inst,
	}
}

func DoMemoryLoad(index, size int, to Arg) *Instruction {
	inst := Inst{}
	inst.Op = MOV
	inst.MemBytes = size
	inst.Args = Args{
		to,
		GetMemoryVarLocation(index, size),
	}
	return &Instruction{
		Inst: inst,
	}
}

func DoRegSetConditional(reg Reg, subreg Reg, size int, op ir.IROperator) []*Instruction {
	instSet := Inst{}
	instSet.MemBytes = size
	if op == ir.IR_OP_EQ {
		instSet.Op = SETE
	} else if op == ir.IR_OP_LT {
		instSet.Op = SETL
	} else if op == ir.IR_OP_LTEQ {
		instSet.Op = SETLE
	} else if op == ir.IR_OP_GT {
		instSet.Op = SETG
	} else if op == ir.IR_OP_GTEQ {
		instSet.Op = SETGE
	} else if op == ir.IR_OP_NOT_EQ {
		instSet.Op = SETNE
	} else {
		panic(fmt.Sprintf("Invalid operation specified for DoRegSetConditional: %v", op))
	}
	instSet.Args = Args{
		subreg,
	}
	instResize := Inst{}
	instResize.MemBytes = size
	instResize.Op = MOVZX
	instResize.Args = Args{
		reg,
		subreg,
	}
	return []*Instruction{
		{
			Inst: instSet,
		},
		{
			Inst: instResize,
		},
	}
}

func DoMemorySetConditional(index, size int, op ir.IROperator) *Instruction {
	inst := Inst{}
	inst.MemBytes = size
	if op == ir.IR_OP_EQ {
		inst.Op = SETE
	} else {
		panic(fmt.Sprintf("Invalid operation specified for DoMemorySetConditional: %v", op))
	}
	inst.Args = Args{
		GetMemoryVarLocation(index, size),
	}
	return &Instruction{
		Inst: inst,
	}
}

func DoCompare(a Arg, b Arg) *Instruction {
	inst := Inst{}
	inst.Op = CMP
	inst.Args = Args{
		a,
		b,
	}
	return &Instruction{
		Inst: inst,
	}
}

func DoMemoryStore(index, size int, src Reg) *Instruction {
	inst := Inst{}
	inst.Op = MOV
	inst.MemBytes = size
	inst.Args = Args{
		GetMemoryVarLocation(index, size),
		src,
	}
	return &Instruction{
		Inst: inst,
	}
}

func DoRegStoreConst(reg Reg, size int, value int64) *Instruction {
	inst := Inst{}
	inst.Op = MOV
	inst.MemBytes = size
	inst.Args = Args{
		reg,
		Imm(value),
	}
	//inst.DataSize = size * 8
	inst.MemBytes = size
	return &Instruction{
		Inst: inst,
	}
}

func DoMemoryStoreConst(index, size int, value int64) *Instruction {
	inst := Inst{}
	inst.Op = MOV
	inst.MemBytes = size
	inst.Args = Args{
		GetMemoryVarLocation(index, size),
		Imm(value),
	}
	//inst.DataSize = size * 8
	inst.MemBytes = size
	return &Instruction{
		Inst: inst,
	}
}

func GetMemoryVarLocation(index int, size int) Mem {
	return Mem{
		Base: RBP,
		Disp: int64(-index - size),
	}
}

func DoRegToMemoryTransfer(index int, size int, target Reg, memoryToReg bool) *Instruction {
	inst := Inst{}
	inst.Op = MOV
	inst.MemBytes = size
	if memoryToReg {
		inst.Args = Args{
			target,
			GetMemoryVarLocation(index, size),
		}
	} else {
		inst.Args = Args{
			GetMemoryVarLocation(index, size),
			target,
		}
	}
	return &Instruction{
		Inst: inst,
	}
}

func DoMemoryToMemoryTransfer(src, srcSize, target, targetSize int, temp Reg) []*Instruction {
	instGet := Inst{}
	instGet.Op = MOV
	instGet.MemBytes = srcSize
	instGet.Args = Args{
		temp,
		GetMemoryVarLocation(src, srcSize),
	}
	instPut := Inst{}
	instPut.Op = MOV
	instPut.MemBytes = targetSize
	instPut.Args = Args{
		GetMemoryVarLocation(target, targetSize),
		temp,
	}
	return []*Instruction{
		{
			Inst: instGet,
		},
		{
			Inst: instPut,
		},
	}
}

func DoIf(label string, labelElse string, hasElse bool, value Arg) []*Instruction {
	instCmp := Inst{}
	instCmp.Op = CMP
	instCmp.Args = Args{
		value,
		Imm(0),
	}
	instJmp := Inst{}
	instJmp.Op = JNE
	instJmp.Args = Args{
		&RelLabel{
			label: label,
		},
	}
	instJmpElse := Inst{}
	instJmpElse.Op = JMP
	instJmpElse.Args = Args{
		&RelLabel{
			label: labelElse,
		},
	}
	ret := []*Instruction{
		{
			Inst: instCmp,
		},
		{
			Inst: instJmp,
		},
	}
	if hasElse {
		ret = append(ret, &Instruction{
			Inst: instJmpElse,
		})
	}
	return ret
}

func Label(name string) *Instruction {
	return &Instruction{
		Label: name,
	}
}
