package ast

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Accessor struct {
	generic_ast.BaseASTNode
	IndexingExpr *Expression `"[" @@ "]"`
	Next         *Accessor   `@@?`
	ParentNode   generic_ast.TraversableNode
}

func (ast *Accessor) ExtractConst() (generic_ast.TraversableNode, bool) {
	return nil, false
}

func (ast *Accessor) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *Accessor) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *Accessor) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Accessor) End() lexer.Position {
	return ast.EndPos
}

func (ast *Accessor) GetNode() interface{} {
	return ast
}

func (ast *Accessor) GetChildren() []generic_ast.TraversableNode {
	if ast.HasNext() {
		return []generic_ast.TraversableNode{
			ast.IndexingExpr,
			ast.Next,
		}
	} else {
		return []generic_ast.TraversableNode{
			ast.IndexingExpr,
		}
	}
}

func (ast *Accessor) HasNext() bool {
	return ast.Next != nil
}

func (ast *Accessor) Print(c *context.ParsingContext) string {
	nextStr := ""
	if ast.HasNext() {
		return ast.Next.Print(c)
	}
	return fmt.Sprintf("[%s]%s", ast.IndexingExpr.Print(c), nextStr)
}

////

func (ast *Accessor) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	if ast.HasNext() {
		return mapper(parent, &Accessor{
			BaseASTNode:  ast.BaseASTNode,
			IndexingExpr: mapper(ast, ast.IndexingExpr, context, false).(*Expression),
			Next:         mapper(ast, ast.Next, context, false).(*Accessor),
			ParentNode:   parent.(generic_ast.TraversableNode),
		}, context, true)
	} else {
		return mapper(parent, &Accessor{
			BaseASTNode:  ast.BaseASTNode,
			IndexingExpr: mapper(ast, ast.IndexingExpr, context, false).(*Expression),
			ParentNode:   parent.(generic_ast.TraversableNode),
		}, context, true)
	}
}

func (ast *Accessor) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(ast, ast.IndexingExpr, context)
	if ast.HasNext() {
		mapper(ast, ast.Next, context)
	}
	mapper(parent, ast, context)
}

func (ast *Accessor) Fn() generic_ast.Expression {
	return &BuiltinFunction{
		BaseASTNode: ast.BaseASTNode,
		name:        "[]",
	}
}

func (ast *Accessor) IsTop() bool {
	_, parentIsAccessor := ast.ParentNode.(*Accessor)
	return !parentIsAccessor
}

func (ast *Accessor) GetLastAccessor() generic_ast.Expression {
	if ast.HasNext() {
		return ast.Next.GetLastAccessor()
	}
	return ast
}

func (ast *Accessor) Body() generic_ast.Expression {
	/*
	 * index [target] -> accessor [index1] -> accessor [index2]
	 *    <proxy>
	 *               getel(target, index1)
	 *                                getel(getel(target, index1), index2)
	 */
	if ast.IsTop() {
		return hindley_milner.Batch{
			Exp: []generic_ast.Expression{
				ast.ParentNode.(*Index).Primary,
				ast.IndexingExpr,
			},
		}
	} else {
		return hindley_milner.Batch{
			Exp: []generic_ast.Expression{
				ast.ParentNode.(generic_ast.Expression),
				ast.IndexingExpr,
			},
		}
	}
}

func (ast *Accessor) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_APPLICATION
}

func (ast *Accessor) BuildType(target hindley_milner.Type) *hindley_milner.SignedTuple {
	if ast.HasNext() {
		return hindley_milner.NewSignedTupleType("array", ast.Next.BuildType(target))
	}
	return hindley_milner.NewSignedTupleType("array", target)
}
