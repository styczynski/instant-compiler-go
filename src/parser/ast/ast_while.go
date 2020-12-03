package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
)

type While struct {
	BaseASTNode
	Condition *Expression `"while" "(" @@ ")"`
	Do *Statement `@@`
}

func (ast *While) Begin() lexer.Position {
	return ast.Pos
}

func (ast *While) End() lexer.Position {
	return ast.EndPos
}

func (ast *While) GetNode() interface{} {
	return ast
}

func (ast *While) Print(c *context.ParsingContext) string {
	c.PrinterConfiguration.SkipStatementIdent = true
	body := ast.Do.Print(c)
	return printNode(c, ast, "while (%s) %s", ast.Condition.Print(c), body)
}

func (ast *While) GetChildren() []TraversableNode {
	return []TraversableNode{
		ast.Condition,
		ast.Do,
	}
}
