package llvm

import (
	"fmt"
	"reflect"

	"github.com/styczynski/latte-compiler/src/compiler"
	"github.com/styczynski/latte-compiler/src/compiler/llvm/llvm_ast"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker"
)

type CompilerLLVMBackend struct {
	state *compiler.CompilerState
}

func (backend CompilerLLVMBackend) RunCompiledCode(runContext compiler.CompiledCodeRunContext, c *context.ParsingContext) ([]string, *compiler.RunError) {
	return nil, nil
}

func (backend CompilerLLVMBackend) getLastInstrTarget(instr []llvm_ast.LLVMInstruction, withType bool) ([]llvm_ast.LLVMInstruction, string) {
	lastInstr := instr[len(instr)-1]
	if tgtOp, ok := lastInstr.(llvm_ast.LLVMTargetableInstruction); ok {
		if tgtOp.IsMovable() {
			instr = instr[:len(instr)-1]
		}
		return instr, tgtOp.GetTarget(withType)
	}
	panic(fmt.Sprintf("getLastInstrTarget() last instruction must be implementing llvm_ast.LLVMTargetableInstruction: %v", lastInstr))
}

func (backend CompilerLLVMBackend) compileExpression(expr generic_ast.Expression) ([]llvm_ast.LLVMInstruction, int64) {
	if expr, ok := (expr.(*ast.LatteProgram)); ok {
		ret := []llvm_ast.LLVMInstruction{}
		maxDepth := int64(0)
		for _, stmt := range expr.Definitions {
			compiledValue, s := backend.compileExpression(stmt)
			if s > maxDepth {
				maxDepth = s
			}
			ret = append(ret, compiledValue...)
		}
		return ret, int64(maxDepth)
	}
	ret := []llvm_ast.LLVMInstruction{}
	if expr, ok := (expr.(*ast.Statement)); ok {
		if expr.IsAssignment() {
			compiledValue, s := backend.compileExpression(&expr.Assignment.Value.Addition)
			_, v := backend.state.DefineAndAlloc(expr.Assignment.TargetName)

			//_, loc := backend.state.DefineAndAlloc(expr.Assignment.TargetName)

			lastInstr := compiledValue[len(compiledValue)-1]
			if tgtOp, ok := lastInstr.(llvm_ast.LLVMTargetableInstruction); ok {
				if tgtOp.IsMovable() {
					ret = append(ret, compiledValue[:len(compiledValue)-1]...)
					ret = append(ret, llvm_ast.CreateLLVMIntOp(fmt.Sprintf("%d", v), "i32", "0", "+", tgtOp.GetTarget(false)))
				} else {
					tgtOp.MoveTarget(fmt.Sprintf("%d", v))
					ret = append(ret, compiledValue...)
				}
			} else {
				panic(fmt.Sprintf("Cannot move instruction for assignment: %v", reflect.TypeOf(lastInstr)))
			}

			//TODO: WHAT?
			// ret = append(ret, &llvm_ast.LLVMStoreInt{
			// 	Index: loc,
			// })
			return ret, s
		} else if expr.IsExpression() {
			compiledValue, s := backend.compileExpression(expr.Expression)
			ret = append(ret, compiledValue...)
			return ret, s + 1
		}
	}
	if expr, ok := (expr.(*ast.Addition)); ok {
		if !expr.HasNext() {
			return backend.compileExpression(expr.Multiplication)
		}
		l, sl := backend.compileExpression(expr.Multiplication)
		r, sr := backend.compileExpression(expr.Next)

		_, v := backend.state.DefineAndAlloc(backend.state.NextUniqueVariableName())

		l, tl := backend.getLastInstrTarget(l, false)
		r, tr := backend.getLastInstrTarget(r, false)

		ret = append(ret, l...)
		ret = append(ret, r...)

		ret = append(ret, llvm_ast.CreateLLVMIntOp(
			fmt.Sprintf("%d", v),
			"i32",
			tl,
			expr.Op,
			tr,
		))
		return ret, sl + sr
	}
	if expr, ok := (expr.(*ast.Expression)); ok {
		a, sa := backend.compileExpression(&expr.Addition)
		a, ta := backend.getLastInstrTarget(a, true)
		ret = append(ret, a...)
		ret = append(ret, &llvm_ast.LLVMPrintInt{
			Target: ta,
		})
		return ret, sa + 1
	}
	if expr, ok := (expr.(*ast.Multiplication)); ok {
		if !expr.HasNext() {
			return backend.compileExpression(expr.Primary)
		}
		l, sl := backend.compileExpression(expr.Primary)
		r, sr := backend.compileExpression(expr.Next)

		l, tl := backend.getLastInstrTarget(l, false)
		r, tr := backend.getLastInstrTarget(r, false)

		_, v := backend.state.DefineAndAlloc(backend.state.NextUniqueVariableName())

		ret = append(ret, l...)
		ret = append(ret, r...)

		ret = append(ret, llvm_ast.CreateLLVMIntOp(
			fmt.Sprintf("%d", v),
			"i32",
			tl,
			expr.Op,
			tr,
		))
		return ret, sl + sr
	}
	if expr, ok := (expr.(*ast.Primary)); ok {
		if expr.IsVariable() {
			loc := backend.state.GetLocationFromScope(*expr.Variable)
			return []llvm_ast.LLVMInstruction{
				&llvm_ast.LLVMVar{
					ID: loc,
				},
			}, 1
		} else if expr.IsInt() {
			return []llvm_ast.LLVMInstruction{
				&llvm_ast.LLVMVal{
					ValueType: "i32",
					Value:     int(*expr.Int),
				},
			}, 1
		}
	}
	panic(fmt.Sprintf("Invalid instruction given to compileExpression(): %s", expr))
}

func (backend CompilerLLVMBackend) Compile(program type_checker.LatteTypecheckedProgram, c *context.ParsingContext, b *compiler.BuildContext) compiler.LatteCompiledProgramPromiseChan {
	ret := make(chan compiler.LatteCompiledProgram)
	go func() {

		ast := program.Program.AST()
		outputCode, _ := backend.compileExpression(ast)

		output := &llvm_ast.LLVMProgram{
			Instructions: outputCode,
		}

		b.WriteOutput("LLVM source", "ll", []byte(output.ProgramToText()))

		ret <- compiler.LatteCompiledProgram{
			Program:          program,
			CompiledProgram:  output,
			CompilationError: nil,
		}
	}()

	return ret
}

func (CompilerLLVMBackend) BackendName() string {
	return "LLVM backend"
}

func CreateCompilerLLVMBackend() compiler.CompilerBackend {
	return CompilerLLVMBackend{
		state: compiler.CreateCompilerState(),
	}
}
