package ast

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Typename struct {
	generic_ast.BaseASTNode
	Expr       *LogicalOperation `"typename" @@`
	ParentNode generic_ast.TraversableNode
}

func (ast *Typename) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *Typename) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *Typename) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Typename) End() lexer.Position {
	return ast.EndPos
}

func (ast *Typename) GetNode() interface{} {
	return ast
}

func (ast *Typename) GetTraversableNode() generic_ast.TraversableNode {
	return ast.Expr
}

func (ast *Typename) Print(c *context.ParsingContext) string {
	return fmt.Sprintf("typename %s", ast.Expr.Print(c))
}

func (ast *Typename) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{
		ast.Expr,
	}
}

////

func (ast *Typename) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	return mapper(parent, &Typename{
		BaseASTNode: ast.BaseASTNode,
		Expr:        mapper(ast, ast.Expr, context, false).(*LogicalOperation),
		ParentNode:  ast.ParentNode,
	}, context, true)
}

func (ast *Typename) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	// TODO
	mapper(ast, ast.Expr, context)
	mapper(parent, ast, context)
}

func (ast *Typename) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_INTROSPECTION
}

func (ast *Typename) Body() generic_ast.Expression {
	return ast.Expr
}

func (ast *Typename) OnTypeReturned(t hindley_milner.Type) {
	parent := ast.Parent()
	expr := (parent.(interface{}).(*Expression))
	expr.LogicalOperation = nil
	expr.Typename = nil
	expr.NewType = nil

	val := fmt.Sprintf("%v", t)
	

	newAST := &LogicalOperation{
		BaseASTNode: ast.BaseASTNode,
		Equality: &Equality{
			BaseASTNode: ast.BaseASTNode,
			Comparison: &Comparison{
				BaseASTNode: ast.BaseASTNode,
				Addition: &Addition{
					BaseASTNode: ast.BaseASTNode,
					Multiplication: &Multiplication{
						BaseASTNode: ast.BaseASTNode,
						Unary: &Unary{
							BaseASTNode: ast.BaseASTNode,
							Op:          "",
							Unary: &Unary{
								BaseASTNode: ast.BaseASTNode,
								Op:          "",
								Unary:       nil,
								UnaryApplication: &UnaryApplication{
									BaseASTNode: ast.BaseASTNode,
									Arguments:   nil,
									Index: &Index{
										BaseASTNode: ast.BaseASTNode,
										Primary: &Primary{
											BaseASTNode:   ast.BaseASTNode,
											Variable:      nil,
											Int:           nil,
											String:        &val,
											Bool:          nil,
											SubExpression: nil,
											ParentNode:    nil,
										},
										IndexingExpr: nil,
										ParentNode:   nil,
									},
									ParentNode: nil,
								},
								ParentNode: nil,
							},
							UnaryApplication: nil,
							ParentNode:       nil,
						},
						Op:         "",
						Next:       nil,
						ParentNode: nil,
					},
					Op:         "",
					Next:       nil,
					ParentNode: nil,
				},
				Op:         "",
				Next:       nil,
				ParentNode: nil,
			},
			Op:         "",
			Next:       nil,
			ParentNode: nil,
		},
		Op:         "",
		Next:       nil,
		ParentNode: nil,
	}

	newAST.Visit(newAST, func(parent generic_ast.Expression, e generic_ast.Expression, context generic_ast.VisitorContext) {
		node := e.(generic_ast.TraversableNode)
		node.OverrideParent(parent.(interface{}).(generic_ast.TraversableNode))
	}, generic_ast.NewEmptyVisitorContext())

	expr.LogicalOperation = newAST
}
