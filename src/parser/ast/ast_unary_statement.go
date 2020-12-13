package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type UnaryStatement struct {
	generic_ast.BaseASTNode
	TargetName *string `@Ident`
	Operation string `@( "+" "+" | "-" "-" ) ";"`
	ParentNode generic_ast.TraversableNode
}

func (ast *UnaryStatement) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *UnaryStatement) OverrideParent(node generic_ast.TraversableNode) {
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

func (ast *UnaryStatement) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{
		generic_ast.MakeTraversableNodeToken(ast, *ast.TargetName, ast.Pos, ast.EndPos),
		generic_ast.MakeTraversableNodeToken(ast, ast.Operation, ast.Pos, ast.EndPos),
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
		ParentNode: parent.(generic_ast.TraversableNode),
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