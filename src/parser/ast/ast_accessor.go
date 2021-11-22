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
	IndexingExpr *Expression `( "[" @@ "]"`
	Property     *string     ` | "." @Ident )`
	Next         *Accessor   `@@?`
	ParentNode   generic_ast.TraversableNode
}

func (ast *Accessor) IsIndex() bool {
	return ast.IndexingExpr != nil
}

func (ast *Accessor) IsProperty() bool {
	return ast.Property != nil
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

func (ast *Accessor) GetIndexingNode() generic_ast.Expression {
	if ast.IsIndex() {
		return ast.IndexingExpr
	} else if ast.IsProperty() {
		return &Primary{
			BaseASTNode: ast.BaseASTNode,
			String:      ast.Property,
		}
	} else {
		panic("Invalid accessor")
	}
}

func (ast *Accessor) GetChildren() []generic_ast.TraversableNode {
	cur := ast.GetIndexingNode()
	if ast.HasNext() {
		return []generic_ast.TraversableNode{
			cur.(generic_ast.TraversableNode),
			ast.Next,
		}
	} else {
		return []generic_ast.TraversableNode{
			cur.(generic_ast.TraversableNode),
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
	if ast.IsIndex() {
		return fmt.Sprintf("[%s]%s", ast.IndexingExpr.Print(c), nextStr)
	} else if ast.IsProperty() {
		return fmt.Sprintf(".%s%s", *ast.Property, nextStr)
	}
	panic("Invalid accessor")
}

////

func (ast *Accessor) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	next := ast.Next
	if ast.HasNext() {
		next = mapper(ast, next, context, false).(*Accessor)
	}
	if ast.IsIndex() {
		return mapper(parent, &Accessor{
			BaseASTNode:  ast.BaseASTNode,
			IndexingExpr: mapper(ast, ast.IndexingExpr, context, false).(*Expression),
			Next:         next,
			ParentNode:   parent.(generic_ast.TraversableNode),
		}, context, true)
	} else if ast.IsProperty() {
		return mapper(parent, ast, context, true)
	}
	panic("Invalid accessor")
}

func (ast *Accessor) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(ast, ast.GetIndexingNode(), context)
	if ast.HasNext() {
		mapper(ast, ast.Next, context)
	}
	mapper(parent, ast, context)
}

func (ast *Accessor) Fn(c hindley_milner.InferContext) generic_ast.Expression {
	operatorName := ""
	if ast.IsIndex() {
		return &BuiltinFunction{
			BaseASTNode: ast.BaseASTNode,
			name:        "[]",
		}
	} else if ast.IsProperty() {
		return &hindley_milner.EmbeddedTypeExpr{
			Source: ast,
			GetType: func() *hindley_milner.Scheme {
				return hindley_milner.NewScheme(
					hindley_milner.TypeVarSet{hindley_milner.TVar('a')},
					hindley_milner.NewFnType(
						hindley_milner.NewSignedStructType("", map[string]hindley_milner.Type{
							*ast.Property: hindley_milner.TVar('a'),
						}),
						CreatePrimitive(T_STRING),
						hindley_milner.TVar('a'),
					))
			},
		}
	} else {
		panic("Invalid accessor")
	}
	return &BuiltinFunction{
		BaseASTNode: ast.BaseASTNode,
		name:        operatorName,
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
	cur := ast.GetIndexingNode()
	if ast.IsTop() {
		return hindley_milner.Batch{
			Exp: []generic_ast.Expression{
				ast.ParentNode.(*Index).Primary,
				cur,
			},
		}
	} else {
		return hindley_milner.Batch{
			Exp: []generic_ast.Expression{
				ast.ParentNode.(generic_ast.Expression),
				cur,
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
