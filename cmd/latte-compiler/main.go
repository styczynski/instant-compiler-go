package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strings"

	//"os"
	//"strings"

	"github.com/styczynski/latte-compiler/cmd/latte-compiler/config"
	"github.com/styczynski/latte-compiler/src/parser"
	context2 "github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/printer"

	//"github.com/styczynski/latte-compiler/src/printer"
	"github.com/styczynski/latte-compiler/src/type_checker"
)

func main() {
	// load application configurations
	if err := config.LoadConfig("./config"); err != nil {
		panic(fmt.Errorf("invalid application configuration: %s", err))
	}

	pr := printer.CreateLattePrinter()
	context := context2.NewParsingContext(pr)
	tc := type_checker.CreateLatteTypeChecker()
	p := parser.CreateLatteParser()
	ast, latteError := p.ParseInput(strings.NewReader(`
// iteracyjnie
int fact (int n) {
  int i,r ;
  i = 2 ;
  r = 1 ;
  while (2) {
    return r;
  }
  return r ;
}
int main (int x) {
  int q = 1;
  printInt(fact(2)) ;
  return 0 ;
}

`), context)
	if latteError != nil {
		fmt.Print(latteError.CliMessage())
		os.Exit(1)
	}

	//content, err := pr.Format(ast, context)
	//if err != nil {
	//	panic(err)
	//}

	fmt.Printf("PARSED\n")
	f, err := os.Create("compiler.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	err = tc.Check(ast, context)
	if err != nil {
		fmt.Print(err.(*type_checker.TypeCheckingError).CliMessage())
		os.Exit(1)
	}

	//tc.Test(context)

	//fmt.Printf("%s", content)
}
