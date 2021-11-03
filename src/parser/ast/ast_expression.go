package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Expression struct {
	generic_ast.ComplexASTNode
	Addition   Addition `@@`
	ParentNode generic_ast.TraversableNode
}

func (ast *Expression) ExtractConst() (generic_ast.TraversableNode, bool) {
	return ast.Addition.ExtractConst()
}

func (ast *Expression) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *Expression) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *Expression) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Expression) End() lexer.Position {
	return ast.EndPos
}

func (ast *Expression) GetNode() interface{} {
	return ast
}

func (ast *Expression) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{&ast.Addition}
}

func printBinaryOperation(c *context.ParsingContext, ast generic_ast.TraversableNode, arg1 string, operator string, arg2 string) string {
	return printNode(c, ast, "%s %s %s", arg1, operator, arg2)
}

func printUnaryOperation(c *context.ParsingContext, ast generic_ast.TraversableNode, operator string, arg string) string {
	return printNode(c, ast, "%s%s", operator, arg)
}

func (ast *Expression) Print(c *context.ParsingContext) string {
	return ast.Addition.Print(c)
}

////

func (ast *Expression) Body() generic_ast.Expression {
	return &ast.Addition
}

func (ast *Expression) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	return mapper(parent, &Expression{
		ComplexASTNode: ast.ComplexASTNode,
		Addition:       ast.Addition,
		ParentNode:     parent.(generic_ast.TraversableNode),
	}, context, true)
}

func (ast *Expression) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(ast, &ast.Addition, context)
	mapper(parent, ast, context)
}

func (ast *Expression) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_PROXY
}
