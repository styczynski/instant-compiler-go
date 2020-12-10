package ast

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
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

///////

func (ast *Block) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_BLOCK
}

func (ast *Block) Map(mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	mappedStmts := []*Statement{}
	for _, stmt := range ast.Statements {
		mappedStmts = append(mappedStmts, mapper(stmt).(*Statement))
	}
	return mapper(&Block{
		BaseASTNode: ast.BaseASTNode,
		Statements:  mappedStmts,
	})
}

func (ast *Block) Visit(mapper hindley_milner.ExpressionMapper) {
	for _, stmt := range ast.Statements {
		mapper(stmt)
	}
	mapper(ast)
}

func (ast *Block) GetContents() hindley_milner.Batch {
	exprs := []hindley_milner.Expression{}
	for _, stmt := range ast.Statements {
		exprs = append(exprs, stmt)
	}
	return hindley_milner.Batch{
		Exp: exprs,
	}
}

func (ast *Block) Expressions() []hindley_milner.Expression {
	return ast.GetContents().Exp
}

func (ast *Block) Body() hindley_milner.Expression {
	panic(fmt.Errorf("Batch Body() method cannot be called."))
}
