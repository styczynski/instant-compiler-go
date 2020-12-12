package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
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

///

func (ast *While) Map(mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	return mapper(&While{
		BaseASTNode: ast.BaseASTNode,
		Condition: mapper(ast.Condition).(*Expression),
		Do: mapper(ast.Do).(*Statement),
	})
}

func (ast *While) Visit(mapper hindley_milner.ExpressionMapper) {
	mapper(ast.Condition)
	mapper(ast.Do)
	mapper(ast)
}

func (ast *While) Fn() hindley_milner.Expression {
	return &BuiltinFunction{
		BaseASTNode: ast.BaseASTNode,
		name: "while",
	}
}

func (ast *While) Body() hindley_milner.Expression {
	args := []hindley_milner.Expression{
		ast.Condition,
		ast.Do,
	}
	return hindley_milner.Batch{
		Exp: args,
	}
}

func (ast *While) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_APPLICATION
}