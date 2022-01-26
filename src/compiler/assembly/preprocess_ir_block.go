package assembly

import (
	"fmt"

	"github.com/styczynski/latte-compiler/src/compiler/assembly/allocation"
	"github.com/styczynski/latte-compiler/src/ir"
	"github.com/styczynski/latte-compiler/src/parser/context"
)

func (backend CompilerX86Backend) preprocessIRConst(stmt *ir.IRStatement, instr *ir.IRConst, name string, alloc ir.IRAllocation) (error, []*ir.IRStatement) {
	return nil, []*ir.IRStatement{stmt}
}

func (backend CompilerX86Backend) preprocessIRCopy(stmt *ir.IRStatement, instr *ir.IRCopy, name string, alloc ir.IRAllocation, srcAlloc ir.IRAllocation) (error, []*ir.IRStatement) {
	if _, ok := allocation.IsAllocMem(alloc); ok {
		if _, ok := allocation.IsAllocMem(srcAlloc); ok {
			// Mem <-> Mem transfer
			newVarName := fmt.Sprintf("%s_regload_temp", instr.TargetName)
			return nil, []*ir.IRStatement{
				ir.WrapIRCopy(&ir.IRCopy{
					BaseASTNode: instr.BaseASTNode,
					TargetName:  newVarName,
					Type:        instr.Type,
					Var:         instr.Var,
					ParentNode:  instr.ParentNode,
				}).
					CopyDataForAllocationShadow(stmt).
					SetComment("Transfer %s to registry", instr.Var).
					SetTargetAllocationConstraints(newVarName, ir.IRAllocationConstraints{
						&allocation.AllocConsRequireRegisters{},
					}),
				ir.WrapIRCopy(&ir.IRCopy{
					BaseASTNode: instr.BaseASTNode,
					TargetName:  instr.TargetName,
					Type:        instr.Type,
					Var:         newVarName,
					ParentNode:  instr.ParentNode,
				}).CopyDataForAllocation(stmt),
			}
		} else if _, ok := allocation.IsAllocReg(srcAlloc); ok {
			return nil, []*ir.IRStatement{stmt}
		} else {
			return fmt.Errorf("Unsupported transfer to memory from %s", srcAlloc.String()), nil
		}
	} else if _, ok := allocation.IsAllocReg(alloc); ok {
		if _, ok := allocation.IsAllocMem(srcAlloc); ok {
			return nil, []*ir.IRStatement{stmt}
		} else if _, ok := allocation.IsAllocReg(srcAlloc); ok {
			return nil, []*ir.IRStatement{stmt}
		} else {
			return fmt.Errorf("Unsupported transfer to memory from %s", alloc.String()), nil
		}
	} else {
		return fmt.Errorf("Unsupported transfer to %s", alloc.String()), nil
	}
}

func (backend CompilerX86Backend) preprocessIROpUnary(stmt *ir.IRStatement, instr *ir.IRExpression, op ir.IROperator, name string, alloc ir.IRAllocation, srcAlloc ir.IRAllocation) (error, []*ir.IRStatement) {
	if _, ok := allocation.IsAllocMem(alloc); ok {
		newVarName := fmt.Sprintf("%s_regload_temp", instr.TargetName)
		return nil, []*ir.IRStatement{
			ir.WrapIRCopy(&ir.IRCopy{
				BaseASTNode: instr.BaseASTNode,
				TargetName:  newVarName,
				Type:        instr.Type,
				Var:         instr.TargetName,
			}).
				CopyDataForAllocationShadow(stmt).
				SetComment("Transfer %s to registry", instr.TargetName).
				SetTargetAllocationConstraints(newVarName, ir.IRAllocationConstraints{
					&allocation.AllocConsRequireRegisters{},
				}),
			ir.WrapIRExpression(&ir.IRExpression{
				BaseASTNode:    instr.BaseASTNode,
				TargetName:     newVarName,
				Type:           instr.Type,
				ArgumentsTypes: instr.ArgumentsTypes,
				Arguments:      instr.Arguments,
				Operation:      op,
				ParentNode:     instr.ParentNode,
			}).CopyDataForAllocation(stmt),
			ir.WrapIRCopy(&ir.IRCopy{
				BaseASTNode: instr.BaseASTNode,
				TargetName:  instr.TargetName,
				Type:        instr.Type,
				Var:         newVarName,
			}).
				CopyDataForAllocationShadow(stmt).
				SetComment("Transfer %s to registry", instr.TargetName).
				SetTargetAllocationConstraints(newVarName, ir.IRAllocationConstraints{
					&allocation.AllocConsAllowAll{},
				}),
		}

		// if _, ok := allocation.IsAllocMem(srcAlloc); ok {
		// 	// TODO: FILL THIS CRAP
		// 	return nil, []*ir.IRStatement{stmt}
		// } else {
		// 	return fmt.Errorf("Unsupported opertion %v transfer to memory from %s", op, srcAlloc.String()), nil
		// }
	} else if _, ok := allocation.IsAllocReg(alloc); ok {
		return nil, []*ir.IRStatement{stmt}
	} else {
		return fmt.Errorf("Unsupported operation %v transfer to %s", op, alloc.String()), nil
	}
}

func (backend CompilerX86Backend) preprocessIROpBinary(stmt *ir.IRStatement, instr *ir.IRExpression, op ir.IROperator, name string, alloc ir.IRAllocation, srcAlloc1 ir.IRAllocation, srcAlloc2 ir.IRAllocation) (error, []*ir.IRStatement) {
	if _, ok := allocation.IsAllocReg(srcAlloc1); ok {
		if _, ok := allocation.IsAllocReg(srcAlloc2); ok {
			// Op(REG, REG)
			return nil, []*ir.IRStatement{stmt}
		} else if _, ok := allocation.IsAllocMem(srcAlloc2); ok {
			// Op(REG, MEM)
			newVarName := fmt.Sprintf("%s_regload_temp", instr.Arguments[1])
			return nil, []*ir.IRStatement{
				ir.WrapIRCopy(&ir.IRCopy{
					BaseASTNode: instr.BaseASTNode,
					TargetName:  newVarName,
					Type:        instr.Type,
					Var:         instr.Arguments[1],
				}).
					CopyDataForAllocationShadow(stmt).
					SetComment("Transfer %s to registry", instr.Arguments[1]).
					SetTargetAllocationConstraints(newVarName, ir.IRAllocationConstraints{
						&allocation.AllocConsRequireRegisters{},
					}),
				ir.WrapIRExpression(&ir.IRExpression{
					BaseASTNode:    instr.BaseASTNode,
					TargetName:     instr.TargetName,
					Type:           instr.Type,
					ArgumentsTypes: instr.ArgumentsTypes,
					Arguments:      []string{instr.Arguments[0], newVarName},
					Operation:      op,
					ParentNode:     instr.ParentNode,
				}).CopyDataForAllocation(stmt),
			}
		} else {
			return fmt.Errorf("Unsupported operation on %s = %s(%s, %s)", op, alloc.String(), srcAlloc1.String(), srcAlloc2.String()), nil
		}
	} else if _, ok := allocation.IsAllocReg(srcAlloc1); ok {
		if _, ok := allocation.IsAllocReg(srcAlloc2); ok {
			// Op(MEM, REG)
			newVarName := fmt.Sprintf("%s_regload_temp", instr.Arguments[0])
			return nil, []*ir.IRStatement{
				ir.WrapIRCopy(&ir.IRCopy{
					BaseASTNode: instr.BaseASTNode,
					TargetName:  newVarName,
					Type:        instr.Type,
					Var:         instr.Arguments[0],
				}).
					CopyDataForAllocationShadow(stmt).
					SetComment("Transfer %s to registry", instr.Arguments[0]).
					SetTargetAllocationConstraints(newVarName, ir.IRAllocationConstraints{
						&allocation.AllocConsRequireRegisters{},
					}),
				ir.WrapIRExpression(&ir.IRExpression{
					BaseASTNode:    instr.BaseASTNode,
					TargetName:     instr.TargetName,
					Type:           instr.Type,
					ArgumentsTypes: instr.ArgumentsTypes,
					Arguments:      []string{newVarName, instr.Arguments[1]},
					Operation:      op,
					ParentNode:     instr.ParentNode,
				}).CopyDataForAllocation(stmt),
			}
		} else if _, ok := allocation.IsAllocMem(srcAlloc2); ok {
			// Op(MEM, MEM)
			newVarName0 := fmt.Sprintf("%s_regload_temp", instr.Arguments[0])
			newVarName1 := fmt.Sprintf("%s_regload_temp", instr.Arguments[1])
			return nil, []*ir.IRStatement{
				ir.WrapIRCopy(&ir.IRCopy{
					BaseASTNode: instr.BaseASTNode,
					TargetName:  newVarName1,
					Type:        instr.Type,
					Var:         instr.Arguments[1],
				}).
					CopyDataForAllocationShadow(stmt).
					SetComment("Transfer %s to registry", instr.Arguments[1]).
					SetTargetAllocationConstraints(newVarName1, ir.IRAllocationConstraints{
						&allocation.AllocConsRequireRegisters{},
					}),
				ir.WrapIRCopy(&ir.IRCopy{
					BaseASTNode: instr.BaseASTNode,
					TargetName:  newVarName0,
					Type:        instr.Type,
					Var:         instr.Arguments[0],
				}).
					CopyDataForAllocationShadow(stmt).
					SetComment("Transfer %s to registry", instr.Arguments[0]).
					SetTargetAllocationConstraints(newVarName1, ir.IRAllocationConstraints{
						&allocation.AllocConsRequireRegisters{},
					}),
				ir.WrapIRExpression(&ir.IRExpression{
					BaseASTNode:    instr.BaseASTNode,
					TargetName:     instr.TargetName,
					Type:           instr.Type,
					ArgumentsTypes: instr.ArgumentsTypes,
					Arguments:      []string{newVarName0, newVarName1},
					Operation:      op,
					ParentNode:     instr.ParentNode,
				}).CopyDataForAllocation(stmt),
			}
		} else {
			return fmt.Errorf("Unsupported operation on %s = %s(%s, %s)", op, alloc.String(), srcAlloc1.String(), srcAlloc2.String()), nil
		}
	} else {
		return fmt.Errorf("Unsupported operation on %s = %s(%s, %s)", op, alloc.String(), srcAlloc1.String(), srcAlloc2.String()), nil
	}
	// if _, ok := allocation.IsAllocMem(alloc); ok {
	// 	newVarName := fmt.Sprintf("%s_regload_temp", instr.TargetName)
	// 	return nil, []*ir.IRStatement{
	// 		ir.WrapIRCopy(&ir.IRCopy{
	// 			BaseASTNode: instr.BaseASTNode,
	// 			TargetName:  newVarName,
	// 			Type:        instr.Type,
	// 			Var:         instr.TargetName,
	// 		}).
	// 			CopyDataForAllocationShadow(stmt).
	// 			SetComment("Transfer %s to registry", instr.TargetName).
	// 			SetTargetAllocationConstraints(newVarName, ir.IRAllocationConstraints{
	// 				&allocation.AllocConsRequireRegisters{},
	// 			}),
	// 		ir.WrapIRExpression(&ir.IRExpression{
	// 			BaseASTNode:    instr.BaseASTNode,
	// 			TargetName:     instr.TargetName,
	// 			Type:           instr.Type,
	// 			ArgumentsTypes: instr.ArgumentsTypes,
	// 			Arguments:      []string{newVarName, instr.Arguments[1]},
	// 			Operation:      op,
	// 			ParentNode:     instr.ParentNode,
	// 		}).CopyDataForAllocation(stmt),
	// 	}
	// } else if _, ok := allocation.IsAllocReg(alloc); ok {
	// 	return nil, []*ir.IRStatement{stmt}
	// } else {
	// 	return fmt.Errorf("Unsupported operation %v transfer to %s", op, alloc.String()), nil
	// }
}

func (backend CompilerX86Backend) preprocessIRValuedExit(stmt *ir.IRStatement, instr *ir.IRExit, alloc ir.IRAllocation) (error, []*ir.IRStatement) {
	if _, ok := allocation.IsAllocMem(alloc); ok {
		return nil, []*ir.IRStatement{stmt}
	} else if _, ok := allocation.IsAllocReg(alloc); ok {
		return nil, []*ir.IRStatement{stmt}
	} else {
		return fmt.Errorf("Unsupported exit with variable location %s", alloc.String()), nil
	}
}

func (backend CompilerX86Backend) preprocessIREmptyExit(stmt *ir.IRStatement, instr *ir.IRExit) (error, []*ir.IRStatement) {
	return nil, []*ir.IRStatement{stmt}
}

func (backend CompilerX86Backend) preprocessIRJump(stmt *ir.IRStatement, fn *ir.IRFunction, instr *ir.IRJump) (error, []*ir.IRStatement) {
	return nil, []*ir.IRStatement{stmt}
}

func (backend CompilerX86Backend) preprocessIRIf(stmt *ir.IRStatement, fn *ir.IRFunction, instr *ir.IRIf, alloc ir.IRAllocation) (error, []*ir.IRStatement) {
	if _, ok := allocation.IsAllocMem(alloc); ok {
		return nil, []*ir.IRStatement{stmt}
	} else if _, ok := allocation.IsAllocReg(alloc); ok {
		return nil, []*ir.IRStatement{stmt}
	} else {
		return fmt.Errorf("Unsupported if with variable location %s", alloc.String()), nil
	}
}

func (backend CompilerX86Backend) preprocessIRCall(stmt *ir.IRStatement, fn *ir.IRFunction, instr *ir.IRCall, name string, alloc ir.IRAllocation, argsOrder []string, argsAllocs map[string]ir.IRAllocation) (error, []*ir.IRStatement) {

	ret := []*ir.IRStatement{}

	// ret = append(ret,
	// 	ir.WrapIRMacroCall(&ir.IRMacroCall{
	// 		BaseASTNode: instr.BaseASTNode,
	// 		Type:        instr.Type,
	// 		MacroName:   "PreserveFunctionRegs",
	// 		Var:         "",
	// 		ParentNode:  instr.ParentNode,
	// 	}).
	// 		CopyDataForAllocationShadow(stmt).
	// 		SetComment("Preserve all registries for %s call", instr.TargetName),
	// )

	// for argNo, argName := range argsOrder {
	// 	argAlloc := argsAllocs[argName]
	// 	if argMem, ok := allocation.IsAllocMem(argAlloc); ok {
	// 		newVarName := fmt.Sprintf("%s_call_arg_temp", argName)
	// 		return nil, []*ir.IRStatement{
	// 			ir.WrapIRCopy(&ir.IRCopy{
	// 				BaseASTNode: instr.BaseASTNode,
	// 				TargetName:  newVarName,
	// 				Type:        instr.Type,
	// 				Var:         argName,
	// 				ParentNode:  instr.ParentNode,
	// 			}).
	// 				CopyDataForAllocationShadow(stmt).
	// 				SetComment("Transfer call arg %s to registry", argName).
	// 				SetTargetAllocationConstraints(newVarName, ir.IRAllocationConstraints{
	// 					&allocation.AllocConsRequireRegisters{},
	// 				}),
	// 			ir.WrapIRMacroCall(&ir.IRMacroCall{
	// 				BaseASTNode: instr.BaseASTNode,
	// 				Type:        instr.Type,
	// 				MacroName:   "LoadFunctionArgument",
	// 				Data: map[string]interface{}{
	// 					"Input": argName,
	// 					"Size":  argMem.Size,
	// 					"ArgNo": argNo,
	// 				},
	// 				Var:        argName,
	// 				ParentNode: instr.ParentNode,
	// 			}).
	// 				CopyDataForAllocationShadow(stmt).
	// 				SetComment("Load function argument %s", argName),
	// 		}
	// 	} else if argReg, ok := allocation.IsAllocReg(argAlloc); ok {
	// 		ret = append(ret,
	// 			ir.WrapIRMacroCall(&ir.IRMacroCall{
	// 				BaseASTNode: instr.BaseASTNode,
	// 				Type:        instr.Type,
	// 				MacroName:   "LoadFunctionArgument",
	// 				Data: map[string]interface{}{
	// 					"Input": argName,
	// 					"Size":  argReg.Size,
	// 					"ArgNo": argNo,
	// 				},
	// 				Var:        argName,
	// 				ParentNode: instr.ParentNode,
	// 			}).
	// 				CopyDataForAllocationShadow(stmt).
	// 				SetComment("Load function argument %s", argName),
	// 		)
	// 	}
	// }

	// Original function call
	ret = append(ret, ir.WrapIRCall(&ir.IRCall{
		BaseASTNode:    instr.BaseASTNode,
		TargetName:     instr.TargetName,
		Type:           instr.Type,
		CallTarget:     instr.CallTarget,
		CallTargetType: instr.CallTargetType,
		ArgumentsTypes: instr.ArgumentsTypes,
		Arguments:      instr.Arguments,
		ParentNode:     instr.ParentNode,
		IsBuiltin:      instr.IsBuiltin,
	}).CopyDataForAllocation(stmt))

	//fnResultVarName := fmt.Sprintf("call_result_tmp_%s", instr.TargetName)
	// ret = append(ret,
	// 	ir.WrapIRMacroCall(&ir.IRMacroCall{
	// 		BaseASTNode: instr.BaseASTNode,
	// 		Type:        instr.Type,
	// 		MacroName:   "LoadFunctionResult",
	// 		Var:         "",
	// 		ParentNode:  instr.ParentNode,
	// 		TargetName:  &fnResultVarName,
	// 	}).
	// 		CopyDataForAllocationShadow(stmt).
	// 		SetComment("Load function result into %s", fnResultVarName).
	// 		SetTargetAllocationConstraints(fnResultVarName, ir.IRAllocationConstraints{
	// 			&allocation.AllocConsRequireRegisters{},
	// 		}),
	// )
	// ret = append(ret,
	// 	ir.WrapIRMacroCall(&ir.IRMacroCall{
	// 		BaseASTNode: instr.BaseASTNode,
	// 		Type:        instr.Type,
	// 		MacroName:   "StoreFunctionResult",
	// 		Data: map[string]interface{}{
	// 			"Alloc": alloc,
	// 		},
	// 		Var:        "",
	// 		TargetName: &instr.TargetName,
	// 		ParentNode: instr.ParentNode,
	// 	}).
	// 		CopyDataForAllocationShadow(stmt).
	// 		SetComment("Store function result into %s", instr.TargetName),
	// )
	// ret = append(ret,
	// 	ir.WrapIRMacroCall(&ir.IRMacroCall{
	// 		BaseASTNode: instr.BaseASTNode,
	// 		Type:        instr.Type,
	// 		MacroName:   "ResotreFunctionCall",
	// 		Data: map[string]interface{}{
	// 			"Regs": []x86.Reg{
	// 				x86.RAX,
	// 				x86.RBX,
	// 				x86.RCX,
	// 				x86.RDX,
	// 			},
	// 		},
	// 		Var:        "",
	// 		ParentNode: instr.ParentNode,
	// 	}).
	// 		CopyDataForAllocationShadow(stmt).
	// 		SetComment("Restore registries"),
	//)

	return nil, ret
}

func (backend CompilerX86Backend) preprocessIRBlock(c *context.ParsingContext, fn *ir.IRFunction, code *ir.IRBlock) error {
	ret := []*ir.IRStatement{}

	for _, instr := range code.Statements {
		if instr.IsConst() {
			name, alloc := instr.GetAllocationTarget()
			err, mappedInstrs := backend.preprocessIRConst(instr, instr.Const, name, alloc)
			if err != nil {
				return err
			}
			ret = append(ret, mappedInstrs...)
		} else if instr.IsCopy() {
			name, alloc := instr.GetAllocationTarget()
			srcAlloc := instr.GetAllocationContext()[instr.Copy.Var]
			err, mappedInstrs := backend.preprocessIRCopy(instr, instr.Copy, name, alloc, srcAlloc)
			if err != nil {
				return err
			}
			ret = append(ret, mappedInstrs...)
		} else if instr.IsExpression() {
			exp := instr.Expression
			opSpecs := exp.OperatorSpecs()

			// Unary operator
			if opSpecs.ArgsCount == 1 {
				name, alloc := instr.GetAllocationTarget()
				srcAlloc := instr.GetAllocationContext()[exp.Arguments[0]]
				err, mappedInstrs := backend.preprocessIROpUnary(instr, exp, exp.Operation, name, alloc, srcAlloc)
				if err != nil {
					return err
				}
				ret = append(ret, mappedInstrs...)
			} else {
				// Binary operation
				name, alloc := instr.GetAllocationTarget()
				srcAlloc1 := instr.GetAllocationContext()[exp.Arguments[0]]
				srcAlloc2 := instr.GetAllocationContext()[exp.Arguments[1]]
				err, mappedInstrs := backend.preprocessIROpBinary(instr, exp, exp.Operation, name, alloc, srcAlloc1, srcAlloc2)
				if err != nil {
					return err
				}
				ret = append(ret, mappedInstrs...)
			}
		} else if instr.IsExit() {
			exit := instr.Exit
			if exit.HasValue() {
				alloc := instr.GetAllocationContext()[*exit.Value]
				err, mappedInstrs := backend.preprocessIRValuedExit(instr, exit, alloc)
				if err != nil {
					return err
				}
				ret = append(ret, mappedInstrs...)
			} else {
				err, mappedInstrs := backend.preprocessIREmptyExit(instr, exit)
				if err != nil {
					return err
				}
				ret = append(ret, mappedInstrs...)
			}
		} else if instr.IsIf() {
			ifStmt := instr.If
			alloc := instr.GetAllocationContext()[ifStmt.Condition]
			err, mappedInstrs := backend.preprocessIRIf(instr, fn, ifStmt, alloc)
			if err != nil {
				return err
			}
			ret = append(ret, mappedInstrs...)
		} else if instr.IsJump() {
			jumpStmt := instr.Jump
			err, mappedInstrs := backend.preprocessIRJump(instr, fn, jumpStmt)
			if err != nil {
				return err
			}
			ret = append(ret, mappedInstrs...)
		} else if instr.IsCall() {
			call := instr.Call
			name, alloc := instr.GetAllocationTarget()
			argsAllocs := map[string]ir.IRAllocation{}
			for _, name := range call.Arguments {
				argsAllocs[name] = instr.GetAllocationContext()[name]
			}
			err, mappedInstrs := backend.preprocessIRCall(instr, fn, call, name, alloc, call.Arguments, argsAllocs)
			if err != nil {
				return err
			}
			ret = append(ret, mappedInstrs...)
		} else if instr.IsMacroCall() {
			ret = append(ret, instr)
		} else if instr.IsEmpty() {
			// Do nothing
		} else {
			return fmt.Errorf("Invalid IR code to preprocess: %s", instr.Print(c))
		}
	}
	code.Statements = ret
	return nil
}
