package assembly

import (
	"fmt"
	"strings"

	"github.com/styczynski/latte-compiler/src/compiler"
	"github.com/styczynski/latte-compiler/src/compiler/assembly/allocation"
	"github.com/styczynski/latte-compiler/src/compiler/assembly/x86"
	"github.com/styczynski/latte-compiler/src/config"
	"github.com/styczynski/latte-compiler/src/flow_analysis"
	"github.com/styczynski/latte-compiler/src/ir"
	"github.com/styczynski/latte-compiler/src/parser/context"
)

func init() {
	config.RegisterEntityFactory(config.ENTITY_COMPILER_BACKEND, CompilerX86BackendFactory{})
}

type CompilerX86BackendFactory struct{}

func (CompilerX86BackendFactory) CreateEntity(c config.EntityConfig) interface{} {
	return CreateCompilerX86Backend(
		c.Bool("gcc-docker"),
	)
}

func (CompilerX86BackendFactory) Params(argSpec *config.EntityArgSpec) {
	argSpec.AddBool("gcc-docker", false, "Prefer using dockerized GCC x64")
}

func (CompilerX86BackendFactory) EntityName() string {
	return "x86"
}

type CompilerX86Backend struct {
	state           *compiler.CompilerState
	preferGCCDocker bool
}

func (backend CompilerX86Backend) RunCompiledCode(runContext compiler.CompiledCodeRunContext, c *context.ParsingContext) ([]string, *compiler.RunError) {
	//return []string{}, nil

	runCmd := "bash"
	runArgs := []interface{}{
		"-c",
		"chmod +x ./$INPUT_FILE_BASE && ./$INPUT_FILE_BASE",
	}
	//(echo \"Output:\" ; ./$INPUT_FILE_BASE ; echo \"Exit code: $?\" ; true)
	if backend.preferGCCDocker {
		runCmd = "docker"
		runArgs = []interface{}{
			"run",
			"-t",
			"-v", "$OUTPUT_DIR:/code",
			"-w", "/code",
			"gcc:11.2.0",
			"bash", "-c", "chmod +x ./$INPUT_FILE_BASE && ./$INPUT_FILE_BASE",
		}
	}

	callOut, err := runContext.Call(runCmd, "rror", runArgs...)
	if err != nil {
		return nil, err
	}

	out := []string{}
	for _, line := range callOut {
		// && !strings.Contains(line, "_JAVA_OPTIONS")
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			out = append(out, line)
		}
	}
	return out, nil
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

		// Initial processing
		for _, fn := range program.IR.Statements {
			for _, fnBlock := range fn.FunctionBody {
				for _, stmt := range fnBlock.Statements {
					if stmt.IsMacroCall() {
						macroCall := stmt.MacroCall
						argNo := macroCall.Data["ArgNo"].(int)
						if macroCall.MacroName == "LoadInputFunctionArgument" {
							argReg := x86.GetRegisterForFunctionArg(argNo)
							if argReg == nil {
								// Do not allocate anything
								pos, _ := x86.GetMemoryForFunctionArg(argNo)
								stmt.SetTargetAllocationConstraints(*macroCall.TargetName, ir.IRAllocationConstraints{
									&allocation.AllocConsRequireMemoryStackTop{
										Offset: pos + 16,
									},
								})
							} else {
								stmt.SetTargetAllocationConstraints(*macroCall.TargetName, ir.IRAllocationConstraints{
									&allocation.AllocConsRequireSpecificRegisters{
										AllowedRegisters: []x86.Reg{
											argReg.Normalized,
										},
									},
								})
							}
						}
					}
				}
			}
		}

		var alloc allocation.Allocator = (&allocation.LinearScanAllocator{}).Lock([]x86.Reg{
			x86.EAX,
			x86.EBX,
		})
		allocation.RunAllocator(program.IR, alloc)
		fmt.Printf("?BASED?\n%s", program.IR.Print(c))

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

		fmt.Printf("\nAFTER PREPROCESSING STEP:\n%s\n==========END=========\n\n", program.IR.Print(c))

		err, entries := backend.compileIR(c, program.IR)
		if err != nil {
			ret <- compiler.LatteCompiledProgram{
				Program:          program,
				CompiledProgram:  nil,
				CompilationError: compiler.CreateCompilationError("Code emitter error", err.Error()),
			}
			return
		}

		for i := 0; i < 5; i++ {
			entries = x86.Optimize(entries)
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

		if 0 == 0 {
			assemblyGccCmd := "gcc"
			assemblyGccArgs := []interface{}{
				"-c", "$BUILD_DIR/code.s", "-o", "$BUILD_DIR/code.o",
			}
			if backend.preferGCCDocker {
				assemblyGccCmd = "docker"
				assemblyGccArgs = []interface{}{
					"run", "-v", "$BUILD_DIR:/code", "-w", "/code", "gcc:11.2.0", "gcc", "-c", "/code/code.s", "-o", "/code/code.o",
				}
			}
			validationErr = b.Call(assemblyGccCmd, "rror", assemblyGccArgs...)
			if validationErr != nil {
				ret <- compiler.LatteCompiledProgram{
					Program:          program,
					CompiledProgram:  &output,
					CompilationError: validationErr,
				}
				return
			}

			linkerGccCmd := "gcc"
			linkerGccArgs := []interface{}{
				"$BUILD_DIR/code.o", "-o", "$BUILD_DIR/code_exe",
			}
			if backend.preferGCCDocker {
				linkerGccCmd = "docker"
				linkerGccArgs = []interface{}{
					"run", "-v", "$BUILD_DIR:/code", "-w", "/code", "gcc:11.2.0", "gcc", "/code/code.o", "-o", "/code/code_exe",
				}
			}
			validationErr = b.Call(linkerGccCmd, "rror", linkerGccArgs...)

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
		}

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

func CreateCompilerX86Backend(preferGCCDocker bool) compiler.CompilerBackend {
	return CompilerX86Backend{
		state:           compiler.CreateCompilerState(),
		preferGCCDocker: preferGCCDocker,
	}
}
