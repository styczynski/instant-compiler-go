package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Unary struct {
	generic_ast.BaseASTNode
	Op               string            `  ( @( "!" | "-" )`
	Unary            *Unary            `    @@ )`
	UnaryApplication *UnaryApplication `| @@`
	ParentNode       generic_ast.TraversableNode
	ResolvedType     hindley_milner.Type
}

func (ast *Unary) OnTypeReturned(t hindley_milner.Type) {
	ast.ResolvedType = t
}

func (ast *Unary) ExtractConst() (generic_ast.TraversableNode, bool) {
	if ast.IsUnaryApplication() {
		return ast.UnaryApplication.ExtractConst()
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
	} else if ast.IsUnaryApplication() {
		return []generic_ast.TraversableNode{
			ast.UnaryApplication,
		}
	}
	return []generic_ast.TraversableNode{}
}

func (ast *Unary) IsOperation() bool {
	return ast.Unary != nil
}

func (ast *Unary) IsUnaryApplication() bool {
	return ast.UnaryApplication != nil
}

func (ast *Unary) Print(c *context.ParsingContext) string {
	if ast.IsOperation() {
		return printUnaryOperation(c, ast, ast.Op, ast.Unary.Print(c))
	} else if ast.IsUnaryApplication() {
		return ast.UnaryApplication.Print(c)
	}
	panic("Unvalid Unary value")
}

////

func (ast *Unary) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	if ast.IsOperation() {
		return mapper(parent, &Unary{
			BaseASTNode: ast.BaseASTNode,
			Op:          ast.Op,
			Unary:       mapper(ast, ast.Unary, context, false).(*Unary),
			ParentNode:  parent.(generic_ast.TraversableNode),
		}, context, true)
	} else if ast.IsUnaryApplication() {
		return mapper(parent, &Unary{
			BaseASTNode:      ast.BaseASTNode,
			UnaryApplication: mapper(ast, ast.UnaryApplication, context, false).(*UnaryApplication),
			ParentNode:       parent.(generic_ast.TraversableNode),
		}, context, true)
	}
	panic("Invalid Unary operation type")
}

func (ast *Unary) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	if ast.IsOperation() {
		mapper(ast, ast.Unary, context)
	} else if ast.IsUnaryApplication() {
		mapper(ast, ast.UnaryApplication, context)
	}
	mapper(parent, ast, context)
}

func (ast *Unary) Fn(c hindley_milner.InferContext) generic_ast.Expression {
	return &BuiltinFunction{
		BaseASTNode: ast.BaseASTNode,
		name:        "unary_" + ast.Op,
	}
}

func (ast *Unary) Body() generic_ast.Expression {
	if ast.IsUnaryApplication() {
		return ast.UnaryApplication
	}
	return hindley_milner.Batch{
		Exp: []generic_ast.Expression{
			ast.Unary,
		},
	}
}

func (ast *Unary) ExpressionType() hindley_milner.ExpressionType {
	if ast.IsUnaryApplication() {
		return hindley_milner.E_PROXY
	}
	return hindley_milner.E_APPLICATION
}
