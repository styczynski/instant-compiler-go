package ast

import (
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type VarName struct {
	BaseASTNode
	name string
}

func (ast *VarName) Name() hindley_milner.NameGroup     { return hindley_milner.Name(ast.name) }

func (ast *VarName) Body() hindley_milner.Expression { return ast }

func (ast *VarName) Map(mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	return mapper(ast)
}

func (ast *VarName) Visit(mapper hindley_milner.ExpressionMapper) {
	mapper(ast)
}

func (ast *VarName) Type() hindley_milner.Type {
	return nil
}

func (ast *VarName) ExpressionType() hindley_milner.ExpressionType { return hindley_milner.E_LITERAL }

