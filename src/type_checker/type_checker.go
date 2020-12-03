package type_checker

import "github.com/styczynski/latte-compiler/src/parser/context"

type LatteTypeChecker struct {}

func CreateLatteTypeChecker() *LatteTypeChecker {
	return &LatteTypeChecker{}
}

func (tc *LatteTypeChecker) Test(c *context.ParsingContext) {
	Example_greenspun()
}