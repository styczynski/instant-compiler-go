package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type LogicalOperation struct {
	BaseASTNode
	Equality *Equality `@@`
	Op         string      `[ @( "|" "|" | "&" "&" )`
	Next       *LogicalOperation   `  @@ ]`
}

func (ast *LogicalOperation) Begin() lexer.Position {
	return ast.Pos
}

func (ast *LogicalOperation) End() lexer.Position {
	return ast.EndPos
}

func (ast *LogicalOperation) GetNode() interface{} {
	return ast
}

func (ast *LogicalOperation) GetChildren() []TraversableNode {
	return []TraversableNode{
		ast.Equality,
		MakeTraversableNodeToken(ast.Op, ast.Pos, ast.EndPos),
		ast.Next,
	}
}

func (ast *LogicalOperation) HasNext() bool {
	return ast.Next != nil
}

func (ast *LogicalOperation) Print(c *context.ParsingContext) string {
	if ast.HasNext() {
		return printBinaryOperation(c, ast, ast.Equality.Print(c), ast.Op, ast.Next.Print(c))
	}
	return ast.Equality.Print(c)
}

/////

func (ast *LogicalOperation) Map(mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	next := ast.Next
	if ast.HasNext() {
		next = mapper(ast.Next).(*LogicalOperation)
	}
	return mapper(&LogicalOperation{
		BaseASTNode: ast.BaseASTNode,
		Equality:    mapper(ast.Equality).(*Equality),
		Op:          ast.Op,
		Next:        next,
	})
}

func (ast *LogicalOperation) Visit(mapper hindley_milner.ExpressionMapper) {
	mapper(ast.Equality)
	if ast.HasNext() {
		mapper(ast.Next)
	}
	mapper(ast)
}

func (ast *LogicalOperation) Fn() hindley_milner.Expression {
	return &BuiltinFunction{
		BaseASTNode: ast.BaseASTNode,
		name: ast.Op,
	}
}

func (ast *LogicalOperation) Body() hindley_milner.Expression {
	if !ast.HasNext() {
		return ast.Equality
	}
	return hindley_milner.Batch{
		Exp: []hindley_milner.Expression{
			ast.Equality,
			ast.Next,
		},
	}
}

func (ast *LogicalOperation) ExpressionType() hindley_milner.ExpressionType {
	if !ast.HasNext() {
		return hindley_milner.E_PROXY
	}
	return hindley_milner.E_APPLICATION
}