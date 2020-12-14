package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type For struct {
	generic_ast.BaseASTNode
	ElementType *Type `"for" "(" @@`
	Destructor *ForDestructor `@@ ")"`
	Do *Statement `@@`
	ParentNode generic_ast.TraversableNode
}

func (ast *For) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *For) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *For) Begin() lexer.Position {
	return ast.Pos
}

func (ast *For) End() lexer.Position {
	return ast.EndPos
}

func (ast *For) GetNode() interface{} {
	return ast
}

func (ast *For) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{
		ast.Destructor,
		ast.Do,
	}
}

func (ast *For) Print(c *context.ParsingContext) string {
	return printNode(c, ast, "for (%s %s: %s) %s",
		ast.ElementType.Print(c),
		ast.Destructor.Print(c),
		ast.Do.Print(c),
	)
}

///

func (ast *For) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	// TODO
	return ast
}
func (ast *For) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) {
	mapper(ast, ast.Destructor, context)
	mapper(ast, ast.Do, context)
	mapper(parent, ast, context)
}

func (ast *For) Var() hindley_milner.NameGroup {
	types := map[string]*hindley_milner.Scheme{}
	types[ast.Destructor.ElementVar] = ast.ElementType.GetType()
	return hindley_milner.NamesWithTypes([]string{ ast.Destructor.ElementVar }, types)
}

func (ast *For) Def() generic_ast.Expression {
	return ast.Destructor
}

func (ast *For) Body() generic_ast.Expression {
	return hindley_milner.Batch{Exp: []generic_ast.Expression{
		ast.Do,
	}}
}

func (ast *For) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_LET_RECURSIVE
}
