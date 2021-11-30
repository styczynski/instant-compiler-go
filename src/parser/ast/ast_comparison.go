package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Comparison struct {
	generic_ast.BaseASTNode
	Addition   *Addition   `@@`
	Op         string      `[ @( ">" "=" | "<" "=" | "=" "=" | ">" | "<" )`
	Next       *Comparison `  @@ ]`
	ParentNode generic_ast.TraversableNode
}

func (ast *Comparison) ExtractConst() (generic_ast.TraversableNode, bool) {
	if ast.HasNext() {
		return nil, false
	}
	return ast.Addition.ExtractConst()
}

func (ast *Comparison) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *Comparison) OverrideParent(node generic_ast.TraversableNode) {
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

func (ast *Comparison) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{
		ast.Addition,
		generic_ast.MakeTraversableNodeToken(ast, ast.Op, ast.Pos, ast.EndPos),
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

func (ast *Comparison) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	next := ast.Next
	if ast.HasNext() {
		next = mapper(ast, ast.Next, context, false).(*Comparison)
	}
	return mapper(parent, &Comparison{
		BaseASTNode: ast.BaseASTNode,
		Addition:    mapper(ast, ast.Addition, context, false).(*Addition),
		Op:          ast.Op,
		Next:        next,
		ParentNode:  parent.(generic_ast.TraversableNode),
	}, context, true)
}

func (ast *Comparison) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(ast, ast.Addition, context)
	if ast.HasNext() {
		mapper(ast, ast.Next, context)
	}
	mapper(parent, ast, context)
}

func (ast *Comparison) Fn(c hindley_milner.InferContext) generic_ast.Expression {
	return &BuiltinFunction{
		BaseASTNode: ast.BaseASTNode,
		name:        ast.Op,
	}
}

func (ast *Comparison) Body() generic_ast.Expression {
	if !ast.HasNext() {
		return ast.Addition
	}
	return hindley_milner.Batch{
		Exp: []generic_ast.Expression{
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

//

func (ast *Comparison) ConstFold() generic_ast.TraversableNode {
	if ast.HasNext() {
		const1, ok1 := ast.Addition.ExtractConst()
		const2, ok2 := ast.Next.Addition.ExtractConst()
		if ok1 && ok2 {
			p1 := const1.(*Primary)
			p2 := const2.(*Primary)
			v := p1.Compare(p2, ast.Op)
			if v == nil {
				return ast
			}
			// Change pointers
			ast.Addition.Multiplication.Unary.UnaryApplication.Index.Primary = v
			ast.Op = ast.Next.Op
			ast.Next = ast.Next.Next
			return ast
		}
	}
	return ast
}
