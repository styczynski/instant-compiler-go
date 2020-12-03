package ast

import (
	"strings"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
)

type Block struct {
	BaseASTNode
	Statements []*Statement `"{" @@* "}"`
}

func (ast *Block) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Block) End() lexer.Position {
	return ast.EndPos
}

func (ast *Block) GetNode() interface{} {
	return ast
}

func (ast *Block) Print(c *context.ParsingContext) string {
	statementsList := []string{}
	if c.PrinterConfiguration.SkipStatementIdent {
		c.PrinterConfiguration.SkipStatementIdent = false
	}
	c.BlockPush()
	for _, statement := range ast.Statements {
		statementsList = append(statementsList, statement.Print(c))
	}
	c.BlockPop()
	return printNode(c, ast, "{\n%s\n%s}", strings.Join(statementsList, "\n"), strings.Repeat("  ", c.BlockDepth))
}

func (ast *Block) GetChildren() []TraversableNode {
	nodes := make([]TraversableNode, len(ast.Statements))
	for _, child := range ast.Statements {
		nodes = append(nodes, child)
	}
	return nodes
}