package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type While struct {
	 generic_ast.BaseASTNode
	Condition *Expression `"while" "(" @@ ")"`
	Do *Statement `@@`
	ParentNode generic_ast.TraversableNode
}

func (ast *While) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *While) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *While) Begin() lexer.Position {
	return ast.Pos
}

func (ast *While) End() lexer.Position {
	return ast.EndPos
}

func (ast *While) GetNode() interface{} {
	return ast
}

func (ast *While) Print(c *context.ParsingContext) string {
	c.PrinterConfiguration.SkipStatementIdent = true
	body := ast.Do.Print(c)
	return printNode(c, ast, "while (%s) %s", ast.Condition.Print(c), body)
}

func (ast *While) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{
		ast.Condition,
		ast.Do,
	}
}

///

func (ast *While) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	return mapper(parent, &While{
		BaseASTNode: ast.BaseASTNode,
		Condition: mapper(ast, ast.Condition, context).(*Expression),
		Do: mapper(ast, ast.Do, context).(*Statement),
		ParentNode: parent.(generic_ast.TraversableNode),
	}, context)
}

func (ast *While) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) {
	mapper(ast, ast.Condition, context)
	mapper(ast, ast.Do, context)
	mapper(parent, ast, context)
}

func (ast *While) Fn() generic_ast.Expression {
	return &hindley_milner.EmbeddedTypeExpr{GetType: func() *hindley_milner.Scheme {
		return hindley_milner.NewScheme(
			hindley_milner.TypeVarSet{hindley_milner.TVar('a')},
			hindley_milner.NewFnType(CreatePrimitive(T_BOOL), hindley_milner.TVar('a'), CreatePrimitive(T_VOID)))
	}}
}

func (ast *While) Body() generic_ast.Expression {
	args := []generic_ast.Expression{
		ast.Condition,
		ast.Do,
	}
	return hindley_milner.Batch{
		Exp: args,
	}
}

func (ast *While) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_APPLICATION
}