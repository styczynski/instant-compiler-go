package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type New struct {
	BaseASTNode
	Type       *Type       `"new" ( @@`
	Class      *string       `| @Ident )`
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

func (ast *New) IsTypeConstructor() bool {
	return ast.Type != nil
}

func (ast *New) IsClassConstructor() bool {
	return ast.Class != nil
}

func (ast *New) GetTraversableNode() TraversableNode {
	if ast.IsTypeConstructor() {
		return ast.Type
	} else if ast.IsClassConstructor() {
		return MakeTraversableNodeValue(*ast.Class, "ident", ast.Pos, ast.EndPos)
	}
	panic("Invalid New type")
}

func (ast *New) Print(c *context.ParsingContext) string {
	if ast.IsTypeConstructor() {
		return ast.Type.Print(c)
	} else if ast.IsClassConstructor() {
		return *ast.Class
	}
	panic("Invalid New type")
}

////

func (ast *New) Map(mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	// TODO
	return ast
}

func (ast *New) Visit(mapper hindley_milner.ExpressionMapper) {
	// TODO
}

func (ast *New) EmbeddedType() *hindley_milner.Scheme {
	return ast.Type.GetType()
}

func (ast *New) ExpressionType() hindley_milner.ExpressionType {
	if ast.IsTypeConstructor() {
		return hindley_milner.E_TYPE
	} else if ast.IsClassConstructor() {
		return hindley_milner.E_APPLICATION
	}
	panic("Invalid New type")
}

func (ast *New) Fn() hindley_milner.Expression {
	return hindley_milner.ExpressionSignedTupleGet("class", 1, 0, &VarName{
		BaseASTNode: ast.BaseASTNode,
		name: *ast.Class,
	})
}

func (ast *New) Body() hindley_milner.Expression {
	return hindley_milner.Batch{
		Exp: []hindley_milner.Expression{
			hindley_milner.EmbeddedTypeExpr{
				GetType: func() *hindley_milner.Scheme {
					return hindley_milner.NewScheme(nil, CreatePrimitive(T_VOID),)
				},
			},
		},
	}
}