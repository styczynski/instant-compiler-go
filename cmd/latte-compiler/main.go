package main

import (
	"fmt"
	"strings"

	"github.com/styczynski/latte-compiler/cmd/latte-compiler/config"
	"github.com/styczynski/latte-compiler/src/parser"
	"github.com/styczynski/latte-compiler/src/printer"
)

func main() {
	// load application configurations
	if err := config.LoadConfig("./config"); err != nil {
		panic(fmt.Errorf("invalid application configuration: %s", err))
	}

	fmt.Println(config.Config.ConfigVar)

	p := parser.CreateLatteParser()
	ast, err := p.ParseInput(strings.NewReader("int test(int y) {" +
		"2*2;" +
		"}"))
	if err != nil {
		panic(err)
	}

	pr := printer.CreateLattePrinter()
	fmt.Printf("%s", pr.StructRepr(ast))
}
