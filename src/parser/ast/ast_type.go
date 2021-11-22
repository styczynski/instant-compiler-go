package ast

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Type struct {
	generic_ast.BaseASTNode
	Name       *string   `@Ident`
	Dimensions *Accessor `( @@ )?`
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

func (ast *Type) IsBasePrimitive() bool {
	typeName := *ast.Name
	if typeName == "string" {
		return true
	} else if typeName == "bool" {
		return true
	} else if typeName == "int" {
		return true
	}
	return false
}

func (ast *Type) GetType(c hindley_milner.InferContext) *hindley_milner.Scheme {
	var baseType hindley_milner.Type
	if ast.IsBasePrimitive() {
		baseType = PrimitiveType{
			name: *ast.Name,
		}
	} else {
		if c != nil {
			t, err := c.TypeOf(&New{
				BaseASTNode: ast.BaseASTNode,
				ParentNode:  ast.ParentNode,
				Class:       ast.Name,
			})
			fmt.Printf("GOT CLASS => %s [%v]\n", *ast.Name, t)
			if err != nil {
				panic(err.Error())
			}
			if ast.Dimensions != nil {
				return hindley_milner.NewScheme(hindley_milner.TypeVarSet{}, ast.Dimensions.BuildType(t))
			}
			return hindley_milner.NewScheme(hindley_milner.TypeVarSet{}, t)
		}
		baseType = hindley_milner.NewSignedStructType(*ast.Name, map[string]hindley_milner.Type{})
	}

	if ast.Dimensions != nil {
		return hindley_milner.NewScheme(nil, ast.Dimensions.BuildType(baseType))
	}
	return hindley_milner.NewScheme(nil, baseType)
}
