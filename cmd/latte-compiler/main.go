package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/styczynski/latte-compiler/cmd/latte-compiler/config"
	"github.com/styczynski/latte-compiler/src/compiler"
	"github.com/styczynski/latte-compiler/src/compiler/jvm"
	"github.com/styczynski/latte-compiler/src/compiler/llvm"
	"github.com/styczynski/latte-compiler/src/events_collector"
	"github.com/styczynski/latte-compiler/src/input_reader"
	"github.com/styczynski/latte-compiler/src/parser"
	context2 "github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/printer"
	"github.com/styczynski/latte-compiler/src/runner"
	"github.com/styczynski/latte-compiler/src/type_checker"
)

func ActionCompile(c *cli.Context) error {
	pr := printer.CreateLattePrinter()
	eventsCollector := events_collector.StartEventsCollector(events_collector.CreateStatusUpdater(false))
	context := context2.NewParsingContext(pr, eventsCollector)

	tc := type_checker.CreateLatteTypeChecker()
	p := parser.CreateLatteParser()

	inputPaths := c.Args().Slice()
	reader := input_reader.CreateLatteInputReader(inputPaths)

	backend := jvm.CreateCompilerJVMBackend()
	backendName := c.String("backend")
	if backendName == "llvm" {
		backend = llvm.CreateCompilerLLVMBackend()
	} else if backendName == "jvm" {
		backend = jvm.CreateCompilerJVMBackend()
	} else {
		log.Fatal(fmt.Sprintf("Invalid compiler backend was specified: %s", backendName))
	}

	comp := compiler.CreateLatteCompiler(backend)
	run := runner.CreateLatteCompiledCodeRunner()
	ast := p.ParseInput(reader, context)

	checkedProgram := tc.Check(ast, context)
	compiledProgram := comp.Compile(checkedProgram, context)
	runnedProgram := run.RunCompiledProgram(compiledProgram, context)

	var summary events_collector.Summarizer = events_collector.CreateCliSummary(-1)
	if c.Bool("short") {
		summary = events_collector.CreateCliSummaryShortStatus()
	}
	message, ok := eventsCollector.SummarizeCompiledCodeRunning(summary, runnedProgram, context)

	if !ok {
		os.Stderr.WriteString("ERROR\n")
		fmt.Print(message)
		os.Exit(1)
	} else {
		os.Stderr.WriteString("OK\n")
		fmt.Print(message)
	}

	//f, err := os.Create("compiler.prof")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//pprof.StartCPUProfile(f)
	//defer pprof.StopCPUProfile()

	return nil
}

func main() {
	// load application configurations
	if err := config.LoadConfig("./config"); err != nil {
		panic(fmt.Errorf("invalid application configuration: %s", err))
	}

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "short",
				Usage: "Display only basic success/error information (useful when testing files with glob expression or in bulk)",
				Value: false,
			},
			&cli.StringFlag{
				Name:  "backend",
				Usage: "Use specific compiler backend. Supported options are: jvm/llvm",
				Value: "jvm",
			},
		},
		Action: ActionCompile,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
