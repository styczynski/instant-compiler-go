package ast

import (
	"strings"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Statement struct {
	BaseASTNode
	Empty *string `";"`
	BlockStatement *Block `| @@`
	Declaration *Declaration `| @@`
	Assignment *Assignment `| @@`
	UnaryStatement *UnaryStatement `| @@`
	Return *Return `| @@`
	If *If `| @@`
	While *While `| @@`
	Expression *Expression `| @@ ";"`
}

func (ast *Statement) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Statement) End() lexer.Position {
	return ast.EndPos
}

func (ast *Statement) GetNode() interface{} {
	return ast
}

func (ast *Statement) IsEmpty() bool {
	return ast.Empty != nil
}

func (ast *Statement) IsBlockStatement() bool {
	return ast.BlockStatement != nil
}

func (ast *Statement) IsDeclaration() bool {
	return ast.Declaration != nil
}

func (ast *Statement) IsAssignment() bool {
	return ast.Assignment != nil
}

func (ast *Statement) IsUnaryStatement() bool {
	return ast.UnaryStatement != nil
}

func (ast *Statement) IsReturn() bool {
	return ast.Return != nil
}

func (ast *Statement) IsIf() bool {
	return ast.If != nil
}

func (ast *Statement) IsWhile() bool {
	return ast.While != nil
}

func (ast *Statement) IsExpression() bool {
	return ast.Expression != nil
}

func (ast *Statement) GetChildren() []TraversableNode {
	if ast.IsEmpty() {
		return []TraversableNode{ MakeTraversableNodeToken(*ast.Empty, ast.Pos, ast.EndPos) }
	} else if ast.IsBlockStatement() {
		return []TraversableNode{ ast.BlockStatement }
	} else if ast.IsDeclaration() {
		return []TraversableNode{ ast.Declaration }
	} else if ast.IsAssignment() {
		return []TraversableNode{ ast.Assignment }
	} else if ast.IsUnaryStatement() {
		return []TraversableNode{ ast.UnaryStatement }
	} else if ast.IsReturn() {
		return []TraversableNode{ ast.Return }
	} else if ast.IsIf() {
		return []TraversableNode{ ast.If }
	} else if ast.IsWhile() {
		return []TraversableNode{ ast.While }
	} else if ast.IsExpression() {
		return []TraversableNode{ ast.Expression }
	}
	return []TraversableNode{}
}

func (ast *Statement) formatStatementInstruction(statement string, c *context.ParsingContext) string {
	if c.PrinterConfiguration.SkipStatementIdent {
		c.PrinterConfiguration.SkipStatementIdent = false
		return statement
	}
	return printNode(c, ast, "%s%s", strings.Repeat("  ", c.BlockDepth), statement)
}

func (ast *Statement) Print(c *context.ParsingContext) string {
	ret := "UNKNOWN"
	propagateSkipStatementIdent := false
	if ast.IsEmpty() {
		ret = ";"
	} else if ast.IsBlockStatement() {
		if c.PrinterConfiguration.SkipStatementIdent {
			propagateSkipStatementIdent = true
		}
		ret = ast.BlockStatement.Print(c)
	} else if ast.IsDeclaration() {
		ret = ast.Declaration.Print(c)
	} else if ast.IsAssignment() {
		ret = ast.Assignment.Print(c)
	} else if ast.IsUnaryStatement() {
		ret = ast.UnaryStatement.Print(c)
	} else if ast.IsReturn() {
		ret = ast.Return.Print(c)
	} else if ast.IsIf() {
		ret =  ast.If.Print(c)
	} else if ast.IsWhile() {
		ret = ast.While.Print(c)
	} else if ast.IsExpression() {
		ret = printNode(c, ast, "%s;", ast.Expression.Print(c))
	}
	if propagateSkipStatementIdent {
		c.PrinterConfiguration.SkipStatementIdent = true
		ret = ast.formatStatementInstruction(ret, c)
	} else {
		ret = ast.formatStatementInstruction(ret, c)
	}
	return ret
}

//////

func (ast *Statement) Body() hindley_milner.Expression { return ast }
func (ast *Statement) Map(mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	return mapper(ast)
}
func (ast *Statement) Visit(mapper hindley_milner.ExpressionMapper) {
	mapper(ast)
}
func (ast *Statement) Type() hindley_milner.Type {
	return PrimitiveType{
		name:    "void",
	}
}
func (ast *Statement) ExpressionType() hindley_milner.ExpressionType { return hindley_milner.E_LITERAL }