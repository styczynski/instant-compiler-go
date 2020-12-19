package ast

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
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
	Variable   *string   `@Ident`
	Int        *int64    `| @Int`
	String        *string     `| @String`
	Bool          *bool       `| @( "true" | "false" )`
	SubExpression *Expression `| ( "(" @@ ")" )`
	ParentNode generic_ast.TraversableNode
	Invalid *PrimaryInvalid
}

func (ast *Primary) ExtractConst() (generic_ast.TraversableNode, bool) {
	if ast.IsInvalid() {
		return ast, false
	} else if ast.IsSubexpression() {
		return ast.SubExpression.ExtractConst()
	} else if ast.IsVariable() {
		return nil, false
	}
	return ast, true
}

// Basic arithmetic functions

func (a *Primary) IsInvalid() bool {
	return a.Invalid != nil
}

func (a *Primary) ValidateConstFold() (error, generic_ast.TraversableNode) {
	if a.IsInvalid() {
		return fmt.Errorf(a.Invalid.Reason), a.Invalid.Source
	}
	return nil, nil
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
			BaseASTNode:   a.BaseASTNode,
			Int: &v,
		}
	} else if a.IsString() && b.IsString() {
		v := ""
		if Op == "+" {
			v = *a.String + *b.String
		} else {
			panic("Invalid add operation")
		}
		return &Primary{
			BaseASTNode:   a.BaseASTNode,
			String: &v,
		}
	}
	panic("Invalid addition")
}


func (a *Primary) Compare(b *Primary, Op string) *Primary {
	if a.IsString() && b.IsString() {
		v := false
		if Op == ">=" {
			v = *a.String >= *b.String
		} else if Op == "<=" {
			v = *a.String <= *b.String
		} else if Op == "==" {
			v = *a.String == *b.String
		} else if Op == "!=" {
			v = *a.String != *b.String
		} else if Op == "<" {
			v = *a.String < *b.String
		} else if Op == ">" {
			v = *a.String > *b.String
		} else {
			panic("Invalid comaprison type")
		}
		return &Primary{
			BaseASTNode:   a.BaseASTNode,
			Bool: &v,
		}
	} else if a.IsBool() && b.IsBool() {
		v := false
		if Op == "==" {
			v = *a.Bool == *b.Bool
		}
		return &Primary{
			BaseASTNode:   a.BaseASTNode,
			Bool: &v,
		}
	} else if a.IsInt() && b.IsInt() {
		v := false
		if Op == ">=" {
			v = *a.Int >= *b.Int
		} else if Op == "<=" {
			v = *a.Int <= *b.Int
		} else if Op == "==" {
			v = *a.Int == *b.Int
		} else if Op == "!=" {
			v = *a.Int != *b.Int
		} else if Op == "<" {
			v = *a.Int < *b.Int
		} else if Op == ">" {
			v = *a.Int > *b.Int
		} else {
			panic("Invalid comaprison type")
		}
		return &Primary{
			BaseASTNode:   a.BaseASTNode,
			Bool: &v,
		}
	}
	panic("Invalid addition")
}

func (a *Primary) And(b *Primary, Op string) *Primary {
	if a.IsBool() && b.IsBool() {
		v := false
		if Op == "&&" {
			v = *a.Bool && *b.Bool
		} else if Op == "||" {
			v = *a.Bool || *b.Bool
		} else {
			panic("Invalid and operator")
		}
		return &Primary{
			BaseASTNode:   a.BaseASTNode,
			Bool: &v,
		}
	}
	panic("Invalid and")
}

func (a *Primary) Mul(b *Primary, Op string) *Primary {
	if a.IsInt() && b.IsInt() {
		v := int64(0)
		if Op == "/" {
			if *b.Int == 0 {
				return &Primary{
					BaseASTNode:   a.BaseASTNode,
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
					BaseASTNode:   a.BaseASTNode,
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
			BaseASTNode:   a.BaseASTNode,
			Int: &v,
		}
	}
	panic("Invalid addition")
}

func (ast *Primary) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *Primary) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node

	// Normalize
	if ast.IsVariable() {
		if *ast.Variable == "true" {
			ast.Variable = nil
			v := true
			ast.Bool = &v
		} else if *ast.Variable == "false" {
			ast.Variable = nil
			v := false
			ast.Bool = &v
		}
	}
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
	} else if ast.IsString() {
		return []generic_ast.TraversableNode{
			generic_ast.MakeTraversableNodeValue(ast, *ast.String, "string", ast.Pos, ast.EndPos),
		}
	} else if ast.IsBool() {
		return []generic_ast.TraversableNode{
			generic_ast.MakeTraversableNodeValue(ast, *ast.Bool, "bool", ast.Pos, ast.EndPos),
		}
	} else if ast.IsSubexpression() {
		return []generic_ast.TraversableNode{
			ast.SubExpression,
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

func (ast *Primary) IsString() bool {
	return ast.String != nil
}

func (ast *Primary) IsBool() bool {
	return ast.Bool != nil
}

func (ast *Primary) IsSubexpression() bool {
	return ast.SubExpression != nil
}

func (ast *Primary) Print(c *context.ParsingContext) string {
	if ast.IsVariable() {
		return printNode(c, ast, "%s", *ast.Variable)
	} else if ast.IsInt() {
		return printNode(c, ast, "%d", *ast.Int)
	} else if ast.IsString() {
		return printNode(c, ast, "\"%s\"", *ast.String)
	} else if ast.IsBool() {
		return printNode(c, ast, "%v", *ast.Bool)
	} else if ast.IsSubexpression() {
		return printNode(c, ast, "(%s)", ast.SubExpression.Print(c))
	}
	panic("Unvalid Expression value")
	return "UNKNOWN"
}

////

func (ast *Primary) Name() hindley_milner.NameGroup     {
	if ast.IsVariable() {
		return hindley_milner.Name(*ast.Variable)
	}
	panic("Cannot get name for Primary expression which is not a variable")
}
func (ast *Primary) Body() generic_ast.Expression {
	if ast.IsSubexpression() {
		return ast.SubExpression
	}
	return ast
}
func (ast *Primary) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	if ast.IsSubexpression() {
		return mapper(parent, &Primary{
			BaseASTNode:   ast.BaseASTNode,
			SubExpression: mapper(ast, ast.SubExpression, context, false).(*Expression),
			ParentNode:    ast.ParentNode,
		}, context, true)
	}
	return mapper(parent, &Primary{
		BaseASTNode:   ast.BaseASTNode,
		Variable:      ast.Variable,
		Int:           ast.Int,
		String:        ast.String,
		Bool:          ast.Bool,
		SubExpression: nil,
		ParentNode:    ast.ParentNode,
	}, context, true)
}
func (ast *Primary) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	if ast.IsSubexpression() {
		mapper(ast, ast.SubExpression, context)
	}
	mapper(parent, ast, context)
}
func (ast *Primary) Type() hindley_milner.Type {
	if ast.IsVariable() {
		return nil
	} else if ast.IsInt() {
		return CreatePrimitive(T_INT)
	} else if ast.IsString() {
		return CreatePrimitive(T_STRING)
	} else if ast.IsBool() {
		return CreatePrimitive(T_BOOL)
	} else if ast.IsSubexpression() {
		return nil
	}
	panic("Unknown Primary type")
}

func  (ast *Primary)  ExpressionType() hindley_milner.ExpressionType {
	if ast.IsSubexpression() {
		return hindley_milner.E_PROXY
	}
	return hindley_milner.E_LITERAL
}

///

func (ast *Primary) GetUsedVariables(vars cfg.VariableSet) cfg.VariableSet {
	if ast.IsVariable() {
		return cfg.NewVariableSet(cfg.NewVariable(*ast.Variable, ast))
	} else if ast.IsSubexpression() {
		return vars
	}
	return cfg.NewVariableSet()
}

func (ast *Primary) RenameVariables(subst cfg.VariableSubstitution) {
	if ast.IsVariable() {
		v := subst.Replace(*ast.Variable)
		ast.Variable = &v
	}
}

