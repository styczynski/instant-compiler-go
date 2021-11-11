package runner

import (
	"fmt"
	"strings"

	"github.com/styczynski/latte-compiler/src/compiler"
	"github.com/styczynski/latte-compiler/src/config"
	"github.com/styczynski/latte-compiler/src/errors"

	"github.com/styczynski/latte-compiler/src/parser/context"
)

func init() {
	config.RegisterEntityFactory(config.ENTITY_RUNNER, LatteCompiledCodeRunnerFactory{})
}

type LatteCompiledCodeRunnerFactory struct{}

func (LatteCompiledCodeRunnerFactory) CreateEntity(c config.EntityConfig) interface{} {
	return CreateLatteCompiledCodeRunner(c.String("test-extension"), c.Bool("always-run"))
}

func (LatteCompiledCodeRunnerFactory) Params(argSpec *config.EntityArgSpec) {
	argSpec.AddString("test-extension", "output", "Specify test file output extension")
	argSpec.AddBool("always-run", false, "Always run generated code")
}

func (LatteCompiledCodeRunnerFactory) EntityName() string {
	return "runner"
}

type LatteCompiledCodeRunner struct {
	testExtension string
	alwaysRun     bool
}

func CreateLatteCompiledCodeRunner(testExtension string, alwaysRun bool) *LatteCompiledCodeRunner {
	return &LatteCompiledCodeRunner{
		testExtension: testExtension,
		alwaysRun:     alwaysRun,
	}
}

type LatteRunnedProgram struct {
	Program       compiler.LatteCompiledProgram
	filename      string
	ProgramOutput []string
	RunError      *compiler.RunError
}

func (p LatteRunnedProgram) Filename() string {
	return p.filename
}

func (p LatteRunnedProgram) Resolve() LatteRunnedProgram {
	return p
}

type LatteRunnedProgramPromise interface {
	Resolve() LatteRunnedProgram
}

type LatteRunnedProgramPromiseChan <-chan LatteRunnedProgram

func (p LatteRunnedProgramPromiseChan) Resolve() LatteRunnedProgram {
	return <-p
}

func (tc *LatteCompiledCodeRunner) checkAsync(programPromise compiler.LatteCompiledProgramPromise, c *context.ParsingContext) LatteRunnedProgramPromise {
	r := make(chan LatteRunnedProgram)
	ctx := c.Copy()
	go func() {
		program := programPromise.Resolve()
		defer errors.GeneralRecovery(ctx, "Running compiled program", program.Filename(), func(message string, textMessage string) {
			r <- LatteRunnedProgram{
				Program:  program,
				filename: program.Filename(),
			}
		}, func() {
			close(r)
		})

		if program.Program.FlowAnalysisError != nil {
			r <- LatteRunnedProgram{
				Program:  program,
				filename: program.Filename(),
			}
			return
		}

		if program.CompilationError != nil {
			r <- LatteRunnedProgram{
				Program:  program,
				filename: program.Filename(),
			}
			return
		}

		if program.Program.Program.TypeCheckingError != nil {
			r <- LatteRunnedProgram{
				Program:  program,
				filename: program.Filename(),
			}
			return
		}

		if program.Program.Program.Program.ParsingError() != nil {
			r <- LatteRunnedProgram{
				Program:  program,
				filename: program.Filename(),
			}
			return
		}

		backendProcessDescription := fmt.Sprintf("Running compiled program: %v", tc.alwaysRun)

		c.EventsCollectorStream.Start(backendProcessDescription, c, program)
		defer c.EventsCollectorStream.End(backendProcessDescription, c, program)

		runContext := CreateCompiledCodeRunContext(program)

		expectedContentBytes, outFileErr := runContext.ReadFileByExt(tc.testExtension)
		if outFileErr == nil || tc.alwaysRun {
			runContext := CreateCompiledCodeRunContext(program)
			out, err := program.Backend.RunCompiledCode(runContext, c)

			if err != nil {
				r <- LatteRunnedProgram{
					Program:  program,
					filename: program.Filename(),
					RunError: err,
				}
				return
			}

			if outFileErr == nil {
				expectedContent := strings.Split(string(expectedContentBytes), "\n")
				testDescription := runContext.Substitute("Test $INPUT_FILE_BASE.%s", tc.testExtension)
				c.EventsCollectorStream.Start(testDescription, c, program)
				for lineNo, line := range out {
					if line != expectedContent[lineNo] {
						c.EventsCollectorStream.End(testDescription, c, program)
						r <- LatteRunnedProgram{
							Program:       program,
							filename:      program.Filename(),
							ProgramOutput: out,
							RunError:      compiler.CreateRunError("Failed test", runContext.Substitute("Test $INPUT_FILE_BASE.out in directory $INPUT_FILE_LOC has failed.\n    Line:    %d\n    Expected: %s\n    Got:     %s", lineNo, expectedContent[lineNo], line)),
						}
						return
					}
				}
				c.EventsCollectorStream.End(testDescription, c, program)
			}

			r <- LatteRunnedProgram{
				Program:       program,
				filename:      program.Filename(),
				ProgramOutput: out,
				RunError:      nil,
			}
			return
		}

		r <- LatteRunnedProgram{
			Program:  program,
			filename: program.Filename(),
			RunError: nil,
		}
	}()

	return LatteRunnedProgramPromiseChan(r)
}

func (tc *LatteCompiledCodeRunner) RunCompiledProgram(programs []compiler.LatteCompiledProgramPromise, c *context.ParsingContext) []LatteRunnedProgramPromise {
	ret := []LatteRunnedProgramPromise{}
	for _, programPromise := range programs {
		ret = append(ret, tc.checkAsync(programPromise, c))
	}
	return ret
}
