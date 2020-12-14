package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type New struct {
	generic_ast.BaseASTNode
	Type       *Type       `"new" ( @@`
	Class      *string       `| @Ident )`
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

func (ast *New) IsTypeConstructor() bool {
	return ast.Type != nil
}

func (ast *New) IsClassConstructor() bool {
	return ast.Class != nil
}

func (ast *New) GetTraversableNode() generic_ast.TraversableNode {
	if ast.IsTypeConstructor() {
		return ast.Type
	} else if ast.IsClassConstructor() {
		return generic_ast.MakeTraversableNodeValue(ast.GetTraversableNode(), *ast.Class, "ident", ast.Pos, ast.EndPos)
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

func (ast *New) GetChildren() []generic_ast.TraversableNode {
	if ast.IsTypeConstructor() {
		return []generic_ast.TraversableNode{
			ast.Type,
		}
	} else if ast.IsClassConstructor() {
		return []generic_ast.TraversableNode{}
	}
	return []generic_ast.TraversableNode{}
}

////

func (ast *New) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	// TODO
	return ast
}

func (ast *New) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) {
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

func (ast *New) Fn() generic_ast.Expression {
	return hindley_milner.ExpressionSignedTupleGet("class", 1, 0, &VarName{
		BaseASTNode: ast.BaseASTNode,
		name: *ast.Class,
	})
}

func (ast *New) Body() generic_ast.Expression {
	return hindley_milner.Batch{
		Exp: []generic_ast.Expression{
			hindley_milner.EmbeddedTypeExpr{
				GetType: func() *hindley_milner.Scheme {
					return hindley_milner.NewScheme(nil, CreatePrimitive(T_VOID),)
				},
			},
		},
	}
}