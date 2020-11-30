package parser

import (
	"fmt"
	"strings"
)

type Expression struct {
	LogicalOperation *LogicalOperation `@@`
}

func printBinaryOperation(c *ParsingContext, arg1 string, operator string, arg2 string) string{
	return fmt.Sprintf("%s %s %s", arg1, operator, arg2)
}

func printUnaryOperation(c *ParsingContext, operator string, arg string) string{
	return fmt.Sprintf("%s%s", operator, arg)
}

func (ast *Expression) Print(c *ParsingContext) string {
	return ast.LogicalOperation.Print(c)
}

type LogicalOperation struct {
	Equality *Equality `@@`
	Op         string      `[ @( "|" "|" | "&" "&" )`
	Next       *LogicalOperation   `  @@ ]`
}

func (ast *LogicalOperation) HasNext() bool {
	return ast.Next != nil
}

func (ast *LogicalOperation) Print(c *ParsingContext) string {
	if ast.HasNext() {
		return printBinaryOperation(c, ast.Equality.Print(c), ast.Op, ast.Next.Print(c))
	}
	return ast.Equality.Print(c)
}

type Equality struct {
	Comparison *Comparison `@@`
	Op         string      `[ @( "!" "=" | "=" "=" )`
	Next       *Equality   `  @@ ]`
}

func (ast *Equality) HasNext() bool {
	return ast.Next != nil
}

func (ast *Equality) Print(c *ParsingContext) string {
	if ast.HasNext() {
		return printBinaryOperation(c, ast.Comparison.Print(c), ast.Op, ast.Next.Print(c))
	}
	return ast.Comparison.Print(c)
}

type Comparison struct {
	Addition *Addition   `@@`
	Op       string      `[ @( ">" | ">" "=" | "<" | "<" "=" )`
	Next     *Comparison `  @@ ]`
}

func (ast *Comparison) HasNext() bool {
	return ast.Next != nil
}

func (ast *Comparison) Print(c *ParsingContext) string {
	if ast.HasNext() {
		return printBinaryOperation(c, ast.Addition.Print(c), ast.Op, ast.Next.Print(c))
	}
	return ast.Addition.Print(c)
}

type Addition struct {
	Multiplication *Multiplication `@@`
	Op             string          `[ @( "-" | "+" )`
	Next           *Addition       `  @@ ]`
}

func (ast *Addition) HasNext() bool {
	return ast.Next != nil
}

func (ast *Addition) Print(c *ParsingContext) string {
	if ast.HasNext() {
		return printBinaryOperation(c, ast.Multiplication.Print(c), ast.Op, ast.Next.Print(c))
	}
	return ast.Multiplication.Print(c)
}

type Multiplication struct {
	Unary *Unary          `@@`
	Op    string          `[ @( "/" | "*" )`
	Next  *Multiplication `  @@ ]`
}

func (ast *Multiplication) HasNext() bool {
	return ast.Next != nil
}

func (ast *Multiplication) Print(c *ParsingContext) string {
	if ast.HasNext() {
		return printBinaryOperation(c, ast.Unary.Print(c), ast.Op, ast.Next.Print(c))
	}
	return ast.Unary.Print(c)
}

type Unary struct {
	Op      string   `  ( @( "!" | "-" )`
	Unary   *Unary   `    @@ )`
	UnaryApplication *UnaryApplication `| @@`
}

func (ast *Unary) IsOperation() bool {
	return ast.Unary != nil
}

func (ast *Unary) IsUnaryApplication() bool {
	return ast.UnaryApplication != nil
}

func (ast *Unary) Print(c *ParsingContext) string {
	if ast.IsOperation() {
		return printUnaryOperation(c, ast.Op, ast.Unary.Print(c))
	} else if ast.IsUnaryApplication() {
		return ast.UnaryApplication.Print(c)
	}
	return "UNKNOWN"
}

type UnaryApplication struct {
	Target *string   `( @Ident`
	Arguments []*Expression   `"(" (@@ ("," @@)*)? ")" )`
	Primary *Primary `| @@`
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
		return fmt.Sprintf("%s(%s)", *ast.Target, strings.Join(args, ", "))
	} else if ast.IsPrimary() {
		return ast.Primary.Print(c)
	}
	return "UNKNOWN"
}


type Primary struct {
	Variable   *string   `@Ident`
	Int        *int64    `| @Int`
	String        *string     `| @String`
	Bool          *bool       `| ( @"true" | "false" )`
	SubExpression *Expression `| "(" @@ ")" `
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
		return fmt.Sprintf("%s", *ast.Variable)
	} else if ast.IsInt() {
		return fmt.Sprintf("%d", *ast.Int)
	} else if ast.IsString() {
		return fmt.Sprintf("\"%s\"", *ast.String)
	} else if ast.IsBool() {
		return fmt.Sprintf("%b", *ast.Bool)
	} else if ast.IsSubexpression() {
		return fmt.Sprintf("(%s)", ast.SubExpression.Print(c))
	}
	return "UNKNOWN"
}