package main

import (
	"fmt"
	"os"
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

	context := context2.NewParsingContext()
	tc := type_checker.CreateLatteTypeChecker()
	pr := printer.CreateLattePrinter()
	p := parser.CreateLatteParser(pr)
	ast, latteError := p.ParseInput(strings.NewReader(`
int main (int a) {
 bool e = !(2<3);
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

	tc.Check(ast, context)

	//tc.Test(context)

	//fmt.Printf("%s", content)
}
