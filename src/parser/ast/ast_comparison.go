package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Comparison struct {
	BaseASTNode
	Addition *Addition   `@@`
	Op       string      `[ @( ">" "=" | "<" "=" | "=" "=" | ">" | "<" )`
	Next     *Comparison `  @@ ]`
	ParentNode TraversableNode
}

func (ast *Comparison) Parent() TraversableNode {
	return ast.ParentNode
}

func (ast *Comparison) OverrideParent(node TraversableNode) {
	ast.ParentNode = node
}

func (ast *Comparison) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Comparison) End() lexer.Position {
	return ast.EndPos
}

func (ast *Comparison) GetNode() interface{} {
	return ast
}

func (ast *Comparison) GetChildren() []TraversableNode {
	return []TraversableNode{
		ast.Addition,
		MakeTraversableNodeToken(ast, ast.Op, ast.Pos, ast.EndPos),
		ast.Next,
	}
}

func (ast *Comparison) HasNext() bool {
	return ast.Next != nil
}

func (ast *Comparison) Print(c *context.ParsingContext) string {
	if ast.HasNext() {
		return printBinaryOperation(c, ast, ast.Addition.Print(c), ast.Op, ast.Next.Print(c))
	}
	return ast.Addition.Print(c)
}


////

func (ast *Comparison) Map(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	next := ast.Next
	if ast.HasNext() {
		next = mapper(ast, ast.Next).(*Comparison)
	}
	return mapper(parent, &Comparison{
		BaseASTNode: ast.BaseASTNode,
		Addition:    mapper(ast, ast.Addition).(*Addition),
		Op:          ast.Op,
		Next:        next,
		ParentNode: parent.(TraversableNode),
	})
}

func (ast *Comparison) Visit(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) {
	mapper(ast, ast.Addition)
	if ast.HasNext() {
		mapper(ast, ast.Next)
	}
	mapper(parent, ast)
}

func (ast *Comparison) Fn() hindley_milner.Expression {
	return &BuiltinFunction{
		BaseASTNode: ast.BaseASTNode,
		name: ast.Op,
	}
}

func (ast *Comparison) Body() hindley_milner.Expression {
	if !ast.HasNext() {
		return ast.Addition
	}
	return hindley_milner.Batch{
		Exp: []hindley_milner.Expression{
			ast.Addition,
			ast.Next,
		},
	}
}

func (ast *Comparison) ExpressionType() hindley_milner.ExpressionType {
	if !ast.HasNext() {
		return hindley_milner.E_PROXY
	}
	return hindley_milner.E_APPLICATION
}