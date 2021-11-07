package compiler_pipeline

import (
	"fmt"
	"os"
	"strings"

	"github.com/styczynski/latte-compiler/src/compiler"
	"github.com/styczynski/latte-compiler/src/config"
	"github.com/styczynski/latte-compiler/src/events_collector"
	"github.com/styczynski/latte-compiler/src/input_reader"
	"github.com/styczynski/latte-compiler/src/parser"
	context2 "github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/printer"
	"github.com/styczynski/latte-compiler/src/runner"
	"github.com/styczynski/latte-compiler/src/type_checker"
)

func init() {
	config.RegisterEntityFactory(config.ENTITY_COMPILER_PIPELINE, CompilerPipelineFactory{})
}

type CompilerPipelineFactory struct{}

func (CompilerPipelineFactory) CreateEntity(c config.EntityConfig) interface{} {
	return CreateCompilerPipeline()
}

func (CompilerPipelineFactory) Params(argSpec *config.EntityArgSpec) {
	argSpec.AddString("backend",
		config.GetEntityNamesList(config.ENTITY_COMPILER_BACKEND)[0],
		fmt.Sprintf("Use specific compiler backend. Supported options are: %s", strings.Join(config.GetEntityNamesList(config.ENTITY_COMPILER_BACKEND), ", ")))
	argSpec.AddString("summary",
		config.GetEntityNamesList(config.ENTITY_SUMMARIZER)[0],
		fmt.Sprintf("Use specific summarizer. Supported options are: %s", strings.Join(config.GetEntityNamesList(config.ENTITY_SUMMARIZER), ", ")))
	argSpec.AddString("status-updater",
		config.GetEntityNamesList(config.ENTITY_STATUS_UPDATER)[0],
		fmt.Sprintf("Use specific live status updater. Supported options are: %s", strings.Join(config.GetEntityNamesList(config.ENTITY_STATUS_UPDATER), ", ")))
}

func (CompilerPipelineFactory) EntityName() string {
	return "compiler-pipeline"
}

type CompilerPipeline struct{}

func CreateCompilerPipeline() CompilerPipeline {
	return CompilerPipeline{}
}

func (CompilerPipeline) RunPipeline(c config.EntityConfig, inputPaths []string) error {
	pr := printer.CreateLattePrinter()
	eventsCollector := events_collector.StartEventsCollector(
		config.CreateEntity(config.ENTITY_STATUS_UPDATER, c.String("status-updater"), c).(events_collector.StatusUpdater))
	context := context2.NewParsingContext(pr, eventsCollector)

	tc := type_checker.CreateLatteTypeChecker()
	p := parser.CreateLatteParser()

	reader := input_reader.CreateLatteInputReader(inputPaths)
	backend := config.CreateEntity(config.ENTITY_COMPILER_BACKEND, c.String("backend"), c).(compiler.CompilerBackend)

	comp := compiler.CreateLatteCompiler(backend)
	run := config.CreateEntity(config.ENTITY_RUNNER, "runner", c).(*runner.LatteCompiledCodeRunner)
	ast := p.ParseInput(reader, context)

	checkedProgram := tc.Check(ast, context)
	compiledProgram := comp.Compile(checkedProgram, context)
	runnedProgram := run.RunCompiledProgram(compiledProgram, context)

	var summary events_collector.Summarizer = config.CreateEntity(config.ENTITY_SUMMARIZER, c.String("summary"), c).(events_collector.Summarizer)
	message, ok := eventsCollector.SummarizeCompiledCodeRunning(summary, runnedProgram, context)

	if !ok {
		os.Stderr.WriteString("ERROR\n")
		fmt.Print(message)
		os.Exit(1)
	} else {
		os.Stderr.WriteString("OK\n")
		fmt.Print(message)
	}
	return nil
}
