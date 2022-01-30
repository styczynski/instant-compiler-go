package assembly

import (
	"github.com/styczynski/latte-compiler/src/compiler/assembly/allocation"
	"github.com/styczynski/latte-compiler/src/compiler/assembly/x86"
	"github.com/styczynski/latte-compiler/src/ir"
	"github.com/styczynski/latte-compiler/src/parser/context"
)

func (CompilerX86Backend) fnHeader(fn *ir.IRFunction, compiledFunction []*x86.Instruction) []*x86.Instruction {
	fnMeta := fn.GetMeta().(allocation.AssemblyFunctionMeta)
	if fnMeta.VarLen > 0 {
		ret := []*x86.Instruction{
			x86.DoPush(x86.RBP, 4),
			x86.DoRegistryCopy(x86.RBP, x86.RSP, 8),
			x86.DoSub(x86.RSP, x86.Imm(fnMeta.VarLen), 4),
		}
		ret = append(ret, compiledFunction...)
		return ret
	} else {
		// We must remove all leaves from this function and shift all non-negative movs from memory by 8 bytes
		newBody := []*x86.Instruction{}
		for _, instr := range compiledFunction {
			if !instr.IsLabel() {
				if instr.Op == x86.LEAVE {
					continue
				}
				for i, arg := range instr.Args {
					if argMem, ok := arg.(x86.Mem); ok {
						if argMem.Base == x86.RSP && argMem.Disp >= 0 {
							argMem.Disp -= 8
						}
						arg = argMem
					}
					instr.Args[i] = arg
				}
			}
			newBody = append(newBody, instr)
		}
		return newBody
	}
}

func (backend CompilerX86Backend) compileIR(c *context.ParsingContext, code *ir.IRProgram) (error, []x86.Entry) {
	data := EmptyAssemblyDataSection()
	ret := []x86.Entry{
		data,
	}
	for _, fn := range code.Statements {
		fnName := fn.Name
		//for _, fnBlock := range fn.FunctionBody {
		// Compile function
		// CDECL header
		bodyExprs := []*x86.Instruction{}
		// // Function body
		for _, fnBlock := range fn.FunctionBody {
			err, instrs := backend.compileIRBlock(c, fn, fnBlock, data)
			if err != nil {
				return err, nil
			}
			bodyExprs = append(bodyExprs, instrs...)
		}

		allBodyExprs := backend.fnHeader(fn, bodyExprs)

		ret = append(ret, &x86.Function{
			Name:   fnName,
			Source: fn.BaseASTNode,
			Body:   allBodyExprs,
		})

	}
	return nil, ret
}
