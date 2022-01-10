package ir

import (
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type IREmpty struct {
	generic_ast.BaseASTNode
	content string `@(";")`
}

func (ast *IREmpty) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_BLOCK
}

func (ast *IREmpty) Parent() generic_ast.TraversableNode {
	return nil
}

func (ast *IREmpty) Print(c *context.ParsingContext) string {
	return ""
}

func (ast *IREmpty) OverrideParent(node generic_ast.TraversableNode) {

}

func (ast *IREmpty) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	return mapper(parent, &IREmpty{
		content: ast.content,
	}, context, true)
}

func (ast *IREmpty) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(parent, ast, context)
}

func (ast *IREmpty) GetContents() hindley_milner.Batch {
	return hindley_milner.Batch{}
}

func (ast *IREmpty) IsBlock() bool {
	return true
}

func (ast *IREmpty) Expressions() []generic_ast.Expression {
	return []generic_ast.Expression{}
}

func (ast *IREmpty) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{}
}

func (ast *IREmpty) GetNode() interface{} {
	return ast
}

func (ast *IREmpty) Body() generic_ast.Expression {
	return hindley_milner.Batch{}
}
