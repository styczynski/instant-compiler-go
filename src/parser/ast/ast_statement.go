package ast

import (
	"strings"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Statement struct {
	generic_ast.BaseASTNode
	Empty          *string         `";"`
	Assignment     *Assignment     `| @@`
	Expression     *Expression     `| @@ ";"`
	ParentNode     generic_ast.TraversableNode
}

func (ast *Statement) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *Statement) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *Statement) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Statement) End() lexer.Position {
	return ast.EndPos
}

func (ast *Statement) GetNode() interface{} {
	return ast
}

func (ast *Statement) IsEmpty() bool {
	return (ast.Empty != nil || (
		!ast.IsAssignment() &&
		!ast.IsExpression()))
}

func (ast *Statement) IsAssignment() bool {
	return ast.Assignment != nil
}

func (ast *Statement) IsExpression() bool {
	return ast.Expression != nil
}

func (ast *Statement) GetChildren() []generic_ast.TraversableNode {
	if ast.IsEmpty() {
		return []generic_ast.TraversableNode{}
	} else if ast.IsAssignment() {
		return []generic_ast.TraversableNode{ast.Assignment}
	} else if ast.IsExpression() {
		return []generic_ast.TraversableNode{ast.Expression}
	}
	return []generic_ast.TraversableNode{}
}

func (ast *Statement) formatStatementInstruction(statement string, c *context.ParsingContext) string {
	if c.PrinterConfiguration.SkipStatementIdent {
		c.PrinterConfiguration.SkipStatementIdent = false
		return statement
	}
	return printNode(c, ast, "%s%s", strings.Repeat("  ", c.BlockDepth), statement)
}

func (ast *Statement) Print(c *context.ParsingContext) string {
	ret := "UNKNOWN"
	propagateSkipStatementIdent := false
	if ast.IsEmpty() {
		ret = ";"
	} else if ast.IsAssignment() {
		ret = ast.Assignment.Print(c)
	} else if ast.IsExpression() {
		ret = printNode(c, ast, "%s;", ast.Expression.Print(c))
	}
	if propagateSkipStatementIdent {
		c.PrinterConfiguration.SkipStatementIdent = true
		ret = ast.formatStatementInstruction(ret, c)
	} else {
		ret = ast.formatStatementInstruction(ret, c)
	}
	return ret
}

//////

func (ast *Statement) Body() generic_ast.Expression {
	if ast.IsEmpty() {
		return &Empty{}
		//return hindley_milner.Batch{Exp: []generic_ast.Expression{}}
	} else if ast.IsAssignment() {
		return ast.Assignment
	} else if ast.IsExpression() {
		return ast.Expression
	}
	ast.BaseASTNode = generic_ast.BaseASTNode{}
	ast.ParentNode = nil
	//fmt.Printf("Failed for node %s\n", repr.String(ast))
	panic("Invalid Statement type")
}

func feedExpressionIntoStatement(node interface{}, base generic_ast.BaseASTNode) *Statement {
	if assignment, ok := node.(*Assignment); ok {
		return &Statement{
			BaseASTNode: base,
			Assignment:  assignment,
		}
	} else if expr, ok := node.(*Expression); ok {
		return &Statement{
			BaseASTNode: base,
			Expression:  expr,
		}
	}
	panic("Invalid statement type")
}

func (ast *Statement) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	mappedStatement := feedExpressionIntoStatement(mapper(ast, ast.Body(), context, false), ast.BaseASTNode)
	return mapper(parent, mappedStatement, context, true)
}

func (ast *Statement) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(ast, ast.Body(), context)
	mapper(parent, ast, context)
}

func (ast *Statement) ExpressionType() hindley_milner.ExpressionType { return hindley_milner.E_PROXY }
