package ast

import (
	"fmt"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
)

var SUGGESTED_KEYWORDS = []string{
	"int",
	"string",
	"bool",
	"true",
	"false",
	"void",
	"return",
	"if",
	"while",
	"else",
}

func printNode(c *context.ParsingContext, ast generic_ast.TraversableNode,format string, args ...interface{}) string {
	if c.PrinterConfiguration.MaxPrintPosition != nil {
		if ast.Begin().Line > c.PrinterConfiguration.MaxPrintPosition.Line {
			return ""
		}
		if ast.Begin().Line == c.PrinterConfiguration.MaxPrintPosition.Line && ast.Begin().Column > c.PrinterConfiguration.MaxPrintPosition.Column {
			return ""
		}
	}
	return fmt.Sprintf(format, args...)
}