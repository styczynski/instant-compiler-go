package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
)

type If struct {
	BaseASTNode
	Condition *Expression `"if" "(" @@ ")"`
	Then *Statement `@@`
	Else *Statement `( "else" @@ )?`
}

func (ast *If) Begin() lexer.Position {
	return ast.Pos
}

func (ast *If) End() lexer.Position {
	return ast.EndPos
}

func (ast *If) GetNode() interface{} {
	return ast
}

func (ast *If) GetChildren() []TraversableNode {
	return []TraversableNode{
		ast.Condition,
		ast.Then,
		ast.Else,
	}
}

func (ast *If) HasElseBlock() bool {
	return ast.Else != nil
}

func (ast *If) Print(c *context.ParsingContext) string {
	if ast.HasElseBlock(){
		return printNode(c, ast, "if (%s) %s else %s", ast.Condition.Print(c), makeBlockFromStatement(ast.Then).Print(c), makeBlockFromStatement(ast.Else).Print(c))
	}
	return printNode(c, ast, "if (%s) %s", ast.Condition.Print(c), ast.Then.Print(c))
}
