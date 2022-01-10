package assembly

import (
	"fmt"
	"reflect"

	"github.com/styczynski/latte-compiler/src/compiler"
	"github.com/styczynski/latte-compiler/src/compiler/assembly/x86"
	"github.com/styczynski/latte-compiler/src/config"
	"github.com/styczynski/latte-compiler/src/flow_analysis"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"golang.org/x/arch/x86/x86asm"
)

func init() {
	config.RegisterEntityFactory(config.ENTITY_COMPILER_BACKEND, CompilerX86BackendFactory{})
}

type CompilerX86BackendFactory struct{}

func (CompilerX86BackendFactory) CreateEntity(c config.EntityConfig) interface{} {
	return CreateCompilerX86Backend()
}

func (CompilerX86BackendFactory) Params(argSpec *config.EntityArgSpec) {

}

func (CompilerX86BackendFactory) EntityName() string {
	return "x86"
}

type CompilerX86Backend struct {
	state *compiler.CompilerState
}

func (backend CompilerX86Backend) compileExpression(expr generic_ast.Expression) []x86.Entry {
	if _, ok := (expr.(*ast.Empty)); ok {
		return []x86.Entry{}
	}
	if expr, ok := (expr.(*ast.LatteProgram)); ok {
		ret := []x86.Entry{}
		for _, stmt := range expr.Definitions {
			compiledValue := backend.compileExpression(stmt)
			ret = append(ret, compiledValue...)
		}
		return ret
	}
	if block, ok := (expr.(*ast.Block)); ok {
		blockExprs := []x86.Entry{}
		for _, expr := range block.Expressions() {
			fmt.Printf("COMPILE BODY EXPR: %s\n", expr)
			blockExprs = append(blockExprs, backend.compileExpression(expr)...)
		}
		return blockExprs
	}
	if topDef, ok := (expr.(*ast.TopDef)); ok {
		if topDef.IsFunction() {
			// Compile function
			// CDECL header
			bodyExprs := []*x86.Instruction{
				x86.DoPush(x86asm.EBP),
				x86.DoMov(x86asm.EBP, x86asm.ESP),
			}
			// Function body
			for _, entry := range backend.compileExpression(topDef.Function.Body()) {
				bodyExprs = append(bodyExprs, entry.(*x86.Instruction))
			}
			return []x86.Entry{
				&x86.Function{
					Name:   topDef.Function.Name,
					Type:   topDef.Function.GetDeclarationType().Concrete(),
					Source: topDef.BaseASTNode,
					Body:   bodyExprs,
				},
			}
		}
	}
	if _, ok := (expr.(*ast.Return)); ok {
		inst := x86asm.Inst{}
		inst.Op = x86asm.RET

		return []x86.Entry{
			x86.DoPop(x86asm.EBP),
			x86.DoRet(),
		}
	}
	if stmt, ok := (expr.(*ast.Statement)); ok {
		if stmt.IsReturn() {
			return backend.compileExpression(stmt.Return)
		}
		return []x86.Entry{}
	}
	panic(fmt.Sprintf("Invalid instruction given to compileExpression(): %s", reflect.TypeOf(expr)))
}

func (backend CompilerX86Backend) RunCompiledCode(runContext compiler.CompiledCodeRunContext, c *context.ParsingContext) ([]string, *compiler.RunError) {
	return []string{}, nil
}

func (backend CompilerX86Backend) Compile(program flow_analysis.LatteAnalyzedProgram, c *context.ParsingContext, b *compiler.BuildContext) compiler.LatteCompiledProgramPromiseChan {
	ret := make(chan compiler.LatteCompiledProgram)
	go func() {

		//ast := program.Program.AST()

		// output := x86.JasmineProgram{
		// 	StackLimit:   maxStack,
		// 	LocalsLimit:  int64(backend.state.ScopeSize()),
		// 	Instructions: outputCode,
		// }

		output := x86.Program{
			Entries: backend.compileExpression(program.Program.Program.AST()),
		}

		var validationErr *compiler.CompilationError
		// validationErr := output.Validate()

		// if validationErr != nil {
		// 	ret <- compiler.LatteCompiledProgram{
		// 		Program:          program,
		// 		CompiledProgram:  &output,
		// 		CompilationError: validationErr,
		// 	}
		// 	return
		// }

		// b.WriteBuildFile("code.jasmine", []byte(output.ProgramToText()))

		// validationErr = b.Call("java", "rror", "-jar", "$ROOT/lib/jasmin.jar", "-d", "$BUILD_DIR/out", "$BUILD_DIR/code.jasmine")
		// if validationErr != nil {
		// 	ret <- compiler.LatteCompiledProgram{
		// 		Program:          program,
		// 		CompiledProgram:  &output,
		// 		CompilationError: validationErr,
		// 	}
		// 	return
		// }

		// outputX86Bytecode := b.ReadBuildFile("out/%s.class", className)
		b.WriteOutput("X86 assembly source", "asm", []byte(output.ProgramToText()))

		ret <- compiler.LatteCompiledProgram{
			Program:          program,
			CompiledProgram:  &output,
			CompilationError: validationErr,
		}
	}()

	return ret
}

func (CompilerX86Backend) BackendName() string {
	return "X86 Jasmine backend"
}

func CreateCompilerX86Backend() compiler.CompilerBackend {
	return CompilerX86Backend{
		state: compiler.CreateCompilerState(),
	}
}
