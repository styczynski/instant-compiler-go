package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Multiplication struct {
	generic_ast.BaseASTNode
	Primary    *Primary        `@@`
	Op         string          `[ @( "/" | "*" | "%" )`
	Next       *Multiplication `  @@ ]`
	ParentNode generic_ast.TraversableNode
}

func (ast *Multiplication) ExtractConst() (generic_ast.TraversableNode, bool) {
	if ast.HasNext() {
		return nil, false
	}
	return ast.Primary.ExtractConst()
}

func (ast *Multiplication) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *Multiplication) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
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

func (ast *Multiplication) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{
		ast.Primary,
		generic_ast.MakeTraversableNodeToken(ast, ast.Op, ast.Pos, ast.EndPos),
		ast.Next,
	}
}

func (ast *Multiplication) HasNext() bool {
	return ast.Next != nil
}

func (ast *Multiplication) Print(c *context.ParsingContext) string {
	if ast.HasNext() {
		return printBinaryOperation(c, ast, ast.Primary.Print(c), ast.Op, ast.Next.Print(c))
	}
	return ast.Primary.Print(c)
}

////

func (ast *Multiplication) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	next := ast.Next
	if ast.HasNext() {
		next = mapper(ast, ast.Next, context, false).(*Multiplication)
	}
	return mapper(parent, &Multiplication{
		BaseASTNode: ast.BaseASTNode,
		Primary:     mapper(ast, ast.Primary, context, false).(*Primary),
		Op:          ast.Op,
		Next:        next,
		ParentNode:  parent.(generic_ast.TraversableNode),
	}, context, true)
}

func (ast *Multiplication) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(ast, ast.Primary, context)
	if ast.HasNext() {
		mapper(ast, ast.Next, context)
	}
	mapper(parent, ast, context)
}

func (ast *Multiplication) Fn() generic_ast.Expression {
	return &BuiltinFunction{
		BaseASTNode: ast.BaseASTNode,
		name:        ast.Op,
	}
}

func (ast *Multiplication) Body() generic_ast.Expression {
	if !ast.HasNext() {
		return ast.Primary
	}
	return hindley_milner.Batch{
		Exp: []generic_ast.Expression{
			ast.Primary,
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
