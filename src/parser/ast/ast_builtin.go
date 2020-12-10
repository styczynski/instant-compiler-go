package ast

import (
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type BuiltinFunction struct {
	name string
}

func (ast *BuiltinFunction) Name() hindley_milner.NameGroup     { return hindley_milner.Name(ast.name) }

func (ast *BuiltinFunction) Body() hindley_milner.Expression { return ast }

func (ast *BuiltinFunction) Map(mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	return mapper(ast)
}

func (ast *BuiltinFunction) Visit(mapper hindley_milner.ExpressionMapper) {
	mapper(ast)
}

func (ast *BuiltinFunction) Type() hindley_milner.Type {
	return nil
}

func (ast *BuiltinFunction) ExpressionType() hindley_milner.ExpressionType { return hindley_milner.E_LITERAL }

