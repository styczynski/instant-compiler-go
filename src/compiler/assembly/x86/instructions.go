package x86

import (
	"fmt"

	"github.com/styczynski/latte-compiler/src/ir"
)

type RegistrySpecs struct {
	Normalized          Reg
	Reg8B               Reg
	Reg1B               Reg
	Reg2B               Reg
	Reg4B               Reg
	FnArgIndex          int
	CanStoreFnArg       bool
	FnReturnIndex       int
	CanReturnFn         bool
	IsPreserved         bool
	ForbidForAllocation bool
	DefaultSize         int
}

var ALL_REGS map[Reg]*RegistrySpecs = map[Reg]*RegistrySpecs{
	EAX: {
		Normalized:    EAX,
		Reg8B:         RAX,
		Reg4B:         EAX,
		Reg2B:         AX,
		Reg1B:         AL,
		FnReturnIndex: 0,
		CanReturnFn:   true,
		DefaultSize:   4,
	},
	EBX: {
		Normalized:  EBX,
		Reg8B:       RBX,
		Reg4B:       EBX,
		Reg2B:       BX,
		Reg1B:       BL,
		IsPreserved: true,
		DefaultSize: 4,
	},
	ECX: {
		Normalized:    ECX,
		Reg8B:         RCX,
		Reg4B:         ECX,
		Reg2B:         CX,
		Reg1B:         CL,
		CanStoreFnArg: true,
		FnArgIndex:    3,
		DefaultSize:   4,
	},
	EDX: {
		Normalized:    EDX,
		Reg8B:         RDX,
		Reg4B:         EDX,
		Reg2B:         DX,
		Reg1B:         DL,
		CanStoreFnArg: true,
		FnArgIndex:    2,
		FnReturnIndex: 1,
		CanReturnFn:   true,
		DefaultSize:   4,
	},
	ESI: {
		Normalized:    ESI,
		Reg8B:         RSI,
		Reg4B:         ESI,
		Reg2B:         SI,
		Reg1B:         SIB,
		CanStoreFnArg: true,
		FnArgIndex:    1,
		DefaultSize:   4,
	},
	EDI: {
		Normalized:    EDI,
		Reg8B:         RDI,
		Reg4B:         EDI,
		Reg2B:         DI,
		Reg1B:         DIB,
		CanStoreFnArg: true,
		FnArgIndex:    0,
		DefaultSize:   4,
	},
	EBP: {
		Normalized:          EBP,
		Reg8B:               RBP,
		Reg4B:               EBP,
		Reg2B:               BP,
		Reg1B:               BPB,
		ForbidForAllocation: true,
		DefaultSize:         4,
	},
	ESP: {
		Normalized:          ESP,
		Reg8B:               RSP,
		Reg4B:               ESP,
		Reg2B:               SP,
		Reg1B:               SPB,
		ForbidForAllocation: true,
		DefaultSize:         4,
	},
	R8L: {
		Normalized:    R8L,
		Reg8B:         R8,
		Reg4B:         R8L,
		Reg2B:         R8W,
		Reg1B:         R8B,
		CanStoreFnArg: true,
		FnArgIndex:    4,
		DefaultSize:   4,
	},
	R9L: {
		Normalized:    R9L,
		Reg8B:         R9,
		Reg4B:         R9L,
		Reg2B:         R9W,
		Reg1B:         R9B,
		CanStoreFnArg: true,
		FnArgIndex:    6,
		DefaultSize:   4,
	},
	R10L: {
		Normalized:  R10L,
		Reg8B:       R10,
		Reg4B:       R10L,
		Reg2B:       R10W,
		Reg1B:       R10B,
		DefaultSize: 4,
	},
	R11L: {
		Normalized:  R11L,
		Reg8B:       R11,
		Reg4B:       R11L,
		Reg2B:       R11W,
		Reg1B:       R11B,
		DefaultSize: 4,
	},
	R12L: {
		Normalized:  R12L,
		Reg8B:       R12,
		Reg4B:       R12L,
		Reg2B:       R12W,
		Reg1B:       R12B,
		IsPreserved: true,
		DefaultSize: 4,
	},
	R13L: {
		Normalized:  R13L,
		Reg8B:       R13,
		Reg4B:       R13L,
		Reg2B:       R13W,
		Reg1B:       R13B,
		IsPreserved: true,
		DefaultSize: 4,
	},
	R14L: {
		Normalized:  R14L,
		Reg8B:       R14,
		Reg4B:       R14L,
		Reg2B:       R14W,
		Reg1B:       R14B,
		IsPreserved: true,
		DefaultSize: 4,
	},
	R15L: {
		Normalized:  R15L,
		Reg8B:       R15,
		Reg4B:       R15L,
		Reg2B:       R15W,
		Reg1B:       R15B,
		IsPreserved: true,
		DefaultSize: 4,
	},
}

func ResizeReg(a Reg, size int) (effectiveSize int, effectiveReg Reg) {
	defer func() {
		fmt.Printf("[?] Resize %v (size %d) into %v (size %d)\n", a, size, effectiveReg, effectiveSize)
	}()
	for reg, regSpecs := range ALL_REGS {
		if AreRegsColliding(&a, &reg) {
			// Matching reg
			if size == 1 && regSpecs.Reg1B != 0 {
				effectiveSize, effectiveReg = 1, regSpecs.Reg1B
				return
			} else if size == 2 && regSpecs.Reg2B != 0 {
				effectiveSize, effectiveReg = 2, regSpecs.Reg2B
				return
			} else if size == 4 && regSpecs.Reg4B != 0 {
				effectiveSize, effectiveReg = 4, regSpecs.Reg4B
				return
			} else if size == 8 && regSpecs.Reg8B != 0 {
				effectiveSize, effectiveReg = 8, regSpecs.Reg8B
				return
			} else {
				if regSpecs.Reg4B != 0 {
					effectiveSize, effectiveReg = 4, regSpecs.Reg4B
					return
				} else if regSpecs.Reg8B != 0 {
					effectiveSize, effectiveReg = 8, regSpecs.Reg8B
					return
				} else if regSpecs.Reg2B != 0 {
					effectiveSize, effectiveReg = 2, regSpecs.Reg2B
					return
				} else if regSpecs.Reg1B != 0 {
					effectiveSize, effectiveReg = 1, regSpecs.Reg1B
					return
				} else {
					panic("Invalid registry")
				}
			}
		}
	}
	panic("Invalid registry")
}

func AreRegsCollidingConst(a *Reg, b Reg) bool {
	bp := b
	return AreRegsColliding(a, &bp)
}

func AreRegsColliding(a *Reg, b *Reg) bool {
	if a == nil || b == nil {
		return false
	}
	var usedSpecs *RegistrySpecs = nil
	var compReg *Reg = nil
	if specs, ok := ALL_REGS[*a]; ok && usedSpecs == nil {
		usedSpecs = specs
		compReg = b
	} else if specs, ok := ALL_REGS[*b]; ok && usedSpecs == nil {
		usedSpecs = specs
		compReg = a
	}
	if usedSpecs == nil {
		return false
	}
	return *a == *b || *compReg == usedSpecs.Reg1B || *compReg == usedSpecs.Reg2B || *compReg == usedSpecs.Reg4B || *compReg == usedSpecs.Reg8B
}

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
	instPreserveStack.Op = LEAVE
	//instPreserveStack.MemBytes = 16
	//instPreserveStack.Args = Args{
	//	RBP,
	//}

	return []*Instruction{
		{
			Inst: instPreserveStack,
		},
		{
			Inst: instRet,
		},
	}
}

func doRawMov(to Arg, from Arg, size int) *Instruction {
	if to == from {
		return DoNop()
	}
	instMove := Inst{}
	instMove.Op = MOV
	instMove.MemBytes = size
	instMove.Args = Args{
		to,
		from,
	}
	return &Instruction{
		Inst: instMove,
	}
}

func doRawSwap(a Arg, b Arg, size int) *Instruction {
	if a == b {
		return DoNop()
	}
	instSwap := Inst{}
	instSwap.Op = XCHG
	instSwap.MemBytes = size
	instSwap.Args = Args{
		a,
		b,
	}
	return &Instruction{
		Inst: instSwap,
	}
}

func DoRetMain() []*Instruction {
	// inst := Inst{}
	// inst.Op = RET
	// return &Instruction{
	// 	Inst: inst,
	// }
	instSyscallNo := doRawMov(EBX, Imm(1), 4)

	instSwap := doRawSwap(EBX, EAX, 4)

	instInterrupt := Inst{}
	instInterrupt.Op = INT
	instInterrupt.Args = Args{
		Imm(128),
	}

	instRet := Inst{}
	instRet.Op = RET

	return []*Instruction{
		instSyscallNo,
		instSwap,
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

func DoEmptyPop() *Instruction {
	inst := Inst{}
	inst.Op = ADD
	inst.MemBytes = 32
	inst.Args = Args{
		RSP,
		Imm(8),
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

func DoSwapRegistryWithMemory(index, size int, from Reg) *Instruction {
	fromSize, fromResized := ResizeReg(from, size)
	return doRawSwap(fromResized, GetMemoryVarLocation(index, size), fromSize)
}

func DoSwap(to Reg, from Reg, size int) *Instruction {
	fromSize, fromResized := ResizeReg(from, size)
	toSize, toResized := ResizeReg(to, size)
	// if fromSize != toSize {
	// 	fromSize, fromResized = ResizeReg(to, toSize)
	// }
	if fromSize != toSize {
		panic(fmt.Sprintf("Couldn't match sized of two registries: %v and %v", to, from))
	}

	return doRawSwap(toResized, fromResized, fromSize)
}

func DoUnaryOp(self Arg, val Arg, operation ir.IROperator, size int, argType ir.IRType) []*Instruction {
	inst := Inst{}
	inst.MemBytes = size

	if operation == ir.IR_OP_SELF_ADD {
		inst.Op = ADD

		if argType == ir.IR_STRING {
			// String addition
			instFirst := doRawMov(EDI, self, size)
			instSecond := doRawMov(ESI, val, size)

			instCall := Inst{}
			instCall.Op = CALL
			instCall.Args = Args{
				&RelLabel{
					label: "AddStrings",
				},
			}

			instResult := doRawMov(self, EAX, size)

			return []*Instruction{
				instFirst,
				instSecond,
				{
					Inst: instCall,
				},
				instResult,
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

func DoNop() *Instruction {
	inst := Inst{}
	inst.Op = NOP
	return &Instruction{
		Inst: inst,
	}
}

func DoRegistryCopy(to Reg, from Reg, size int) *Instruction {
	fromSize, fromResized := ResizeReg(from, size)
	toSize, toResized := ResizeReg(to, size)
	// if fromSize != toSize {
	// 	fromSize, fromResized = ResizeReg(from, toSize)
	// }
	if fromSize != toSize {
		panic(fmt.Sprintf("Couldn't match sized of two registries: %v and %v", to, from))
	}
	return doRawMov(toResized, fromResized, fromSize)
}

func DoZeroRegistry(to Reg, size int) *Instruction {
	toSize, toResized := ResizeReg(to, size)
	return doRawMov(toResized, Imm(0), toSize)
}

func DoMemoryLoad(index, size int, to Reg) *Instruction {
	toSize, toResized := ResizeReg(to, size)
	return doRawMov(toResized, GetMemoryVarLocation(index, size), toSize)
}

func DoRegSetConditional(reg Reg, size int, op ir.IROperator) []*Instruction {
	instSet := Inst{}
	instSet.MemBytes = size

	_, resizedReg := ResizeReg(reg, 1)

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
		resizedReg,
	}
	instResize := Inst{}
	instResize.MemBytes = size
	instResize.Op = MOVZX
	instResize.Args = Args{
		reg,
		resizedReg,
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

func DoLoadDataIntoRegister(to Reg, size int, label string) *Instruction {
	toSize, toResized := ResizeReg(to, size)
	return doRawMov(toResized, CreateRelLabel(label), toSize)
}

func DoLoadDataIntoMemory(index, size int, label string) *Instruction {
	return doRawMov(GetMemoryVarLocation(index, size), CreateRelLabel(label), size)
}

func DoMemoryStore(index, size int, from Reg) *Instruction {
	fromSize, fromResized := ResizeReg(from, size)
	return doRawMov(GetMemoryVarLocation(index, size), fromResized, fromSize)
}

func DoRegStoreConst(reg Reg, size int, value int64) *Instruction {
	return doRawMov(reg, Imm(value), size)
}

func DoMemoryStoreConst(index, size int, value int64) *Instruction {
	return doRawMov(GetMemoryVarLocation(index, size), Imm(value), size)
}

func GetMemoryVarLocation(index int, size int) Mem {
	return Mem{
		Base: RBP,
		Disp: int64(-index - size),
	}
}

func DoRegToMemoryTransfer(index int, size int, target Reg, memoryToReg bool) *Instruction {
	if memoryToReg {
		return doRawMov(
			target,
			GetMemoryVarLocation(index, size),
			size,
		)
	} else {
		return doRawMov(
			GetMemoryVarLocation(index, size),
			target,
			size,
		)
	}
}

func DoMemoryToMemoryTransfer(src, srcSize, target, targetSize int, temp Reg) []*Instruction {
	instGet := doRawMov(temp, GetMemoryVarLocation(src, srcSize), srcSize)
	instPut := doRawMov(GetMemoryVarLocation(target, targetSize), temp, targetSize)
	return []*Instruction{
		instGet,
		instPut,
	}
}

func DoJump(label string) *Instruction {
	instJmp := Inst{}
	instJmp.Op = JMP
	instJmp.Args = Args{
		CreateRelLabel(label),
	}
	return &Instruction{
		Inst: instJmp,
	}
}

func DoIf(label string, labelElse string, hasElse bool, value Arg, negated bool) []*Instruction {
	instCmp := Inst{}
	instCmp.Op = CMP
	instCmp.Args = Args{
		value,
		Imm(0),
	}
	instJmp := Inst{}
	instJmp.Op = JNE
	if negated {
		instJmp.Op = JE
	}
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

func GetRegisterForFunctionArg(index int) *RegistrySpecs {
	for _, reg := range ALL_REGS {
		if reg.CanStoreFnArg && reg.FnArgIndex == index {
			return reg
		}
	}
	return nil
}

func GetRegisterForFunctionReturn(index int) *RegistrySpecs {
	for _, reg := range ALL_REGS {
		if reg.CanReturnFn && reg.FnReturnIndex == index {
			return reg
		}
	}
	return nil
}
