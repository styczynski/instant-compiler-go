package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Assignment struct {
	BaseASTNode
	TargetName string `@Ident`
	Value *Expression `"=" @@ ";"`
}

func (ast *Assignment) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Assignment) End() lexer.Position {
	return ast.EndPos
}

func (ast *Assignment) GetNode() interface{} {
	return ast
}

func (ast *Assignment) GetChildren() []TraversableNode {
	return []TraversableNode{
		MakeTraversableNodeToken(ast.TargetName, ast.Pos, ast.EndPos),
		ast.Value,
	}
}

func (ast *Assignment) Print(c *context.ParsingContext) string {
	return printNode(c, ast, "%s = %s;", ast.TargetName, ast.Value.Print(c))
}

//

func (ast *Assignment) Map(mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	return mapper(&Assignment{
		BaseASTNode: ast.BaseASTNode,
		Value:    mapper(ast.Value).(*Expression),
	})
}

func (ast *Assignment) Visit(mapper hindley_milner.ExpressionMapper) {
	mapper(ast.Value)
	mapper(ast)
}

func (ast *Assignment) Fn() hindley_milner.Expression {
	return &BuiltinFunction{
		BaseASTNode: ast.BaseASTNode,
		name: "=",
	}
}

func (ast *Assignment) Body() hindley_milner.Expression {
	return hindley_milner.Batch{
		Exp: []hindley_milner.Expression{
			&VarName{
				BaseASTNode: ast.BaseASTNode,
				name: ast.TargetName,
			},
			ast.Value,
		},
	}
}

func (ast *Assignment) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_APPLICATION
}