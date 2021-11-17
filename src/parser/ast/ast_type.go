package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Type struct {
	generic_ast.BaseASTNode
	Name *string `@( "string" | "boolean" | "int" | "void" )`
	Dimensions *string `(@( "["`
	Size *Expression `@@? "]" ))?`
	ParentNode generic_ast.TraversableNode
}

func (ast *Type) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *Type) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
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

func (ast *Type) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{
		generic_ast.MakeTraversableNodeToken(ast, *ast.Name, ast.Pos, ast.EndPos),
	}
}

func (ast *Type) Print(c *context.ParsingContext) string {
	return printNode(c, ast, "%s", *ast.Name)
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