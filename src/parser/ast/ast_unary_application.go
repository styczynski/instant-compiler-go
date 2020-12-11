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
	Index *Index `| @@`
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
	} else if ast.IsIndex() {
		return []TraversableNode{
			ast.Index,
		}
	}
	return []TraversableNode{}
}

func (ast *UnaryApplication) IsApplication() bool {
	return ast.Target != nil
}

func (ast *UnaryApplication) IsIndex() bool {
	return ast.Index != nil
}

func (ast *UnaryApplication) Print(c *context.ParsingContext) string {
	if ast.IsApplication() {
		args := []string{}
		for _, argument := range ast.Arguments {
			args = append(args, argument.Print(c))
		}
		return printNode(c, ast, "%s(%s)", *ast.Target, strings.Join(args, ", "))
	} else if ast.IsIndex() {
		return ast.Index.Print(c)
	}
	return "UNKNOWN"
}

////

func (ast *UnaryApplication) Map(mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	if ast.IsIndex() {
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
			Index: mapper(ast.Index).(*Index),
		})
	}
	panic("Invalid UnaryApplication operation type")
}

func (ast *UnaryApplication) Visit(mapper hindley_milner.ExpressionMapper) {
	if ast.IsIndex() {
		mapper(ast.Index)
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
	if ast.IsIndex() {
		return ast.Index
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
	if ast.IsIndex() {
		return hindley_milner.E_PROXY
	}
	return hindley_milner.E_APPLICATION
}


