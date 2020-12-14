package ast

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type ForDestructor struct {
	 generic_ast.BaseASTNode
	ElementVar string `@Ident`
	Target *Expression `":" @@`
	ParentNode generic_ast.TraversableNode
}

func (ast *ForDestructor) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *ForDestructor) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *ForDestructor) Begin() lexer.Position {
	return ast.Pos
}

func (ast *ForDestructor) End() lexer.Position {
	return ast.EndPos
}

func (ast *ForDestructor) GetNode() interface{} {
	return ast
}

func (ast *ForDestructor) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{
		ast.Target,
	}
}

func (ast *ForDestructor) Print(c *context.ParsingContext) string {
	return fmt.Sprintf("%s: %s",
		ast.ElementVar,
		ast.Target.Print(c))
}

////


func (ast *ForDestructor) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	return mapper(parent, &ForDestructor{
		BaseASTNode: ast.BaseASTNode,
		ElementVar:  ast.ElementVar,
		Target:      mapper(ast, ast.Target, context).(*Expression),
		ParentNode: parent.(generic_ast.TraversableNode),
	}, context)
}

func (ast *ForDestructor) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) {
	mapper(ast, ast.Target, context)
	mapper(parent, ast, context)
}

func (ast *ForDestructor) Fn() generic_ast.Expression {
	return &BuiltinFunction{
		BaseASTNode: ast.BaseASTNode,
		name: "[_]",
	}
}

func (ast *ForDestructor) Body() generic_ast.Expression {
	return hindley_milner.Batch{
		Exp: []generic_ast.Expression{
			ast.Target,
		},
	}
}

func (ast *ForDestructor) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_APPLICATION
}