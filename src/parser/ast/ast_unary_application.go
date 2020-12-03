package ast

import (
	"strings"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
)

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

func (ast *UnaryApplication) Print(c *context.ParsingContext) string {
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
