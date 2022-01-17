package assembly

import (
	"github.com/styczynski/latte-compiler/src/compiler"
	"github.com/styczynski/latte-compiler/src/compiler/assembly/allocation"
	"github.com/styczynski/latte-compiler/src/compiler/assembly/x86"
	"github.com/styczynski/latte-compiler/src/config"
	"github.com/styczynski/latte-compiler/src/flow_analysis"
	"github.com/styczynski/latte-compiler/src/parser/context"
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

		var alloc allocation.Allocator = (&allocation.LinearScanAllocator{}).Lock([]x86.Reg{
			x86.EAX,
			x86.EBX,
		})
		allocation.RunAllocator(program.IR, alloc)
		// fmt.Printf("?BASED?\n%s", program.IR.Print(c))

		err := backend.preprocessIR(c, program.IR)
		if err != nil {
			ret <- compiler.LatteCompiledProgram{
				Program:          program,
				CompiledProgram:  nil,
				CompilationError: compiler.CreateCompilationError("IR Preprocessing error", err.Error()),
			}
			return
		}
		alloc.ResetSettings()
		allocation.RunAllocator(program.IR, alloc)

		// fmt.Printf("\nAFTER PREPROCESSING STEP:\n%s\n==========END=========\n\n", program.IR.Print(c))

		err, entries := backend.compileIR(c, program.IR)
		if err != nil {
			ret <- compiler.LatteCompiledProgram{
				Program:          program,
				CompiledProgram:  nil,
				CompilationError: compiler.CreateCompilationError("Code emitter error", err.Error()),
			}
			return
		}

		output := x86.Program{
			Entries: entries,
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

		b.WriteOutput("X86 assembly source", "s", []byte(output.ProgramToText()))
		b.WriteBuildFile("code.s", []byte(output.ProgramToText()))

		validationErr = b.Call("gcc", "rror", "-c", "$BUILD_DIR/code.s", "-o", "$BUILD_DIR/code.o")
		if validationErr != nil {
			ret <- compiler.LatteCompiledProgram{
				Program:          program,
				CompiledProgram:  &output,
				CompilationError: validationErr,
			}
			return
		}

		validationErr = b.Call("gcc", "rror", "$BUILD_DIR/code.o", "-o", "$BUILD_DIR/code_exe")
		if validationErr != nil {
			ret <- compiler.LatteCompiledProgram{
				Program:          program,
				CompiledProgram:  &output,
				CompilationError: validationErr,
			}
			return
		}

		outputX86Executable := b.ReadBuildFile("code_exe")

		b.WriteOutput("X86 program", "", outputX86Executable)

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
