package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Addition struct {
	generic_ast.BaseASTNode
	Multiplication *Multiplication `@@`
	Op             string          `[ @( "-" | "+" )`
	Next           *Addition       `  @@ ]`
	ParentNode     generic_ast.TraversableNode
}

func (ast *Addition) ExtractConst() (generic_ast.TraversableNode, bool) {
	if ast.HasNext() {
		return nil, false
	}
	return ast.Multiplication.ExtractConst()
}

func (ast *Addition) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *Addition) OverrideParent(node generic_ast.TraversableNode) {
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

func (ast *Addition) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{
		ast.Multiplication,
		generic_ast.MakeTraversableNodeToken(ast, ast.Op, ast.Pos, ast.EndPos),
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

func (ast *Addition) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	next := ast.Next
	if ast.HasNext() {
		next = mapper(ast, ast.Next, context, false).(*Addition)
	}
	return mapper(parent, &Addition{
		BaseASTNode:    ast.BaseASTNode,
		Multiplication: mapper(ast, ast.Multiplication, context, false).(*Multiplication),
		Op:             ast.Op,
		Next:           next,
		ParentNode:     parent.(generic_ast.TraversableNode),
	}, context, true)
}

func (ast *Addition) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(ast, ast.Multiplication, context)
	if ast.HasNext() {
		mapper(ast, ast.Next, context)
	}
	mapper(parent, ast, context)
}

func (ast *Addition) Fn(c hindley_milner.InferContext) generic_ast.Expression {
	return &BuiltinFunction{
		BaseASTNode: ast.BaseASTNode,
		name:        ast.Op,
	}
}

func (ast *Addition) Body() generic_ast.Expression {
	if !ast.HasNext() {
		return ast.Multiplication
	}
	return hindley_milner.Batch{
		Exp: []generic_ast.Expression{
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

//

func (ast *Addition) ConstFold() generic_ast.TraversableNode {
	if ast.HasNext() {
		const1, ok1 := ast.Multiplication.ExtractConst()
		const2, ok2 := ast.Next.Multiplication.ExtractConst()
		if ok1 && ok2 {
			p1 := const1.(*Primary)
			p2 := const2.(*Primary)
			v := p1.Add(p2, ast.Op)
			// Change pointers
			ast.Multiplication.Unary.UnaryApplication.Index.Primary = v
			ast.Op = ast.Next.Op
			ast.Next = ast.Next.Next
			return ast
		}
	}
	return ast
}
