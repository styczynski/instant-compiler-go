package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Addition struct {
	BaseASTNode
	Multiplication *Multiplication `@@`
	Op             string          `[ @( "-" | "+" )`
	Next           *Addition       `  @@ ]`
	ParentNode TraversableNode
}

func (ast *Addition) Parent() TraversableNode {
	return ast.ParentNode
}

func (ast *Addition) OverrideParent(node TraversableNode) {
	ast.ParentNode = node
}

func (ast *Addition) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Addition) End() lexer.Position {
	return ast.EndPos
}

func (ast *Addition) GetNode() interface{} {
	return ast
}

func (ast *Addition) GetChildren() []TraversableNode {
	return []TraversableNode{
		ast.Multiplication,
		MakeTraversableNodeToken(ast, ast.Op, ast.Pos, ast.EndPos),
		ast.Next,
	}
}

func (ast *Addition) HasNext() bool {
	return ast.Next != nil
}

func (ast *Addition) Print(c *context.ParsingContext) string {
	if ast.HasNext() {
		return printBinaryOperation(c, ast, ast.Multiplication.Print(c), ast.Op, ast.Next.Print(c))
	}
	return ast.Multiplication.Print(c)
}

///


func (ast *Addition) Map(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	next := ast.Next
	if ast.HasNext() {
		next = mapper(ast, ast.Next).(*Addition)
	}
	return mapper(parent, &Addition{
		BaseASTNode: ast.BaseASTNode,
		Multiplication:    mapper(ast, ast.Multiplication).(*Multiplication),
		Op:          ast.Op,
		Next:        next,
		ParentNode: parent.(TraversableNode),
	})
}

func (ast *Addition) Visit(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) {
	mapper(ast, ast.Multiplication)
	if ast.HasNext() {
		mapper(ast, ast.Next)
	}
	mapper(parent, ast)
}

func (ast *Addition) Fn() hindley_milner.Expression {
	return &BuiltinFunction{
		BaseASTNode: ast.BaseASTNode,
		name: ast.Op,
	}
}

func (ast *Addition) Body() hindley_milner.Expression {
	if !ast.HasNext() {
		return ast.Multiplication
	}
	return hindley_milner.Batch{
		Exp: []hindley_milner.Expression{
			ast.Multiplication,
			ast.Next,
		},
	}
}

func (ast *Addition) ExpressionType() hindley_milner.ExpressionType {
	if !ast.HasNext() {
		return hindley_milner.E_PROXY
	}
	return hindley_milner.E_APPLICATION
}
