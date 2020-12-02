package type_checker

import "github.com/styczynski/latte-compiler/src/parser"

type LatteTypeChecker struct {}

func CreateLatteTypeChecker() *LatteTypeChecker {
	return &LatteTypeChecker{}
}

func (tc *LatteTypeChecker) Test(c *parser.ParsingContext) {
	Example_greenspun()
}