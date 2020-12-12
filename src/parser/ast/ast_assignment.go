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
	ParentNode TraversableNode
}

func (ast *Assignment) Parent() TraversableNode {
	return ast.ParentNode
}

func (ast *Assignment) OverrideParent(node TraversableNode) {
	ast.ParentNode = node
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
		MakeTraversableNodeToken(ast, ast.TargetName, ast.Pos, ast.EndPos),
		ast.Value,
	}
}

func (ast *Assignment) Print(c *context.ParsingContext) string {
	return printNode(c, ast, "%s = %s;", ast.TargetName, ast.Value.Print(c))
}

//

func (ast *Assignment) Map(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	return mapper(parent, &Assignment{
		BaseASTNode: ast.BaseASTNode,
		Value:    mapper(ast, ast.Value).(*Expression),
		ParentNode: parent.(TraversableNode),
	})
}

func (ast *Assignment) Visit(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) {
	mapper(ast, ast.Value)
	mapper(parent, ast)
}

func (ast *Assignment) Fn() hindley_milner.Expression {
	//return &BuiltinFunction{
	//	BaseASTNode: ast.BaseASTNode,
	//	name: "=",
	//}
	return &hindley_milner.EmbeddedTypeExpr{GetType: func() *hindley_milner.Scheme {
		return hindley_milner.NewScheme(
			hindley_milner.TypeVarSet{hindley_milner.TVar('a')},
			hindley_milner.NewFnType(hindley_milner.TVar('a'), hindley_milner.TVar('a'), hindley_milner.TVar('a')))
	}}
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