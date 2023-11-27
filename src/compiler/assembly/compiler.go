package assembly

import (

	"github.com/styczynski/latte-compiler/src/compiler"
	"github.com/styczynski/latte-compiler/src/config"
	"github.com/styczynski/latte-compiler/src/flow_analysis"
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
	return []string{}, nil
}

func (backend CompilerX86Backend) Compile(program flow_analysis.LatteAnalyzedProgram, c *context.ParsingContext, b *compiler.BuildContext) compiler.LatteCompiledProgramPromiseChan {
	ret := make(chan compiler.LatteCompiledProgram)
	go func() {

		ret <- compiler.LatteCompiledProgram{
			Program:          program,
			CompiledProgram:  &compiler.CompiledProgramEmpty{},
			CompilationError: nil,
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
