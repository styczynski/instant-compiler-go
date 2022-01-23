package x86

import (
	"fmt"

	"github.com/styczynski/latte-compiler/src/ir"
)

type RegistrySpecs struct {
	Size                int
	TopReg              Reg
	Reg                 Reg
	SubRegSize1         Reg
	SubRegSize2         Reg
	FnArgIndex          int
	CanStoreFnArg       bool
	FnReturnIndex       int
	CanReturnFn         bool
	IsPreserved         bool
	ForbidForAllocation bool
}

var ALL_REGS map[Reg]*RegistrySpecs = map[Reg]*RegistrySpecs{
	ESP: {
		Size:                4,
		TopReg:              RSP,
		Reg:                 ESP,
		SubRegSize1:         SP,
		SubRegSize2:         ESP,
		ForbidForAllocation: true,
	},
	EBP: {
		Size:                4,
		TopReg:              RBP,
		Reg:                 EBP,
		SubRegSize1:         BP,
		SubRegSize2:         EBP,
		ForbidForAllocation: true,
	},
	EAX: {
		Size:          4,
		TopReg:        RAX,
		Reg:           EAX,
		SubRegSize1:   AL,
		SubRegSize2:   EAX,
		FnReturnIndex: 0,
		CanReturnFn:   true,
	},
	EBX: {
		Size:        4,
		TopReg:      RBX,
		Reg:         EBX,
		SubRegSize1: BL,
		SubRegSize2: EBX,
		IsPreserved: true,
	},
	ECX: {
		Size:          4,
		TopReg:        RCX,
		Reg:           ECX,
		SubRegSize1:   CL,
		SubRegSize2:   ECX,
		CanStoreFnArg: true,
		FnArgIndex:    3,
	},
	EDX: {
		Size:          4,
		TopReg:        RDX,
		Reg:           EDX,
		SubRegSize1:   DL,
		SubRegSize2:   EDX,
		CanStoreFnArg: true,
		FnArgIndex:    2,
		FnReturnIndex: 1,
		CanReturnFn:   true,
	},
	EDI: {
		Size:          4,
		TopReg:        RDI,
		Reg:           EDI,
		SubRegSize1:   DI,
		SubRegSize2:   EDI,
		CanStoreFnArg: true,
		FnArgIndex:    0,
	},
	ESI: {
		Size:          4,
		TopReg:        RSI,
		Reg:           ESI,
		SubRegSize1:   SI,
		SubRegSize2:   ESI,
		CanStoreFnArg: true,
		FnArgIndex:    1,
	},
	R8L: {
		Size:          4,
		TopReg:        R8,
		Reg:           R8L,
		SubRegSize1:   R8W,
		SubRegSize2:   R8L,
		CanStoreFnArg: true,
		FnArgIndex:    4,
	},
	R9L: {
		Size:          4,
		TopReg:        R9,
		Reg:           R9L,
		SubRegSize1:   R9W,
		SubRegSize2:   R9L,
		CanStoreFnArg: true,
		FnArgIndex:    6,
	},
	R10L: {
		Size:        4,
		TopReg:      R10,
		Reg:         R10L,
		SubRegSize1: R10W,
		SubRegSize2: R10L,
	},
	R11L: {
		Size:        4,
		TopReg:      R11,
		Reg:         R11L,
		SubRegSize1: R11W,
		SubRegSize2: R11L,
	},
	R12L: {
		Size:        4,
		TopReg:      R12,
		Reg:         R12L,
		SubRegSize1: R12W,
		SubRegSize2: R12L,
		IsPreserved: true,
	},
	R13L: {
		Size:        4,
		TopReg:      R13,
		Reg:         R13L,
		SubRegSize1: R13W,
		SubRegSize2: R13L,
		IsPreserved: true,
	},
	R14L: {
		Size:        4,
		TopReg:      R14,
		Reg:         R14L,
		SubRegSize1: R14W,
		SubRegSize2: R14L,
		IsPreserved: true,
	},
	R15L: {
		Size:        4,
		TopReg:      R15,
		Reg:         R15L,
		SubRegSize1: R15W,
		SubRegSize2: R15L,
		IsPreserved: true,
	},
}

func ResizeReg(a Reg, size int) (effectiveSize int, effectiveReg Reg) {
	defer func() {
		fmt.Printf("[?] Resize %v (size %d) into %v (size %d)\n", a, size, effectiveReg, effectiveSize)
	}()
	for reg, regSpecs := range ALL_REGS {
		if AreRegsColliding(&a, &reg) {
			// Matching reg
			if size == 2 && regSpecs.SubRegSize1 != reg && regSpecs.SubRegSize1 != 0 {
				effectiveSize, effectiveReg = 2, regSpecs.SubRegSize1
				return
			} else if size == 4 && regSpecs.SubRegSize2 != reg && regSpecs.SubRegSize2 != 0 {
				effectiveSize, effectiveReg = 4, regSpecs.SubRegSize2
				return
			} else if size == 8 && regSpecs.TopReg != reg && regSpecs.TopReg != 0 {
				effectiveSize, effectiveReg = 8, regSpecs.TopReg
				return
			} else {
				effectiveSize, effectiveReg = regSpecs.Size, reg
				return
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
	return *a == *b || *compReg == usedSpecs.Reg || *compReg == usedSpecs.SubRegSize1 || *compReg == usedSpecs.SubRegSize2 || *compReg == usedSpecs.TopReg
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
