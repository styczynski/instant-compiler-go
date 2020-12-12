package ast

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Index struct {
	BaseASTNode
	Primary   *Primary   ` @@ `
	IndexingExpr *Expression `( "[" @@ "]" )?`
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

func (ast *Index) GetChildren() []TraversableNode {
	if ast.HasIndexingExpr() {
		return []TraversableNode{
			ast.Primary,
			ast.IndexingExpr,
		}
	} else {
		return []TraversableNode{
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


func (ast *Index) Map(mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	if ast.HasIndexingExpr() {
		return mapper(&Index{
			BaseASTNode:      ast.BaseASTNode,
			Primary: mapper(ast.Primary).(*Primary),
			IndexingExpr: mapper(ast.IndexingExpr).(*Expression),
		})
	} else {
		return mapper(&Index{
			BaseASTNode:      ast.BaseASTNode,
			Primary: mapper(ast.Primary).(*Primary),
		})
	}
}

func (ast *Index) Visit(mapper hindley_milner.ExpressionMapper) {
	mapper(ast.Primary)
	if ast.HasIndexingExpr() {
		mapper(ast.IndexingExpr)
	}
	mapper(ast)
}

func (ast *Index) Fn() hindley_milner.Expression {
	return &BuiltinFunction{
		BaseASTNode: ast.BaseASTNode,
		name: "[]",
	}
}

func (ast *Index) Body() hindley_milner.Expression {
	if !ast.HasIndexingExpr() {
		return ast.Primary
	}
	return hindley_milner.Batch{
		Exp: []hindley_milner.Expression{
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