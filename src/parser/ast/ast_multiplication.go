package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Multiplication struct {
	BaseASTNode
	Unary *Unary          `@@`
	Op    string          `[ @( "/" | "*" )`
	Next  *Multiplication `  @@ ]`
}

func (ast *Multiplication) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Multiplication) End() lexer.Position {
	return ast.EndPos
}

func (ast *Multiplication) GetNode() interface{} {
	return ast
}

func (ast *Multiplication) GetChildren() []TraversableNode {
	return []TraversableNode{
		ast.Unary,
		MakeTraversableNodeToken(ast.Op, ast.Pos, ast.EndPos),
		ast.Next,
	}
}

func (ast *Multiplication) HasNext() bool {
	return ast.Next != nil
}

func (ast *Multiplication) Print(c *context.ParsingContext) string {
	if ast.HasNext() {
		return printBinaryOperation(c, ast, ast.Unary.Print(c), ast.Op, ast.Next.Print(c))
	}
	return ast.Unary.Print(c)
}


////


func (ast *Multiplication) Map(mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	next := ast.Next
	if ast.HasNext() {
		next = mapper(ast.Next).(*Multiplication)
	}
	return mapper(&Multiplication{
		BaseASTNode: ast.BaseASTNode,
		Unary:    mapper(ast.Unary).(*Unary),
		Op:          ast.Op,
		Next:        next,
	})
}

func (ast *Multiplication) Visit(mapper hindley_milner.ExpressionMapper) {
	mapper(ast.Unary)
	if ast.HasNext() {
		mapper(ast.Next)
	}
	mapper(ast)
}

func (ast *Multiplication) Fn() hindley_milner.Expression {
	return &BuiltinFunction{
		name: ast.Op,
	}
}

func (ast *Multiplication) Body() hindley_milner.Expression {
	if !ast.HasNext() {
		return ast.Unary
	}
	return hindley_milner.Batch{
		Exp: []hindley_milner.Expression{
			ast.Unary,
			ast.Next,
		},
	}
}

func (ast *Multiplication) ExpressionType() hindley_milner.ExpressionType {
	if !ast.HasNext() {
		return hindley_milner.E_PROXY
	}
	return hindley_milner.E_APPLICATION
}