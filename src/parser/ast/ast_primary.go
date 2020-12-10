package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)


type Primary struct {
	BaseASTNode
	Variable   *string   `@Ident`
	Int        *int64    `| @Int`
	String        *string     `| @String`
	Bool          *bool       `| ( @"true" | "false" )`
	SubExpression *Expression `| "(" @@ ")" `
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

func (ast *Primary) GetChildren() []TraversableNode {
	if ast.IsVariable() {
		return []TraversableNode{
			MakeTraversableNodeValue(*ast.Variable, "ident", ast.Pos, ast.EndPos),
		}
	} else if ast.IsInt() {
		return []TraversableNode{
			MakeTraversableNodeValue(*ast.Int, "int", ast.Pos, ast.EndPos),
		}
	} else if ast.IsString() {
		return []TraversableNode{
			MakeTraversableNodeValue(*ast.String, "string", ast.Pos, ast.EndPos),
		}
	} else if ast.IsBool() {
		return []TraversableNode{
			MakeTraversableNodeValue(*ast.Bool, "bool", ast.Pos, ast.EndPos),
		}
	} else if ast.IsSubexpression() {
		return []TraversableNode{
			ast.SubExpression,
		}
	}
	return []TraversableNode{}
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
		return printNode(c, ast, "%b", *ast.Bool)
	} else if ast.IsSubexpression() {
		return printNode(c, ast, "(%s)", ast.SubExpression.Print(c))
	}
	return "UNKNOWN"
}

////

func (ast *Primary) Name() hindley_milner.NameGroup     {
	if ast.IsVariable() {
		return hindley_milner.Name(*ast.Variable)
	}
	panic("Cannot get name for Primary expression which is not a variable")
}
func (ast *Primary) Body() hindley_milner.Expression {
	if ast.IsSubexpression() {
		return ast.SubExpression
	}
	return ast
}
func (ast *Primary) Map(mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	// TODO
	return ast
}
func (ast *Primary) Visit(mapper hindley_milner.ExpressionMapper) {
	// TODO
	mapper(ast)
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