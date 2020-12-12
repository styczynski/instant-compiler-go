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

	fmt.Printf("PROGRAM:\n%s\n", ast.Print(context))

	//tc.Test(context)

	//fmt.Printf("%s", content)
}
