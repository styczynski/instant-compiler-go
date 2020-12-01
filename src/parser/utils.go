package parser

import "fmt"

func TraverseAST(node TraversableNode, visitor func(ast TraversableNode)) {
	children := node.GetChildren()
	for _, child := range children {
		visitor(child)
		TraverseAST(child, visitor)
	}
}

func printNode(c *ParsingContext, ast TraversableNode,format string, args ...interface{}) string {
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