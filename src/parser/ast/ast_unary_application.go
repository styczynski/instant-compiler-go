package ast

import (
	"strings"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
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

////

func (ast *UnaryApplication) Map(mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	if ast.IsPrimary() {
		args := []*Expression{}
		for _, arg := range ast.Arguments {
			args = append(args, mapper(arg).(*Expression))
		}
		return mapper(&UnaryApplication{
			BaseASTNode: ast.BaseASTNode,
			Target:      ast.Target,
			Arguments:   args,
		})
	} else if ast.IsApplication() {
		return mapper(&UnaryApplication{
			BaseASTNode: ast.BaseASTNode,
			Primary: mapper(ast.Primary).(*Primary),
		})
	}
	panic("Invalid UnaryApplication operation type")
}

func (ast *UnaryApplication) Visit(mapper hindley_milner.ExpressionMapper) {
	if ast.IsPrimary() {
		mapper(ast.Primary)
	} else if ast.IsApplication() {
		for _, arg := range ast.Arguments {
			mapper(arg)
		}
	}
	mapper(ast)
}

func (ast *UnaryApplication) Fn() hindley_milner.Expression {
	return &VarName{
		BaseASTNode: ast.BaseASTNode,
		name: *ast.Target,
	}
}

func (ast *UnaryApplication) Body() hindley_milner.Expression {
	if ast.IsPrimary() {
		return ast.Primary
	}
	args := []hindley_milner.Expression{}
	for _, arg := range ast.Arguments {
		args = append(args, arg)
	}
	return hindley_milner.Batch{
		Exp: args,
	}
}

func (ast *UnaryApplication) ExpressionType() hindley_milner.ExpressionType {
	if ast.IsPrimary() {
		return hindley_milner.E_PROXY
	}
	return hindley_milner.E_APPLICATION
}


