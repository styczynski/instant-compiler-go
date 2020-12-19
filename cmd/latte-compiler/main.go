package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/styczynski/latte-compiler/cmd/latte-compiler/config"
	"github.com/styczynski/latte-compiler/src/compiler"
	"github.com/styczynski/latte-compiler/src/events_collector"
	"github.com/styczynski/latte-compiler/src/flow_analysis"
	"github.com/styczynski/latte-compiler/src/input_reader"
	"github.com/styczynski/latte-compiler/src/parser"
	context2 "github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/printer"
	"github.com/styczynski/latte-compiler/src/type_checker"
)

/**
`

class x {
	int a;
	int b;
}

class y {
	int c;
}

//// iteracyjnie
//int fact (int n) {
//  int i,r;
//  int[] q = new int[2];
//  //q = (new int[]);
//  int ff;
//  ff = 2 + 2;
//  //for(int c: q) a = q[a] + 9;
//  return r;
//}
//
//int main (int r) {
//  printInt(fact(7));
//  return 0;
//}

// rekurencyjnie
int factr (int n) {
  if (n < 2)
    return 1 ;
  else
    return (n * factr(n-1)) ;
}

int main() {
	x inst;
    int a;
	string b;
	a = 2;
	b = typename main;
    inst = new x;
}
`
 */

func ActionCompile(c *cli.Context) error {
	pr := printer.CreateLattePrinter()
	eventsCollector := events_collector.StartEventsCollector()
	context := context2.NewParsingContext(pr, eventsCollector)

	tc := type_checker.CreateLatteTypeChecker()
	p := parser.CreateLatteParser()

	inputPaths := c.Args().Slice()
	reader := input_reader.CreateLatteInputReader(inputPaths)
	comp := compiler.CreateLatteCompiler()
	analyzer := flow_analysis.CreateLatteFlowAnalyzer()
	ast := p.ParseInput(reader, context)

	checkedProgram := tc.Check(ast, context)
	analyzedProgram := analyzer.Analyze(checkedProgram, context)
	compiledProgram := comp.Compile(analyzedProgram, context)
	//summary := events_collector.CreateCliSummaryShortStatus()
	summary := events_collector.CreateCliSummary(-1)
	message, ok := eventsCollector.SummarizeCompilation(summary, compiledProgram, context)

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
		},
		Action: ActionCompile,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
