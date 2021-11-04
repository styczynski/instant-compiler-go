package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type PrimaryInvalid struct {
	Reason string
	Source generic_ast.TraversableNode
}

type Primary struct {
	generic_ast.BaseASTNode
	Variable      *string     `@Ident`
	Int           *int64      `| @Int`
	ParentNode    generic_ast.TraversableNode
	Invalid       *PrimaryInvalid
}

func (ast *Primary) ExtractConst() (generic_ast.TraversableNode, bool) {
	if ast.IsInvalid() {
		return ast, false
	} else if ast.IsVariable() {
		return nil, false
	}
	return ast, true
}

// Basic arithmetic functions

func (a *Primary) IsInvalid() bool {
	return a.Invalid != nil
}

func (a *Primary) GetInvalidReason() PrimaryInvalid {
	return *a.Invalid
}

func (a *Primary) Add(b *Primary, Op string) *Primary {
	if a.IsInt() && b.IsInt() {
		v := int64(0)
		if Op == "+" {
			v = *a.Int + *b.Int
		} else if Op == "-" {
			v = *a.Int - *b.Int
		} else {
			panic("Invalid add operation")
		}
		return &Primary{
			BaseASTNode: a.BaseASTNode,
			Int:         &v,
		}
	}
	panic("Invalid addition")
}

func (a *Primary) Mul(b *Primary, Op string) *Primary {
	if a.IsInt() && b.IsInt() {
		v := int64(0)
		if Op == "/" {
			if *b.Int == 0 {
				return &Primary{
					BaseASTNode: a.BaseASTNode,
					Invalid: &PrimaryInvalid{
						Reason: "Detected division by 0.",
						Source: b,
					},
					Int: &v,
				}
			} else {
				v = *a.Int / *b.Int
			}
		} else if Op == "*" {
			v = *a.Int * *b.Int
		} else if Op == "%" {
			if *b.Int == 0 {
				return &Primary{
					BaseASTNode: a.BaseASTNode,
					Invalid: &PrimaryInvalid{
						Reason: "Detected modulo operation with 0 base.",
						Source: b,
					},
					Int: &v,
				}
			} else {
				v = *a.Int % *b.Int
			}
		} else {
			panic("Invalid multiplication type")
		}
		return &Primary{
			BaseASTNode: a.BaseASTNode,
			Int:         &v,
		}
	}
	panic("Invalid addition")
}

func (ast *Primary) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *Primary) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *Primary) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Primary) End() lexer.Position {
	return ast.EndPos
}

func (ast *Primary) GetNode() interface{} {
	return ast
}

func (ast *Primary) GetChildren() []generic_ast.TraversableNode {
	if ast.IsVariable() {
		return []generic_ast.TraversableNode{
			generic_ast.MakeTraversableNodeValue(ast, *ast.Variable, "ident", ast.Pos, ast.EndPos),
		}
	} else if ast.IsInt() {
		return []generic_ast.TraversableNode{
			generic_ast.MakeTraversableNodeValue(ast, *ast.Int, "int", ast.Pos, ast.EndPos),
		}
	}
	return []generic_ast.TraversableNode{}
}

func (ast *Primary) IsVariable() bool {
	return ast.Variable != nil
}

func (ast *Primary) IsInt() bool {
	return ast.Int != nil
}

func (ast *Primary) Print(c *context.ParsingContext) string {
	if ast.IsVariable() {
		return printNode(c, ast, "%s", *ast.Variable)
	} else if ast.IsInt() {
		return printNode(c, ast, "%d", *ast.Int)
	}
	panic("Unvalid Expression value")
	return "UNKNOWN"
}

////

func (ast *Primary) Name() hindley_milner.NameGroup {
	if ast.IsVariable() {
		return hindley_milner.Name(*ast.Variable)
	}
	panic("Cannot get name for Primary expression which is not a variable")
}
func (ast *Primary) Body() generic_ast.Expression {
	return ast
}
func (ast *Primary) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	return mapper(parent, &Primary{
		BaseASTNode:   ast.BaseASTNode,
		Variable:      ast.Variable,
		Int:           ast.Int,
		ParentNode:    ast.ParentNode,
	}, context, true)
}
func (ast *Primary) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(parent, ast, context)
}
func (ast *Primary) Type() hindley_milner.Type {
	if ast.IsVariable() {
		return nil
	} else if ast.IsInt() {
		return CreatePrimitive(T_INT)
	}
	panic("Unknown Primary type")
}

func (ast *Primary) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_LITERAL
}
