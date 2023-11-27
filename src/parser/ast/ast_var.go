package ast

import (
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type VarName struct {
	generic_ast.BaseASTNode
	name         string
	ParentNode   generic_ast.TraversableNode
	ResolvedType hindley_milner.Type
}

func (ast *VarName) OnTypeReturned(t hindley_milner.Type) {
	ast.ResolvedType = t
}

func (ast *VarName) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *VarName) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *VarName) Name() *hindley_milner.NameGroup { return hindley_milner.Name(ast.name) }

func (ast *VarName) Body() generic_ast.Expression { return ast }

func (ast *VarName) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	return mapper(parent, ast, context, false)
}

func (ast *VarName) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(parent, ast, context)
}

func (ast *VarName) Type() hindley_milner.Type {
	return nil
}

func (ast *VarName) ExpressionType() hindley_milner.ExpressionType { return hindley_milner.E_LITERAL }
