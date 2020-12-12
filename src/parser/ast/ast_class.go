package ast

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Class struct {
	BaseASTNode
	Name string `"class" @Ident "{"`
	Fields []*ClassField `(@@ ";")* "}"`
}

func (ast *Class) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Class) End() lexer.Position {
	return ast.EndPos
}

func (ast *Class) GetNode() interface{} {
	return ast
}

func (ast *Class) GetChildren() []TraversableNode {
	ret := []TraversableNode{
		MakeTraversableNodeToken(ast.Name, ast.Pos, ast.EndPos),
	}
	for _, field := range ast.Fields {
		ret = append(ret, field)
	}
	return ret
}

func (ast *Class) Print(c *context.ParsingContext) string {
	classContents := []string{}
	for _, field := range ast.Fields {
		classContents = append(classContents, fmt.Sprintf("%s;", field.Print(c)))
	}
	return printNode(c, ast, "class %s\n{\n    %s\n}\n", ast.Name,strings.Join(classContents, "\n    "))
}

////

func (ast *Class) Body() hindley_milner.Expression {
	return ast
}

func (ast *Class) Map(mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	return mapper(&Class{
		BaseASTNode: ast.BaseASTNode,
		Name: ast.Name,
		Fields: ast.Fields,
	}).(*Class)
}

func (ast *Class) Visit(mapper hindley_milner.ExpressionMapper) {
	mapper(ast)
}

func (ast *Class) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_DECLARATION
}

func (ast *Class) Var() hindley_milner.NameGroup {
	return hindley_milner.NamesWithTypes([]string{
		ast.Name,
	}, map[string]*hindley_milner.Scheme{

	})
	//names := []string{}
	//types := map[string]*hindley_milner.Scheme{}
	//for _, item := range ast.Fields {
	//	names = append(names, item.Name)
	//	types[item.Name] = ast.DeclarationType.GetType()
	//}
	//return hindley_milner.NamesWithTypes(names, types)
}

func (ast *Class) GetClassInstanceType() hindley_milner.Type {
	fields := map[string]hindley_milner.Type{}
	for _, field := range ast.Fields {
		t, _ := field.ClassFieldType.GetType().Type()
		fields[field.Name] = t
	}
	return hindley_milner.NewSignedStructType(ast.Name, fields)
}

func (ast *Class) Def() hindley_milner.Expression {
	return hindley_milner.EmbeddedTypeExpr{
		GetType: func() *hindley_milner.Scheme {
			return hindley_milner.NewScheme(nil, hindley_milner.NewSignedTupleType("class", hindley_milner.NewFnType(
				CreatePrimitive(T_VOID),
				ast.GetClassInstanceType(),
			)))
		},
	}
}