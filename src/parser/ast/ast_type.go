package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Type struct {
	BaseASTNode
	Name *string `@Ident`
	Dimensions *string `(@( "["`
	Size *Expression `@@? "]" ))?`
}

func (ast *Type) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Type) End() lexer.Position {
	return ast.EndPos
}

func (ast *Type) GetNode() interface{} {
	return ast
}

func (ast *Type) GetChildren() []TraversableNode {
	return []TraversableNode{
		MakeTraversableNodeToken(*ast.Name, ast.Pos, ast.EndPos),
	}
}

func (ast *Type) Print(c *context.ParsingContext) string {
	return printNode(c, ast, "%s", ast.Name)
}


/////

func (ast *Type) GetType() *hindley_milner.Scheme {
	if ast.Dimensions != nil {
		return hindley_milner.NewScheme(nil, hindley_milner.NewSignedTupleType("array", PrimitiveType{
			name:    *ast.Name,
		}))
	}
	return hindley_milner.NewScheme(nil, PrimitiveType{
		name:    *ast.Name,
	})
}