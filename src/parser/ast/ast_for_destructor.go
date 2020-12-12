package ast

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type ForDestructor struct {
	BaseASTNode
	ElementVar string `@Ident`
	Target *Expression `":" @@`
	ParentNode TraversableNode
}

func (ast *ForDestructor) Parent() TraversableNode {
	return ast.ParentNode
}

func (ast *ForDestructor) OverrideParent(node TraversableNode) {
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

func (ast *ForDestructor) GetChildren() []TraversableNode {
	return []TraversableNode{
		ast.Target,
	}
}

func (ast *ForDestructor) Print(c *context.ParsingContext) string {
	return fmt.Sprintf("%s: %s",
		ast.ElementVar,
		ast.Target.Print(c))
}

////


func (ast *ForDestructor) Map(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	return mapper(parent, &ForDestructor{
		BaseASTNode: ast.BaseASTNode,
		ElementVar:  ast.ElementVar,
		Target:      mapper(ast, ast.Target).(*Expression),
		ParentNode: parent.(TraversableNode),
	})
}

func (ast *ForDestructor) Visit(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) {
	mapper(ast, ast.Target)
	mapper(parent, ast)
}

func (ast *ForDestructor) Fn() hindley_milner.Expression {
	return &BuiltinFunction{
		BaseASTNode: ast.BaseASTNode,
		name: "[_]",
	}
}

func (ast *ForDestructor) Body() hindley_milner.Expression {
	return hindley_milner.Batch{
		Exp: []hindley_milner.Expression{
			ast.Target,
		},
	}
}

func (ast *ForDestructor) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_APPLICATION
}