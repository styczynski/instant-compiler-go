package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Expression struct {
	ComplexASTNode
	NewType       *Type       `( "new" @@ )`
	LogicalOperation *LogicalOperation `| @@`
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

func (ast *Expression) IsLogicalOperation() bool {
	return ast.LogicalOperation != nil
}

func (ast *Expression) IsNewType() bool {
	return ast.NewType != nil
}

func (ast *Expression) GetChildren() []TraversableNode {
	if ast.IsLogicalOperation() {
		return []TraversableNode{ast.LogicalOperation,}
	} else if ast.IsNewType() {
		return []TraversableNode{ast.NewType,}
	}
	panic("Invalid Expression type")
}

func printBinaryOperation(c *context.ParsingContext, ast TraversableNode, arg1 string, operator string, arg2 string) string{
	return printNode(c, ast, "%s %s %s", arg1, operator, arg2)
}

func printUnaryOperation(c *context.ParsingContext, ast TraversableNode, operator string, arg string) string{
	return printNode(c, ast, "%s%s", operator, arg)
}

func (ast *Expression) Print(c *context.ParsingContext) string {
	if ast.IsLogicalOperation() {
		return ast.LogicalOperation.Print(c)
	} else if ast.IsNewType() {
		return printNode(c, ast, "new %s", ast.NewType.Print(c))
	}
	panic("Invalid Expression type")
}

////

func (ast *Expression) Body() hindley_milner.Expression {
	return ast.LogicalOperation
}

func (ast *Expression) Map(mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	if ast.IsLogicalOperation() {
		return mapper(&Expression{
			ComplexASTNode:   ast.ComplexASTNode,
			LogicalOperation: mapper(ast.LogicalOperation).(*LogicalOperation),
		})
	} else if ast.IsNewType() {
		return mapper(&Expression{
			ComplexASTNode:   ast.ComplexASTNode,
			NewType: ast.NewType,
		})
	}
	panic("Invalid Expression type")
}

func (ast *Expression) Visit(mapper hindley_milner.ExpressionMapper) {
	if ast.IsLogicalOperation() {
		mapper(ast.LogicalOperation)
	}
	mapper(ast)
}

func (ast *Expression) EmbeddedType() *hindley_milner.Scheme {
	return ast.NewType.GetType()
}

func (ast *Expression) ExpressionType() hindley_milner.ExpressionType {
	if ast.IsNewType() {
		return hindley_milner.E_TYPE
	}
	return hindley_milner.E_PROXY
}