package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
)

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

func (ast *LogicalOperation) Print(c *context.ParsingContext) string {
	if ast.HasNext() {
		return printBinaryOperation(c, ast, ast.Equality.Print(c), ast.Op, ast.Next.Print(c))
	}
	return ast.Equality.Print(c)
}

