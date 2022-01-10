package x86

import "golang.org/x/arch/x86/x86asm"

func DoRet() *Instruction {
	inst := x86asm.Inst{}
	inst.Op = x86asm.RET
	return &Instruction{
		Inst: inst,
	}
}

func DoPush(reg x86asm.Reg) *Instruction {
	inst := x86asm.Inst{}
	inst.Op = x86asm.PUSH
	inst.Args = x86asm.Args{
		reg,
	}
	return &Instruction{
		Inst: inst,
	}
}

func DoPop(reg x86asm.Reg) *Instruction {
	inst := x86asm.Inst{}
	inst.Op = x86asm.POP
	inst.Args = x86asm.Args{
		reg,
	}
	return &Instruction{
		Inst: inst,
	}
}

func DoMov(from x86asm.Arg, to x86asm.Arg) *Instruction {
	inst := x86asm.Inst{}
	inst.Op = x86asm.MOV
	inst.Args = x86asm.Args{
		from,
		to,
	}
	return &Instruction{
		Inst: inst,
	}
}
