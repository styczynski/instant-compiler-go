package jvm

import (
	"fmt"
	"strings"

	"github.com/styczynski/latte-compiler/src/compiler"
	"github.com/styczynski/latte-compiler/src/compiler/jvm/jasmine"
	"github.com/styczynski/latte-compiler/src/config"
	"github.com/styczynski/latte-compiler/src/flow_analysis"
	"github.com/styczynski/latte-compiler/src/parser/context"
)

func init() {
	config.RegisterEntityFactory(config.ENTITY_COMPILER_BACKEND, CompilerJVMBackendFactory{})
}

type CompilerJVMBackendFactory struct{}

func (CompilerJVMBackendFactory) CreateEntity(c config.EntityConfig) interface{} {
	return CreateCompilerJVMBackend()
}

func (CompilerJVMBackendFactory) Params(argSpec *config.EntityArgSpec) {

}

func (CompilerJVMBackendFactory) EntityName() string {
	return "jvm"
}

type CompilerJVMBackend struct {
	state *compiler.CompilerState
}

func (backend CompilerJVMBackend) RunCompiledCode(runContext compiler.CompiledCodeRunContext, c *context.ParsingContext) ([]string, *compiler.RunError) {
	callOut, err := runContext.Call("java", "rror", "-cp", "$OUTPUT_DIR", "$CLASS_NAME")
	if err != nil {
		return nil, err
	}

	out := []string{}
	for _, line := range callOut {
		if len(line) > 0 && !strings.Contains(line, "_JAVA_OPTIONS") {
			out = append(out, line)
		}
	}

	return out, nil
}

func (backend CompilerJVMBackend) Compile(program flow_analysis.LatteAnalyzedProgram, c *context.ParsingContext, b *compiler.BuildContext) compiler.LatteCompiledProgramPromiseChan {
	ret := make(chan compiler.LatteCompiledProgram)
	go func() {

		//ast := program.Program.AST()
		maxStack := int64(1)
		outputCode := []jasmine.JasmineInstruction{}

		// output := jasmine.JasmineProgram{
		// 	StackLimit:   maxStack,
		// 	LocalsLimit:  int64(backend.state.ScopeSize()),
		// 	Instructions: outputCode,
		// }

		className := b.GetVariable("INPUT_FILE_BASE")
		b.SetCompilerMeta("CLASS_NAME", className)

		output := jasmine.JasmineProgram{
			Instructions: []jasmine.JasmineInstruction{
				&jasmine.JasmineClass{
					Name:  fmt.Sprintf("public %s", className),
					Super: "java/lang/Object",
					Methods: []*jasmine.JasmineMethod{
						&jasmine.JasmineMethod{
							Name:        "<init>",
							Returns:     "V",
							StackLimit:  1,
							LocalsLimit: 1,
							Body: []jasmine.JasmineInstruction{
								&jasmine.JasmineReferenceLoad{
									Index: 0,
								},
								&jasmine.JasmineInvokeStatic{
									Target:  "java/lang/Object/<init>",
									Special: true,
									Return:  "V",
								},
								&jasmine.JasmineReturn{},
							},
						},
						&jasmine.JasmineMethod{
							Name:        "public static main",
							Args:        []string{"[Ljava/lang/String;"},
							Returns:     "V",
							StackLimit:  maxStack,
							LocalsLimit: int64(backend.state.ScopeSize()) + 1,
							Body:        append(outputCode, &jasmine.JasmineReturn{}),
						},
					},
				},
			},
		}
		output.Normalize()

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

		b.WriteBuildFile("code.jasmine", []byte(output.ProgramToText()))

		validationErr = b.Call("java", "rror", "-jar", "$ROOT/lib/jasmin.jar", "-d", "$BUILD_DIR/out", "$BUILD_DIR/code.jasmine")
		if validationErr != nil {
			ret <- compiler.LatteCompiledProgram{
				Program:          program,
				CompiledProgram:  &output,
				CompilationError: validationErr,
			}
			return
		}

		outputJVMBytecode := b.ReadBuildFile("out/%s.class", className)
		b.WriteOutput("JVM bytecode file", "class", outputJVMBytecode)
		b.WriteOutput("Jasmine source", "j", []byte(output.ProgramToText()))

		ret <- compiler.LatteCompiledProgram{
			Program:          program,
			CompiledProgram:  &output,
			CompilationError: validationErr,
		}
	}()

	return ret
}

func (CompilerJVMBackend) BackendName() string {
	return "JVM Jasmine backend"
}

func CreateCompilerJVMBackend() compiler.CompilerBackend {
	return CompilerJVMBackend{
		state: compiler.CreateCompilerState(),
	}
}
