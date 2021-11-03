package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Unary struct {
	 generic_ast.BaseASTNode
	Op      string   `  ( @( "!" | "-" )`
	Unary   *Unary   `    @@ )`
	Primary *Primary `| @@`
	ParentNode generic_ast.TraversableNode
}

func (ast *Unary) ExtractConst() (generic_ast.TraversableNode, bool) {
	if ast.IsPrimary() {
		return ast.Primary.ExtractConst()
	}
	return nil, false
}

func (ast *Unary) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *Unary) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *Unary) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Unary) End() lexer.Position {
	return ast.EndPos
}

func (ast *Unary) GetNode() interface{} {
	return ast
}

func (ast *Unary) GetChildren() []generic_ast.TraversableNode {
	if ast.IsOperation() {
		return []generic_ast.TraversableNode{
			generic_ast.MakeTraversableNodeToken(ast, ast.Op, ast.Pos, ast.EndPos),
			ast.Unary,
		}
	} else if ast.IsPrimary() {
		return []generic_ast.TraversableNode{
			ast.Primary,
		}
	}
	return []generic_ast.TraversableNode{}
}

func (ast *Unary) IsOperation() bool {
	return ast.Unary != nil
}

func (ast *Unary) IsPrimary() bool {
	return ast.Primary != nil
}

func (ast *Unary) Print(c *context.ParsingContext) string {
	if ast.IsOperation() {
		return printUnaryOperation(c, ast, ast.Op, ast.Unary.Print(c))
	} else if ast.IsPrimary() {
		return ast.Primary.Print(c)
	}
	panic("Unvalid Unary value")
}

////


func (ast *Unary) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	if ast.IsOperation() {
		return mapper(parent, &Unary{
			BaseASTNode:      ast.BaseASTNode,
			Op:               ast.Op,
			Unary:            mapper(ast, ast.Unary, context, false).(*Unary),
			ParentNode: parent.(generic_ast.TraversableNode),
		}, context, true)
	} else if ast.IsPrimary() {
		return mapper(parent, &Unary{
			BaseASTNode:      ast.BaseASTNode,
			Primary: mapper(ast, ast.Primary, context, false).(*Primary),
			ParentNode: parent.(generic_ast.TraversableNode),
		}, context, true)
	}
	panic("Invalid Unary operation type")
}

func (ast *Unary) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	if ast.IsOperation() {
		mapper(ast, ast.Unary, context)
	} else if ast.IsPrimary() {
		mapper(ast, ast.Primary, context)
	}
	mapper(parent, ast, context)
}

func (ast *Unary) Fn() generic_ast.Expression {
	return &BuiltinFunction{
		BaseASTNode: ast.BaseASTNode,
		name: "unary_"+ast.Op,
	}
}

func (ast *Unary) Body() generic_ast.Expression {
	if ast.IsPrimary() {
		return ast.Primary
	}
	return hindley_milner.Batch{
		Exp: []generic_ast.Expression{
			ast.Unary,
		},
	}
}

func (ast *Unary) ExpressionType() hindley_milner.ExpressionType {
	if ast.IsPrimary() {
		return hindley_milner.E_PROXY
	}
	return hindley_milner.E_APPLICATION
}