package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Expression struct {
	generic_ast.ComplexASTNode
	NewType       *New      `@@`
	Typename *Typename `| @@`
	LogicalOperation *LogicalOperation `| @@`
	ParentNode generic_ast.TraversableNode
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

func (ast *Expression) IsLogicalOperation() bool {
	return ast.LogicalOperation != nil
}

func (ast *Expression) IsNewType() bool {
	return ast.NewType != nil
}

func (ast *Expression) IsTypename() bool {
	return ast.Typename != nil
}

func (ast *Expression) GetChildren() []generic_ast.TraversableNode {
	if ast.IsLogicalOperation() {
		return []generic_ast.TraversableNode{ast.LogicalOperation,}
	} else if ast.IsNewType() {
		return []generic_ast.TraversableNode{ ast.NewType.GetTraversableNode(), }
	} else if ast.IsTypename() {
		return []generic_ast.TraversableNode{ ast.Typename.GetTraversableNode(), }
	}
	panic("Invalid Expression type")
}

func printBinaryOperation(c *context.ParsingContext, ast generic_ast.TraversableNode, arg1 string, operator string, arg2 string) string{
	return printNode(c, ast, "%s %s %s", arg1, operator, arg2)
}

func printUnaryOperation(c *context.ParsingContext, ast generic_ast.TraversableNode, operator string, arg string) string{
	return printNode(c, ast, "%s%s", operator, arg)
}

func (ast *Expression) Print(c *context.ParsingContext) string {
	if ast.IsLogicalOperation() {
		return ast.LogicalOperation.Print(c)
	} else if ast.IsNewType() {
		return ast.NewType.Print(c)
	} else if ast.IsTypename() {
		return ast.Typename.Print(c)
	}
	panic("Invalid Expression type")
}

////

func (ast *Expression) Body() hindley_milner.Expression {
	if ast.IsLogicalOperation() {
		return ast.LogicalOperation
	} else if ast.IsNewType() {
		return ast.NewType
	} else if ast.IsTypename() {
		return ast.Typename
	}
	panic("Invalid Expression type")
}

func (ast *Expression) Map(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	if ast.IsLogicalOperation() {
		return mapper(parent, &Expression{
			ComplexASTNode:   ast.ComplexASTNode,
			LogicalOperation: mapper(ast, ast.LogicalOperation).(*LogicalOperation),
			ParentNode: parent.(generic_ast.TraversableNode),
		})
	} else if ast.IsNewType() {
		return mapper(parent, &Expression{
			ComplexASTNode:   ast.ComplexASTNode,
			NewType: ast.NewType,
			ParentNode: parent.(generic_ast.TraversableNode),
		})
	} else if ast.IsTypename() {
		return mapper(parent, &Expression{
			ComplexASTNode:   ast.ComplexASTNode,
			Typename: ast.Typename,
			ParentNode: parent.(generic_ast.TraversableNode),
		})
	}
	panic("Invalid Expression type")
}

func (ast *Expression) Visit(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) {
	if ast.IsLogicalOperation() {
		mapper(ast, ast.LogicalOperation)
	} else if ast.IsNewType() {
		mapper(ast, ast.NewType)
	} else if ast.IsTypename() {
		mapper(ast, ast.Typename)
	}
	mapper(parent, ast)
}

func (ast *Expression) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_PROXY
}