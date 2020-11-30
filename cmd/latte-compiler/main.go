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

	context := parser.NewParsingContext()
	p := parser.CreateLatteParser()
	ast, err := p.ParseInput(strings.NewReader(`
int main () {
  printInt(fact(7)) ;
  printInt(factr(7)) ;
  return 0 ;
}

// iteracyjnie
int fact (int n) {
  int i,r ;
  i = 1 ;
  r = 1 ;
  while (i < n+1) {
    r = r * i ;
    i++ ;
while (i < n+1) {
    r = r * i ;
    i++ ;
  }
  }
  return r ;
}

// rekurencyjnie
int factr (int n) {
  if (n < 2) {
    return 1 ;
 } else {
    return (n * factr(n-1)) ;
 }
}
`), context)
	if err != nil {
		panic(err)
	}

	pr := printer.CreateLattePrinter()
	content, err := pr.Format(ast, context)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", content)
}
