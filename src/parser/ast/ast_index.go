package ast

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Index struct {
	generic_ast.BaseASTNode
	Primary   *Primary   ` @@ `
	IndexingExpr *Expression `( "[" @@ "]" )?`
	ParentNode generic_ast.TraversableNode
}

func (ast *Index) ExtractConst() (generic_ast.TraversableNode, bool) {
	if ast.HasIndexingExpr() {
		return nil, false
	}
	return ast.Primary.ExtractConst()
}

func (ast *Index) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *Index) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *Index) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Index) End() lexer.Position {
	return ast.EndPos
}

func (ast *Index) GetNode() interface{} {
	return ast
}

func (ast *Index) GetChildren() []generic_ast.TraversableNode {
	if ast.HasIndexingExpr() {
		return []generic_ast.TraversableNode{
			ast.Primary,
			ast.IndexingExpr,
		}
	} else {
		return []generic_ast.TraversableNode{
			ast.Primary,
		}
	}
}

func (ast *Index) HasIndexingExpr() bool {
	return ast.IndexingExpr != nil
}

func (ast *Index) Print(c *context.ParsingContext) string {
	if ast.HasIndexingExpr() {
		return fmt.Sprintf("%s[%s]", ast.Primary.Print(c), ast.IndexingExpr.Print(c))
	}
	return ast.Primary.Print(c)
}

////


func (ast *Index) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	if ast.HasIndexingExpr() {
		return mapper(parent, &Index{
			BaseASTNode:      ast.BaseASTNode,
			Primary: mapper(ast, ast.Primary, context, false).(*Primary),
			IndexingExpr: mapper(ast, ast.IndexingExpr, context, false).(*Expression),
			ParentNode: parent.(generic_ast.TraversableNode),
		}, context, true)
	} else {
		return mapper(parent, &Index{
			BaseASTNode:      ast.BaseASTNode,
			Primary: mapper(ast, ast.Primary, context, false).(*Primary),
			ParentNode: parent.(generic_ast.TraversableNode),
		}, context, true)
	}
}

func (ast *Index) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(ast, ast.Primary, context)
	if ast.HasIndexingExpr() {
		mapper(ast, ast.IndexingExpr, context)
	}
	mapper(parent, ast, context)
}

func (ast *Index) Fn() generic_ast.Expression {
	return &BuiltinFunction{
		BaseASTNode: ast.BaseASTNode,
		name: "[]",
	}
}

func (ast *Index) Body() generic_ast.Expression {
	if !ast.HasIndexingExpr() {
		return ast.Primary
	}
	return hindley_milner.Batch{
		Exp: []generic_ast.Expression{
			ast.Primary,
			ast.IndexingExpr,
		},
	}
}

func (ast *Index) ExpressionType() hindley_milner.ExpressionType {
	if !ast.HasIndexingExpr() {
		return hindley_milner.E_PROXY
	}
	return hindley_milner.E_APPLICATION
}