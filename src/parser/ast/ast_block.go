package ast

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Block struct {
	generic_ast.BaseASTNode
	Statements []*Statement `"{" @@* "}"`
	ParentNode generic_ast.TraversableNode
}

func (ast *Block) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *Block) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
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

func (ast *Block) GetChildren() []generic_ast.TraversableNode {
	nodes := []generic_ast.TraversableNode{}
	for _, child := range ast.Statements {
		nodes = append(nodes, child)
	}
	return nodes
}

///////

func (ast *Block) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_BLOCK
}

func (ast *Block) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	mappedStmts := []*Statement{}
	for _, stmt := range ast.Statements {
		mappedStmts = append(mappedStmts, mapper(ast, stmt, context, false).(*Statement))
	}
	return mapper(parent, &Block{
		BaseASTNode: ast.BaseASTNode,
		Statements:  mappedStmts,
		ParentNode:  parent.(generic_ast.TraversableNode),
	}, context, true)
}

func (ast *Block) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	for _, stmt := range ast.Statements {
		mapper(ast, stmt, context)
	}
	mapper(parent, ast, context)
}

func (ast *Block) GetContents() hindley_milner.Batch {
	exprs := []generic_ast.Expression{}
	for _, stmt := range ast.Statements {
		exprs = append(exprs, stmt)
	}
	return hindley_milner.Batch{
		Exp: exprs,
	}
}

func (ast *Block) Expressions() []generic_ast.Expression {
	return ast.GetContents().Exp
}

func (ast *Block) Body() generic_ast.Expression {
	panic(fmt.Errorf("Batch Body() method cannot be called."))
}
