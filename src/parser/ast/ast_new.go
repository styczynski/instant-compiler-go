package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type New struct {
	generic_ast.BaseASTNode
	Class      *string `"new" @Ident "(" ")"`
	ParentNode generic_ast.TraversableNode
}

func (ast *New) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *New) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *New) Begin() lexer.Position {
	return ast.Pos
}

func (ast *New) End() lexer.Position {
	return ast.EndPos
}

func (ast *New) GetNode() interface{} {
	return ast
}

func (ast *New) GetTraversableNode() generic_ast.TraversableNode {
	return generic_ast.MakeTraversableNodeValue(ast, *ast.Class, "ident", ast.Pos, ast.EndPos)
}

func (ast *New) Print(c *context.ParsingContext) string {
	return *ast.Class
}

func (ast *New) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{}
}

////

func (ast *New) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	return mapper(parent, &New{
		BaseASTNode: ast.BaseASTNode,
		Class:       ast.Class,
		ParentNode:  ast.ParentNode,
	}, context, true)
}

func (ast *New) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	// TODO
}

func (ast *New) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_APPLICATION
}

func (ast *New) Fn(c hindley_milner.InferContext) generic_ast.Expression {
	return hindley_milner.ExpressionSignedTupleGet("class", 1, 0, &VarName{
		BaseASTNode: ast.BaseASTNode,
		name:        *ast.Class,
	})
}

func (ast *New) Body() generic_ast.Expression {
	return hindley_milner.Batch{
		Exp: []generic_ast.Expression{
			hindley_milner.EmbeddedTypeExpr{
				GetType: func() *hindley_milner.Scheme {
					return hindley_milner.NewScheme(nil, CreatePrimitive(T_VOID))
				},
				Source: ast,
			},
		},
	}
}
