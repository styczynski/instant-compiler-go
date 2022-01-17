package assembly

import (
	"github.com/styczynski/latte-compiler/src/compiler/assembly/allocation"
	"github.com/styczynski/latte-compiler/src/compiler/assembly/x86"
	"github.com/styczynski/latte-compiler/src/ir"
	"github.com/styczynski/latte-compiler/src/parser/context"
)

func (CompilerX86Backend) fnHeader(fn *ir.IRFunction) []*x86.Instruction {
	fnMeta := fn.GetMeta().(allocation.AssemblyFunctionMeta)
	ret := []*x86.Instruction{
		x86.DoPush(x86.RBP, 4),
		x86.DoMov(x86.RBP, x86.RSP, 4),
	}
	if fnMeta.VarLen > 0 {
		ret = append(ret, x86.DoSub(x86.RSP, x86.Imm(fnMeta.VarLen), 4))
	}
	return ret
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
		bodyExprs := backend.fnHeader(fn)
		// // Function body
		for _, fnBlock := range fn.FunctionBody {
			err, instrs := backend.compileIRBlock(c, fn, fnBlock, data)
			if err != nil {
				return err, nil
			}
			bodyExprs = append(bodyExprs, instrs...)
		}
		ret = append(ret, &x86.Function{
			Name:   fnName,
			Source: fn.BaseASTNode,
			Body:   bodyExprs,
		})

	}
	return nil, ret
}
