package parser

import (
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
)

type Expression struct {
	ComplexASTNode
	LogicalOperation *LogicalOperation `@@`
}

func (ast *Expression) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Expression) End() lexer.Position {
	return ast.EndPos
}

func (ast *Expression) GetNode() interface{} {
	return ast
}

func (ast *Expression) GetChildren() []TraversableNode {
	return []TraversableNode{ ast.LogicalOperation, }
}

func printBinaryOperation(c *ParsingContext, ast TraversableNode, arg1 string, operator string, arg2 string) string{
	return printNode(c, ast, "%s %s %s", arg1, operator, arg2)
}

func printUnaryOperation(c *ParsingContext, ast TraversableNode, operator string, arg string) string{
	return printNode(c, ast, "%s%s", operator, arg)
}

func (ast *Expression) Print(c *ParsingContext) string {
	return ast.LogicalOperation.Print(c)
}

type LogicalOperation struct {
	BaseASTNode
	Equality *Equality `@@`
	Op         string      `[ @( "|" "|" | "&" "&" )`
	Next       *LogicalOperation   `  @@ ]`
}

func (ast *LogicalOperation) Begin() lexer.Position {
	return ast.Pos
}

func (ast *LogicalOperation) End() lexer.Position {
	return ast.EndPos
}

func (ast *LogicalOperation) GetNode() interface{} {
	return ast
}

func (ast *LogicalOperation) GetChildren() []TraversableNode {
	return []TraversableNode{
		ast.Equality,
		MakeTraversableNodeToken(ast.Op, ast.Pos, ast.EndPos),
		ast.Next,
	}
}

func (ast *LogicalOperation) HasNext() bool {
	return ast.Next != nil
}

func (ast *LogicalOperation) Print(c *ParsingContext) string {
	if ast.HasNext() {
		return printBinaryOperation(c, ast, ast.Equality.Print(c), ast.Op, ast.Next.Print(c))
	}
	return ast.Equality.Print(c)
}

type Equality struct {
	BaseASTNode
	Comparison *Comparison `@@`
	Op         string      `[ @( "!" "=" | "=" "=" )`
	Next       *Equality   `  @@ ]`
}

func (ast *Equality) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Equality) End() lexer.Position {
	return ast.EndPos
}

func (ast *Equality) GetNode() interface{} {
	return ast
}

func (ast *Equality) GetChildren() []TraversableNode {
	return []TraversableNode{
		ast.Comparison,
		MakeTraversableNodeToken(ast.Op, ast.Pos, ast.EndPos),
		ast.Next,
	}
}

func (ast *Equality) HasNext() bool {
	return ast.Next != nil
}

func (ast *Equality) Print(c *ParsingContext) string {
	if ast.HasNext() {
		return printBinaryOperation(c, ast, ast.Comparison.Print(c), ast.Op, ast.Next.Print(c))
	}
	return ast.Comparison.Print(c)
}

type Comparison struct {
	BaseASTNode
	Addition *Addition   `@@`
	Op       string      `[ @( ">" | ">" "=" | "<" | "<" "=" )`
	Next     *Comparison `  @@ ]`
}

func (ast *Comparison) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Comparison) End() lexer.Position {
	return ast.EndPos
}

func (ast *Comparison) GetNode() interface{} {
	return ast
}

func (ast *Comparison) GetChildren() []TraversableNode {
	return []TraversableNode{
		ast.Addition,
		MakeTraversableNodeToken(ast.Op, ast.Pos, ast.EndPos),
		ast.Next,
	}
}

func (ast *Comparison) HasNext() bool {
	return ast.Next != nil
}

func (ast *Comparison) Print(c *ParsingContext) string {
	if ast.HasNext() {
		return printBinaryOperation(c, ast, ast.Addition.Print(c), ast.Op, ast.Next.Print(c))
	}
	return ast.Addition.Print(c)
}

type Addition struct {
	BaseASTNode
	Multiplication *Multiplication `@@`
	Op             string          `[ @( "-" | "+" )`
	Next           *Addition       `  @@ ]`
}

func (ast *Addition) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Addition) End() lexer.Position {
	return ast.EndPos
}

func (ast *Addition) GetNode() interface{} {
	return ast
}

func (ast *Addition) GetChildren() []TraversableNode {
	return []TraversableNode{
		ast.Multiplication,
		MakeTraversableNodeToken(ast.Op, ast.Pos, ast.EndPos),
		ast.Next,
	}
}

func (ast *Addition) HasNext() bool {
	return ast.Next != nil
}

func (ast *Addition) Print(c *ParsingContext) string {
	if ast.HasNext() {
		return printBinaryOperation(c, ast, ast.Multiplication.Print(c), ast.Op, ast.Next.Print(c))
	}
	return ast.Multiplication.Print(c)
}

type Multiplication struct {
	BaseASTNode
	Unary *Unary          `@@`
	Op    string          `[ @( "/" | "*" )`
	Next  *Multiplication `  @@ ]`
}

func (ast *Multiplication) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Multiplication) End() lexer.Position {
	return ast.EndPos
}

func (ast *Multiplication) GetNode() interface{} {
	return ast
}

func (ast *Multiplication) GetChildren() []TraversableNode {
	return []TraversableNode{
		ast.Unary,
		MakeTraversableNodeToken(ast.Op, ast.Pos, ast.EndPos),
		ast.Next,
	}
}

func (ast *Multiplication) HasNext() bool {
	return ast.Next != nil
}

func (ast *Multiplication) Print(c *ParsingContext) string {
	if ast.HasNext() {
		return printBinaryOperation(c, ast, ast.Unary.Print(c), ast.Op, ast.Next.Print(c))
	}
	return ast.Unary.Print(c)
}

type Unary struct {
	BaseASTNode
	Op      string   `  ( @( "!" | "-" )`
	Unary   *Unary   `    @@ )`
	UnaryApplication *UnaryApplication `| @@`
}

func (ast *Unary) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Unary) End() lexer.Position {
	return ast.EndPos
}

func (ast *Unary) GetNode() interface{} {
	return ast
}

func (ast *Unary) GetChildren() []TraversableNode {
	if ast.IsOperation() {
		return []TraversableNode{
			MakeTraversableNodeToken(ast.Op, ast.Pos, ast.EndPos),
			ast.Unary,
		}
	} else if ast.IsUnaryApplication() {
		return []TraversableNode{
			ast.UnaryApplication,
		}
	}
	return []TraversableNode{}
}

func (ast *Unary) IsOperation() bool {
	return ast.Unary != nil
}

func (ast *Unary) IsUnaryApplication() bool {
	return ast.UnaryApplication != nil
}

func (ast *Unary) Print(c *ParsingContext) string {
	if ast.IsOperation() {
		return printUnaryOperation(c, ast, ast.Op, ast.Unary.Print(c))
	} else if ast.IsUnaryApplication() {
		return ast.UnaryApplication.Print(c)
	}
	return "UNKNOWN"
}

type UnaryApplication struct {
	BaseASTNode
	Target *string   `( @Ident`
	Arguments []*Expression   `"(" (@@ ("," @@)*)? ")" )`
	Primary *Primary `| @@`
}

func (ast *UnaryApplication) Begin() lexer.Position {
	return ast.Pos
}

func (ast *UnaryApplication) End() lexer.Position {
	return ast.EndPos
}

func (ast *UnaryApplication) GetNode() interface{} {
	return ast
}

func (ast *UnaryApplication) GetChildren() []TraversableNode {
	if ast.IsApplication() {
		nodes := make([]TraversableNode, len(ast.Arguments) + 1)
		nodes = append(nodes, MakeTraversableNodeToken(*ast.Target, ast.Pos, ast.EndPos))
		for _, child := range ast.Arguments {
			nodes = append(nodes, child)
		}
		return nodes
	} else if ast.IsPrimary() {
		return []TraversableNode{
			ast.Primary,
		}
	}
	return []TraversableNode{}
}

func (ast *UnaryApplication) IsApplication() bool {
	return ast.Target != nil
}

func (ast *UnaryApplication) IsPrimary() bool {
	return ast.Primary != nil
}

func (ast *UnaryApplication) Print(c *ParsingContext) string {
	if ast.IsApplication() {
		args := []string{}
		for _, argument := range ast.Arguments {
			args = append(args, argument.Print(c))
		}
		return printNode(c, ast, "%s(%s)", *ast.Target, strings.Join(args, ", "))
	} else if ast.IsPrimary() {
		return ast.Primary.Print(c)
	}
	return "UNKNOWN"
}

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

func (ast *Primary) Print(c *ParsingContext) string {
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