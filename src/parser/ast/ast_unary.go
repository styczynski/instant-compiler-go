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
	UnaryApplication *UnaryApplication `| @@`
	ParentNode generic_ast.TraversableNode
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
	return "UNKNOWN"
}

////


func (ast *Unary) Map(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	if ast.IsOperation() {
		return mapper(parent, &Unary{
			BaseASTNode:      ast.BaseASTNode,
			Op:               ast.Op,
			Unary:            mapper(ast, ast.Unary).(*Unary),
			ParentNode: parent.(generic_ast.TraversableNode),
		})
	} else if ast.IsUnaryApplication() {
		return mapper(parent, &Unary{
			BaseASTNode:      ast.BaseASTNode,
			UnaryApplication: mapper(ast, ast.UnaryApplication).(*UnaryApplication),
			ParentNode: parent.(generic_ast.TraversableNode),
		})
	}
	panic("Invalid Unary operation type")
}

func (ast *Unary) Visit(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) {
	if ast.IsOperation() {
		mapper(ast, ast.Unary)
	} else if ast.IsUnaryApplication() {
		mapper(ast, ast.UnaryApplication)
	}
	mapper(parent, ast)
}

func (ast *Unary) Fn() hindley_milner.Expression {
	return &BuiltinFunction{
		BaseASTNode: ast.BaseASTNode,
		name: "unary_"+ast.Op,
	}
}

func (ast *Unary) Body() hindley_milner.Expression {
	if ast.IsUnaryApplication() {
		return ast.UnaryApplication
	}
	return hindley_milner.Batch{
		Exp: []hindley_milner.Expression{
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