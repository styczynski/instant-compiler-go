package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type For struct {
	BaseASTNode
	ElementType *Type `"for" "(" @@`
	Destructor *ForDestructor `@@ ")"`
	Do *Statement `@@`
	ParentNode TraversableNode
}

func (ast *For) Parent() TraversableNode {
	return ast.ParentNode
}

func (ast *For) OverrideParent(node TraversableNode) {
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

func (ast *For) GetChildren() []TraversableNode {
	return []TraversableNode{
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

func (ast *For) Map(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	// TODO
	return ast
}
func (ast *For) Visit(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) {
	mapper(ast, ast.Destructor)
	mapper(ast, ast.Do)
	mapper(parent, ast)
}

func (ast *For) Var() hindley_milner.NameGroup {
	types := map[string]*hindley_milner.Scheme{}
	types[ast.Destructor.ElementVar] = ast.ElementType.GetType()
	return hindley_milner.NamesWithTypes([]string{ ast.Destructor.ElementVar }, types)
}

func (ast *For) Def() hindley_milner.Expression {
	return ast.Destructor
}

func (ast *For) Body() hindley_milner.Expression {
	return hindley_milner.Batch{Exp: []hindley_milner.Expression{
		ast.Do,
	}}
}

func (ast *For) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_LET_RECURSIVE
}
