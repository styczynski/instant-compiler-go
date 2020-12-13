package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type BuiltinFunction struct {
	generic_ast.BaseASTNode
	name string
	ParentNode generic_ast.TraversableNode
}

func (ast *BuiltinFunction) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *BuiltinFunction) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *BuiltinFunction) Begin() lexer.Position {
	return ast.Pos
}

func (ast *BuiltinFunction) End() lexer.Position {
	return ast.EndPos
}

func (ast *BuiltinFunction) GetNode() interface{} {
	return ast
}

func (ast *BuiltinFunction) Name() hindley_milner.NameGroup     { return hindley_milner.Name(ast.name) }

func (ast *BuiltinFunction) Body() hindley_milner.Expression { return ast }

func (ast *BuiltinFunction) Map(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	return mapper(parent, ast)
}

func (ast *BuiltinFunction) Visit(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) {
	mapper(parent, ast)
}

func (ast *BuiltinFunction) Type() hindley_milner.Type {
	return nil
}

func (ast *BuiltinFunction) ExpressionType() hindley_milner.ExpressionType { return hindley_milner.E_LITERAL }

