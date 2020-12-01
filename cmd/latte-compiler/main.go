package main

import (
	"fmt"
	"os"
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

	context := parser.NewParsingContext()
	pr := printer.CreateLattePrinter()
	p := parser.CreateLatteParser(pr)
	ast, latteError := p.ParseInput(strings.NewReader(`
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
  r;
}

int factr (int n) {
  if (n < 2) {
    return 1 ; // rekurencyjnie
 } else {
    return (n * factr(n-1)) ;
 }
}
`), context)
	if latteError != nil {
		fmt.Print(latteError.CliMessage())
		os.Exit(1)
	}

	content, err := pr.Format(ast, context)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", content)
}
