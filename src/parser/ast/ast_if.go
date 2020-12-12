package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type If struct {
	BaseASTNode
	Condition *Expression `"if" "(" @@ ")"`
	Then *Statement `@@`
	Else *Statement `( "else" @@ )?`
	ParentNode TraversableNode
}

func (ast *If) Parent() TraversableNode {
	return ast.ParentNode
}

func (ast *If) OverrideParent(node TraversableNode) {
	ast.ParentNode = node
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

///

func (ast *If) Map(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	return mapper(parent, &If{
		BaseASTNode: ast.BaseASTNode,
		Condition: mapper(ast, ast.Condition).(*Expression),
		Then: mapper(ast, ast.Then).(*Statement),
		Else: mapper(ast, ast.Else).(*Statement),
		ParentNode: parent.(TraversableNode),
	})
}

func (ast *If) Visit(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) {
	mapper(ast, ast.Condition)
	mapper(ast, ast.Then)
	if ast.HasElseBlock() {
		mapper(ast, ast.Else)
	}
	mapper(parent, ast)
}

func (ast *If) Fn() hindley_milner.Expression {
	return &BuiltinFunction{
		BaseASTNode: ast.BaseASTNode,
		name: "if",
	}
}

func (ast *If) Body() hindley_milner.Expression {
	args := []hindley_milner.Expression{
		ast.Condition,
		ast.Then,
	}
	if ast.HasElseBlock() {
		args = append(args, ast.Else)
	} else {
		args = append(args, ast.Then)
	}
	return hindley_milner.Batch{
		Exp: args,
	}
}

func (ast *If) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_APPLICATION
}