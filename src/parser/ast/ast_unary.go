package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Unary struct {
	BaseASTNode
	Op      string   `  ( @( "!" | "-" )`
	Unary   *Unary   `    @@ )`
	UnaryApplication *UnaryApplication `| @@`
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

func (ast *Unary) GetChildren() []TraversableNode {
	if ast.IsOperation() {
		return []TraversableNode{
			MakeTraversableNodeToken(ast.Op, ast.Pos, ast.EndPos),
			ast.Unary,
		}
	} else if ast.IsUnaryApplication() {
		return []TraversableNode{
			ast.UnaryApplication,
		}
	}
	return []TraversableNode{}
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


func (ast *Unary) Map(mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	if ast.IsOperation() {
		return mapper(&Unary{
			BaseASTNode:      ast.BaseASTNode,
			Op:               ast.Op,
			Unary:            mapper(ast.Unary).(*Unary),
		})
	} else if ast.IsUnaryApplication() {
		return mapper(&Unary{
			BaseASTNode:      ast.BaseASTNode,
			UnaryApplication: mapper(ast.UnaryApplication).(*UnaryApplication),
		})
	}
	panic("Invalid Unary operation type")
}

func (ast *Unary) Visit(mapper hindley_milner.ExpressionMapper) {
	if ast.IsOperation() {
		mapper(ast.Unary)
	} else if ast.IsUnaryApplication() {
		mapper(ast.UnaryApplication)
	}
	mapper(ast)
}

func (ast *Unary) Fn() hindley_milner.Expression {
	return &BuiltinFunction{
		BaseASTNode: ast.BaseASTNode,
		name: ast.Op,
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