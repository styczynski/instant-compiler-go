package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type UnaryStatement struct {
	BaseASTNode
	TargetName *string `@Ident`
	Operation string `@( "+" "+" | "-" "-" ) ";"`
	ParentNode TraversableNode
}

func (ast *UnaryStatement) Parent() TraversableNode {
	return ast.ParentNode
}

func (ast *UnaryStatement) OverrideParent(node TraversableNode) {
	ast.ParentNode = node
}

func (ast *UnaryStatement) Begin() lexer.Position {
	return ast.Pos
}

func (ast *UnaryStatement) End() lexer.Position {
	return ast.EndPos
}

func (ast *UnaryStatement) GetNode() interface{} {
	return ast
}

func (ast *UnaryStatement) GetChildren() []TraversableNode {
	return []TraversableNode{
		MakeTraversableNodeToken(ast, *ast.TargetName, ast.Pos, ast.EndPos),
		MakeTraversableNodeToken(ast, ast.Operation, ast.Pos, ast.EndPos),
	}
}

func (ast *UnaryStatement) Print(c *context.ParsingContext) string {
	return printNode(c, ast, "%s%s;", *ast.TargetName, ast.Operation)
}

///

func (ast *UnaryStatement) Map(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	return mapper(parent, &UnaryStatement{
		BaseASTNode: ast.BaseASTNode,
		TargetName: ast.TargetName,
		Operation: ast.Operation,
		ParentNode: parent.(TraversableNode),
	})
}

func (ast *UnaryStatement) Visit(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) {
	mapper(parent, ast)
}

func (ast *UnaryStatement) Fn() hindley_milner.Expression {
	return &BuiltinFunction{
		BaseASTNode: ast.BaseASTNode,
		name: ast.Operation,
	}
}

func (ast *UnaryStatement) Body() hindley_milner.Expression {
	return hindley_milner.Batch{
		Exp: []hindley_milner.Expression{
			&VarName{
				BaseASTNode: ast.BaseASTNode,
				name: *ast.TargetName,
			},
		},
	}
}

func (ast *UnaryStatement) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_APPLICATION
}