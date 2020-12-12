package ast

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Typename struct {
	BaseASTNode
	Expr      *LogicalOperation       `"typename" @@`
	ParentNode TraversableNode
}

func (ast *Typename) Parent() TraversableNode {
	return ast.ParentNode
}

func (ast *Typename) OverrideParent(node TraversableNode) {
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

func (ast *Typename) GetTraversableNode() TraversableNode {
	return ast.Expr
}

func (ast *Typename) Print(c *context.ParsingContext) string {
	return fmt.Sprintf("typename %s", ast.Expr.Print(c))
}

func (ast *Typename) GetChildren() []TraversableNode {
	return []TraversableNode{
		ast.Expr,
	}
}

////

func (ast *Typename) Map(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	// TODO
	return ast
}

func (ast *Typename) Visit(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) {
	// TODO
	mapper(ast, ast.Expr)
	mapper(parent, ast)
}

func (ast *Typename) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_INTROSPECTION
}

func (ast *Typename) Body() hindley_milner.Expression {
	return ast.Expr
}

func (ast *Typename) OnTypeReturned(t hindley_milner.Type) {
	parent := ast.Parent()
	expr := (parent.(interface{}).(*Expression))
	expr.LogicalOperation = nil
	expr.Typename = nil
	expr.NewType = nil

	val := fmt.Sprintf("%v", t)
	//fmt.Printf("TYPENAME: %s\n", val)

	newAST := &LogicalOperation{
		BaseASTNode: ast.BaseASTNode,
		Equality:    &Equality{
			BaseASTNode: ast.BaseASTNode,
			Comparison:  &Comparison{
				BaseASTNode: ast.BaseASTNode,
				Addition:    &Addition{
					BaseASTNode:    ast.BaseASTNode,
					Multiplication: &Multiplication{
						BaseASTNode: ast.BaseASTNode,
						Unary:      &Unary{
							BaseASTNode:      ast.BaseASTNode,
							Op:               "",
							Unary:            &Unary{
								BaseASTNode:      ast.BaseASTNode,
								Op:               "",
								Unary:            nil,
								UnaryApplication: &UnaryApplication{
									BaseASTNode: ast.BaseASTNode,
									Target:      nil,
									Arguments:   nil,
									Index:       &Index{
										BaseASTNode:  ast.BaseASTNode,
										Primary:      &Primary{
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
									ParentNode:  nil,
								},
								ParentNode:       nil,
							},
							UnaryApplication: nil,
							ParentNode:       nil,
						},
						Op:          "",
						Next:        nil,
						ParentNode:  nil,
					},
					Op:             "",
					Next:           nil,
					ParentNode:     nil,
				},
				Op:          "",
				Next:        nil,
				ParentNode:  nil,
			},
			Op:          "",
			Next:        nil,
			ParentNode:  nil,
		},
		Op:          "",
		Next:        nil,
		ParentNode:  nil,
	}

	newAST.Visit(newAST, func(parent hindley_milner.Expression, e hindley_milner.Expression) hindley_milner.Expression {
		node := e.(TraversableNode)
		node.OverrideParent(parent.(interface{}).(TraversableNode))
		return e
	})

	expr.LogicalOperation = newAST
}