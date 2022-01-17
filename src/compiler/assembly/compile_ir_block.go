package assembly

import (
	"fmt"

	"github.com/styczynski/latte-compiler/src/compiler/assembly/allocation"
	"github.com/styczynski/latte-compiler/src/compiler/assembly/x86"
	"github.com/styczynski/latte-compiler/src/ir"
	"github.com/styczynski/latte-compiler/src/parser/context"
)

func (backend CompilerX86Backend) compileIRConst(ret []*x86.Instruction, instr *ir.IRConst, name string, alloc ir.IRAllocation) (error, []*x86.Instruction) {
	if mem, ok := allocation.IsAllocMem(alloc); ok {
		ret = append(ret, x86.DoMemoryStoreConst(
			mem.Index,
			mem.Size,
			instr.Value,
		))
		return nil, ret
	} else if reg, ok := allocation.IsAllocReg(alloc); ok {
		ret = append(ret, x86.DoRegStoreConst(
			reg.Reg,
			reg.Size,
			instr.Value,
		))
		return nil, ret
	} else {
		return fmt.Errorf("Unsupported transfer to memory from %s", alloc.String()), nil
	}
}

func (backend CompilerX86Backend) compileIRCopy(ret []*x86.Instruction, instr *ir.IRCopy, name string, alloc ir.IRAllocation, srcAlloc ir.IRAllocation) (error, []*x86.Instruction) {
	if mem, ok := allocation.IsAllocMem(alloc); ok {
		if srcMem, ok := allocation.IsAllocMem(srcAlloc); ok {
			// Mem -> Mem transfer
			ret = append(ret, x86.DoMemoryToMemoryTransfer(
				srcMem.Index,
				srcMem.Size,
				mem.Index,
				mem.Size,
				x86.EAX,
			)...)
			return nil, ret
		} else if srcReg, ok := allocation.IsAllocReg(srcAlloc); ok {
			// Mem -> Reg
			ret = append(ret, x86.DoRegToMemoryTransfer(
				mem.Index,
				mem.Size,
				srcReg.Reg,
				false,
			))
			return nil, ret
		} else {
			return fmt.Errorf("Unsupported transfer to memory from %s", srcAlloc.String()), nil
		}
	} else if reg, ok := allocation.IsAllocReg(alloc); ok {
		if srcMem, ok := allocation.IsAllocMem(srcAlloc); ok {
			// Reg -> Mem
			ret = append(ret, x86.DoRegToMemoryTransfer(
				srcMem.Index,
				srcMem.Size,
				reg.Reg,
				true,
			))
			return nil, ret
		} else if srcReg, ok := allocation.IsAllocReg(srcAlloc); ok {
			// Reg -> Reg
			ret = append(ret, x86.DoMov(
				reg.Reg,
				srcReg.Reg,
				srcReg.Size,
			))
			return nil, ret
		} else {
			return fmt.Errorf("Unsupported transfer to memory from %s", srcAlloc.String()), nil
		}
	} else {
		return fmt.Errorf("Unsupported transfer to %s", alloc.String()), nil
	}
}

func (backend CompilerX86Backend) compileIROpUnary(ret []*x86.Instruction, instr *ir.IRExpression, op ir.IROperator, name string, alloc ir.IRAllocation, srcAlloc ir.IRAllocation) (error, []*x86.Instruction) {
	if mem, ok := allocation.IsAllocMem(alloc); ok {
		if srcReg, ok := allocation.IsAllocReg(srcAlloc); ok {
			// MEM = Op(REG)
			ret = append(ret, x86.DoSwap(
				srcReg.Reg,
				x86.GetMemoryVarLocation(mem.Index, mem.Size),
				srcReg.Size,
			))
			ret = append(ret, x86.DoUnaryOp(
				srcReg.Reg,
				x86.GetMemoryVarLocation(mem.Index, mem.Size),
				op,
				srcReg.Size,
			))
			ret = append(ret, x86.DoSwap(
				srcReg.Reg,
				x86.GetMemoryVarLocation(mem.Index, mem.Size),
				srcReg.Size,
			))
			return nil, ret
		} else if _, ok := allocation.IsAllocMem(srcAlloc); ok {
			// MEM = Op(MEM)
			// Impossible
			return nil, ret
		} else {
			return fmt.Errorf("Unsupported operation %v transfer to %s", op, alloc.String()), nil
		}
	} else if reg, ok := allocation.IsAllocReg(alloc); ok {
		// Perform operation
		if srcReg, ok := allocation.IsAllocReg(srcAlloc); ok {
			// REG = Op(REG)
			ret = append(ret, x86.DoUnaryOp(
				reg.Reg,
				srcReg.Reg,
				op,
				reg.Size,
			))
			return nil, ret
		} else if srcMem, ok := allocation.IsAllocMem(srcAlloc); ok {
			// REG = Op(MEM)
			ret = append(ret, x86.DoUnaryOp(
				reg.Reg,
				x86.GetMemoryVarLocation(srcMem.Index, srcMem.Size),
				op,
				reg.Size,
			))
			return nil, ret
		} else {
			return fmt.Errorf("Unsupported operation %v transfer to %s", op, alloc.String()), nil
		}
	} else {
		return fmt.Errorf("Unsupported operation %v transfer to %s", op, alloc.String()), nil
	}
}

func (backend CompilerX86Backend) compileIROpBinary(ret []*x86.Instruction, instr *ir.IRExpression, op ir.IROperator, name string, alloc ir.IRAllocation, srcAlloc1 ir.IRAllocation, srcAlloc2 ir.IRAllocation) (error, []*x86.Instruction) {
	if mem, ok := allocation.IsAllocMem(alloc); ok {
		if srcReg1, ok := allocation.IsAllocReg(srcAlloc1); ok {
			if srcReg2, ok := allocation.IsAllocReg(srcAlloc2); ok {
				ret = append(ret, x86.DoCompare(
					srcReg1.Reg,
					srcReg2.Reg,
				))
				ret = append(ret, x86.DoMemorySetConditional(
					mem.Index,
					mem.Size,
					op,
				))
				return nil, ret
			} else {
				return fmt.Errorf("Unsupported opertion %v transfer to memory from %s (arg 2)", op, srcAlloc2.String()), nil
			}
		} else {
			return fmt.Errorf("Unsupported opertion %v transfer to memory from %s (arg 1)", op, srcAlloc1.String()), nil
		}
	} else if reg, ok := allocation.IsAllocReg(alloc); ok {
		if srcReg1, ok := allocation.IsAllocReg(srcAlloc1); ok {
			if srcReg2, ok := allocation.IsAllocReg(srcAlloc2); ok {
				ret = append(ret, x86.DoCompare(
					srcReg1.Reg,
					srcReg2.Reg,
				))
				ret = append(ret, x86.DoRegSetConditional(
					reg.Reg,
					reg.State.SubRegSize1,
					reg.Size,
					op,
				)...)
				return nil, ret
			} else {
				return fmt.Errorf("Unsupported opertion %v transfer to memory from %s (arg 2)", op, srcAlloc2.String()), nil
			}
		} else {
			return fmt.Errorf("Unsupported opertion %v transfer to memory from %s (arg 1)", op, srcAlloc1.String()), nil
		}
	} else {
		return fmt.Errorf("Unsupported opertion %v transfer to memory from %s (arg 1)", op, srcAlloc1.String()), nil
	}
}

func (backend CompilerX86Backend) compileIRValuedExit(ret []*x86.Instruction, isMain bool, instr *ir.IRExit, alloc ir.IRAllocation) (error, []*x86.Instruction) {
	if mem, ok := allocation.IsAllocMem(alloc); ok {
		ret = append(ret, x86.DoMemoryLoad(
			mem.Index,
			mem.Size,
			x86.EAX,
		))
		if isMain {
			ret = append(ret, x86.DoRetMain()...)
		} else {
			ret = append(ret, x86.DoRet()...)
		}
		return nil, ret
	} else if reg, ok := allocation.IsAllocReg(alloc); ok {
		if reg.Reg != x86.EAX {
			// Return value from reg other than EAX so we move the value
			ret = append(ret, x86.DoMov(
				x86.EAX,
				reg.Reg,
				reg.Size,
			))
		}
		if isMain {
			ret = append(ret, x86.DoRetMain()...)
		} else {
			ret = append(ret, x86.DoRet()...)
		}
		return nil, ret
	} else {
		return fmt.Errorf("Unsupported exit with variable location %s", alloc.String()), nil
	}
}

func (backend CompilerX86Backend) compileIREmptyExit(ret []*x86.Instruction, instr *ir.IRExit) (error, []*x86.Instruction) {
	ret = append(ret, x86.DoMov(
		x86.EAX,
		x86.Imm(0),
		4,
	))
	ret = append(ret, x86.DoRet()...)
	return nil, ret
}

func (backend CompilerX86Backend) compileIRCall(ret []*x86.Instruction, fnName string, instr *ir.IRCall, name string, alloc ir.IRAllocation, argsOrder []string, argsAllocs map[string]ir.IRAllocation) (error, []*x86.Instruction) {
	// for _, argName := range argsOrder {

	// }
	entireStackSize := 4 * 4
	for _, argName := range argsOrder {
		if mem, ok := allocation.IsAllocMem(argsAllocs[argName]); ok {
			entireStackSize += mem.Size
		} else if reg, ok := allocation.IsAllocReg(argsAllocs[argName]); ok {
			entireStackSize += reg.Size
		}
	}
	//ret = append(ret, x86.DoSub(x86.RSP, x86.Imm(entireStackSize), 4))
	ret = append(ret, x86.DoCall(fmt.Sprintf("_%s", instr.CallTarget)))
	return nil, ret
}

func (backend CompilerX86Backend) compileIRMacroCall(ret []*x86.Instruction, fnName string, stmt *ir.IRStatement, instr *ir.IRMacroCall, macroName string, macroData map[string]interface{}) (error, []*x86.Instruction) {
	if macroName == "PreserveFunctionRegs" {
		ret = append(ret, x86.DoPush(x86.RAX, 4))
		ret = append(ret, x86.DoPush(x86.RBX, 4))
		ret = append(ret, x86.DoPush(x86.RCX, 4))
		ret = append(ret, x86.DoPush(x86.RDX, 4))
		return nil, ret
	} else if macroName == "LoadFunctionArgument" {
		alloc := stmt.GetAllocationContext()[instr.Var]
		if reg, ok := allocation.IsAllocReg(alloc); ok {
			ret = append(ret, x86.DoPush(reg.Reg, reg.Size))
			return nil, ret
		} else {
			return fmt.Errorf("Cannot use non-registry parameter for the %s macro", macroName), nil
		}
	} else if macroName == "StoreFunctionResult" {
		allocTgt := stmt.GetAllocationContext()[*instr.TargetName]
		if mem, ok := allocation.IsAllocMem(allocTgt); ok {
			ret = append(ret, x86.DoMov(x86.GetMemoryVarLocation(mem.Index, mem.Size), x86.EAX, 4))
		} else if reg, ok := allocation.IsAllocReg(allocTgt); ok {
			ret = append(ret, x86.DoMov(reg.Reg, x86.EAX, 4))
		} else {
			return fmt.Errorf("Invalid target location for %s macro (%v)", macroName, allocTgt), nil
		}

		// Preserve regs
		ret = append(ret, x86.DoPop(x86.RDX, 4))
		ret = append(ret, x86.DoPop(x86.RCX, 4))
		ret = append(ret, x86.DoPop(x86.RBX, 4))
		ret = append(ret, x86.DoPop(x86.RAX, 4))
		return nil, ret
	} else {
		return fmt.Errorf("Unsupported macro was used: %s", macroName), nil
	}
}

func (backend CompilerX86Backend) compileIRIf(ret []*x86.Instruction, fnName string, instr *ir.IRIf, alloc ir.IRAllocation) (error, []*x86.Instruction) {
	if mem, ok := allocation.IsAllocMem(alloc); ok {
		ret = append(ret, x86.DoMemoryLoad(
			mem.Index,
			mem.Size,
			x86.EAX,
		))
		ret = append(ret, x86.DoIf(
			fmt.Sprintf("%s_block%d", fnName, instr.BlockThen),
			fmt.Sprintf("%s_block%d", fnName, instr.BlockElse),
			instr.HasElseBlock(),
			x86.EAX,
		)...)
		return nil, ret
	} else if reg, ok := allocation.IsAllocReg(alloc); ok {
		ret = append(ret, x86.DoIf(
			fmt.Sprintf("%s_block%d", fnName, instr.BlockThen),
			fmt.Sprintf("%s_block%d", fnName, instr.BlockElse),
			instr.HasElseBlock(),
			reg.Reg,
		)...)
		return nil, ret
	} else {
		return fmt.Errorf("Unsupported if with variable location %s", alloc.String()), nil
	}
}

func (backend CompilerX86Backend) compileIRBlock(c *context.ParsingContext, fn *ir.IRFunction, code *ir.IRBlock) (error, []*x86.Instruction) {
	fnName := fn.Name
	ret := []*x86.Instruction{
		x86.Label(fmt.Sprintf("%s_block%d", fnName, code.BlockID)),
	}
	for _, instr := range code.Statements {
		lastIndex := len(ret) - 1
		if instr.IsConst() {
			name, alloc := instr.GetAllocationTarget()
			err, newRet := backend.compileIRConst(ret, instr.Const, name, alloc)
			if err != nil {
				return err, nil
			}
			ret = newRet
		} else if instr.IsCopy() {
			name, alloc := instr.GetAllocationTarget()
			srcAlloc := instr.GetAllocationContext()[instr.Copy.Var]
			err, newRet := backend.compileIRCopy(ret, instr.Copy, name, alloc, srcAlloc)
			if err != nil {
				return err, nil
			}
			ret = newRet
		} else if instr.IsExpression() {
			exp := instr.Expression
			opSpecs := exp.OperatorSpecs()

			// Unary operator
			if opSpecs.ArgsCount == 1 {
				name, alloc := instr.GetAllocationTarget()
				srcAlloc := instr.GetAllocationContext()[exp.Arguments[0]]
				err, newRet := backend.compileIROpUnary(ret, exp, exp.Operation, name, alloc, srcAlloc)
				if err != nil {
					return err, nil
				}
				ret = newRet
			} else {
				// Binary operation
				name, alloc := instr.GetAllocationTarget()
				srcAlloc1 := instr.GetAllocationContext()[exp.Arguments[0]]
				srcAlloc2 := instr.GetAllocationContext()[exp.Arguments[1]]
				err, newRet := backend.compileIROpBinary(ret, exp, exp.Operation, name, alloc, srcAlloc1, srcAlloc2)
				if err != nil {
					return err, nil
				}
				ret = newRet
			}
		} else if instr.IsExit() {
			exit := instr.Exit
			if exit.HasValue() {
				alloc := instr.GetAllocationContext()[*exit.Value]
				err, newRet := backend.compileIRValuedExit(ret, fnName == "main", exit, alloc)
				if err != nil {
					return err, nil
				}
				ret = newRet
			} else {
				err, newRet := backend.compileIREmptyExit(ret, exit)
				if err != nil {
					return err, nil
				}
				ret = newRet
			}
		} else if instr.IsIf() {
			ifStmt := instr.If
			alloc := instr.GetAllocationContext()[ifStmt.Condition]
			err, newRet := backend.compileIRIf(ret, fnName, ifStmt, alloc)
			if err != nil {
				return err, nil
			}
			ret = newRet
		} else if instr.IsMacroCall() {
			macroCall := instr.MacroCall
			err, newRet := backend.compileIRMacroCall(ret, fnName, instr, macroCall, macroCall.MacroName, macroCall.Data)
			if err != nil {
				return err, nil
			}
			ret = newRet
		} else if instr.IsCall() {
			call := instr.Call
			name, alloc := instr.GetAllocationTarget()
			argsAllocs := map[string]ir.IRAllocation{}
			for _, name := range call.Arguments {
				argsAllocs[name] = instr.GetAllocationContext()[name]
			}

			err, newRet := backend.compileIRCall(ret, fnName, call, name, alloc, call.Arguments, argsAllocs)
			if err != nil {
				return err, nil
			}
			ret = newRet
		} else {
			return fmt.Errorf("Invalid IR code to preprocess: %s", instr.Print(c)), nil
		}
		for i, newInstr := range ret[lastIndex:] {
			newInstr.FromIR(instr)
			if i > 0 {
				newInstr.Comment = ""
			}
		}
	}
	return nil, ret
}
