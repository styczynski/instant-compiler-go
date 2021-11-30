package compiler_pipeline

import (
	"fmt"
	"strings"

	"github.com/styczynski/latte-compiler/src/compiler"
	"github.com/styczynski/latte-compiler/src/config"
	"github.com/styczynski/latte-compiler/src/events_collector"
	"github.com/styczynski/latte-compiler/src/flow_analysis"
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
		"jvm",
		fmt.Sprintf("Use specific compiler backend. Supported options are: %s", strings.Join(config.GetEntityNamesList(config.ENTITY_COMPILER_BACKEND), ", ")))
	argSpec.AddString("summary",
		"summary-cli",
		fmt.Sprintf("Use specific summarizer. Supported options are: %s", strings.Join(config.GetEntityNamesList(config.ENTITY_SUMMARIZER), ", ")))
	argSpec.AddString("status-updater",
		"updater-cli-progress",
		fmt.Sprintf("Use specific live status updater. Supported options are: %s", strings.Join(config.GetEntityNamesList(config.ENTITY_STATUS_UPDATER), ", ")))
}

func (CompilerPipelineFactory) EntityName() string {
	return "compiler-pipeline"
}

type CompilerPipeline struct{}

func CreateCompilerPipeline() CompilerPipeline {
	return CompilerPipeline{}
}

func (CompilerPipeline) RunPipeline(c config.EntityConfig, reader input_reader.InputReader) (string, []runner.LatteRunnedProgram, bool) {
	pr := printer.CreateLattePrinter()
	eventsCollector := events_collector.StartEventsCollector(
		config.CreateEntity(config.ENTITY_STATUS_UPDATER, c.String("status-updater"), c).(events_collector.StatusUpdater))
	context := context2.NewParsingContext(pr, eventsCollector)

	tc := type_checker.CreateLatteTypeChecker()
	p := parser.CreateLatteParser()
	flow_analyzer := flow_analysis.CreateLatteFlowAnalyzer()

	backend := config.CreateEntity(config.ENTITY_COMPILER_BACKEND, c.String("backend"), c).(compiler.CompilerBackend)

	comp := compiler.CreateLatteCompiler(backend)
	run := config.CreateEntity(config.ENTITY_RUNNER, "runner", c).(*runner.LatteCompiledCodeRunner)
	ast := p.ParseInput(reader, context)

	checkedProgram := tc.Check(ast, context)
	programFlow := flow_analyzer.Analyze(checkedProgram, context)
	compiledProgram := comp.Compile(programFlow, context)
	runnedProgram := run.RunCompiledProgram(compiledProgram, context)

	var summary events_collector.Summarizer = config.CreateEntity(config.ENTITY_SUMMARIZER, c.String("summary"), c).(events_collector.Summarizer)
	message, progs, ok := eventsCollector.SummarizeCompiledCodeRunning(summary, runnedProgram, context)

	return message, progs, ok
}
