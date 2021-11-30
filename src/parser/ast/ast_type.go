package ast

import (
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

func IsTypeBasePrimitive(name *string) bool {
	typeName := *name
	if typeName == "string" {
		return true
	} else if typeName == "bool" {
		return true
	} else if typeName == "int" {
		return true
	} else if typeName == "void" {
		return true
	}
	return false
}

func (ast *Type) GetType(c hindley_milner.InferContext) *hindley_milner.Scheme {
	var baseType hindley_milner.Type
	if *ast.Name == "auto" {
		return hindley_milner.TypeHelperAny()
	} else if *ast.Name == "class" {
		return hindley_milner.NewScheme(nil, hindley_milner.NewSignedTupleType("class", hindley_milner.NewFnType(
			CreatePrimitive(T_VOID),
			hindley_milner.NewSignedStructType("", map[string]hindley_milner.Type{}),
		)))
	} else if IsTypeBasePrimitive(ast.Name) {
		baseType = PrimitiveType{
			name: *ast.Name,
		}
	} else {
		baseType = hindley_milner.NewSignedStructType(*ast.Name, map[string]hindley_milner.Type{})
	}

	if ast.Dimensions != nil {
		return hindley_milner.NewScheme(nil, ast.Dimensions.BuildType(baseType))
	}
	return hindley_milner.NewScheme(nil, baseType)
}
