package ast

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type ClassField struct {
	generic_ast.BaseASTNode
	Method         *FnDef  `@@`
	ClassFieldType *Type   `| ( @@`
	Name           *string `@Ident ";" )`
	ParentNode     generic_ast.TraversableNode
}

func (ast *ClassField) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *ClassField) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *ClassField) Begin() lexer.Position {
	return ast.Pos
}

func (ast *ClassField) End() lexer.Position {
	return ast.EndPos
}

func (ast *ClassField) GetNode() interface{} {
	return ast
}

func (ast *ClassField) IsField() bool {
	return ast.Name != nil
}

func (ast *ClassField) IsMethod() bool {
	return ast.Method != nil
}

func (ast *ClassField) GetChildren() []generic_ast.TraversableNode {
	if ast.IsMethod() {
		return []generic_ast.TraversableNode{ast.Method}
	} else if ast.IsField() {
		return []generic_ast.TraversableNode{
			ast.ClassFieldType,
			generic_ast.MakeTraversableNodeToken(ast, *ast.Name, ast.Pos, ast.EndPos),
		}
	}
	panic("Invalid ClassField type")
}

func (ast *ClassField) Print(c *context.ParsingContext) string {
	if ast.IsMethod() {
		return printNode(c, ast, "%s", ast.Method.Print(c))
	} else if ast.IsField() {
		return printNode(c, ast, "%s %s", ast.ClassFieldType.Print(c), *ast.Name)
	}
	panic("Invalid ClassField type")
}

func (ast *ClassField) FieldName() string {
	if ast.IsMethod() {
		return ast.Method.Name
	} else if ast.IsField() {
		return *ast.Name
	}
	panic("Invalid ClassField type")
}

func (ast *ClassField) GetType(c hindley_milner.InferContext) *hindley_milner.Scheme {
	if ast.IsMethod() {
		methodType, err := c.TypeOf(&VarName{
			name: ast.Method.Name,
		}, ast.Method)
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("class method type ==> %s\n", methodType)
		//return hindley_milner.NewScheme(nil, CreatePrimitive(T_VOID))
		return hindley_milner.NewScheme(hindley_milner.TypeVarSet{}, methodType)
	} else {
		return ast.ClassFieldType.GetType(c)
	}
}
