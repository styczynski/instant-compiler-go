package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type LogicalOperation struct {
	BaseASTNode
	Equality *Equality `@@`
	Op         string      `[ @( "|" "|" | "&" "&" )`
	Next       *LogicalOperation   `  @@ ]`
	ParentNode TraversableNode
}

func (ast *LogicalOperation) Parent() TraversableNode {
	return ast.ParentNode
}

func (ast *LogicalOperation) OverrideParent(node TraversableNode) {
	ast.ParentNode = node
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
		MakeTraversableNodeToken(ast, ast.Op, ast.Pos, ast.EndPos),
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

/////

func (ast *LogicalOperation) Map(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	next := ast.Next
	if ast.HasNext() {
		next = mapper(ast, ast.Next).(*LogicalOperation)
	}
	return mapper(parent, &LogicalOperation{
		BaseASTNode: ast.BaseASTNode,
		Equality:    mapper(ast, ast.Equality).(*Equality),
		Op:          ast.Op,
		Next:        next,
		ParentNode: parent.(TraversableNode),
	})
}

func (ast *LogicalOperation) Visit(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) {
	mapper(ast, ast.Equality)
	if ast.HasNext() {
		mapper(ast, ast.Next)
	}
	mapper(parent, ast)
}

func (ast *LogicalOperation) Fn() hindley_milner.Expression {
	return &BuiltinFunction{
		BaseASTNode: ast.BaseASTNode,
		name: ast.Op,
	}
}

func (ast *LogicalOperation) Body() hindley_milner.Expression {
	if !ast.HasNext() {
		return ast.Equality
	}
	return hindley_milner.Batch{
		Exp: []hindley_milner.Expression{
			ast.Equality,
			ast.Next,
		},
	}
}

func (ast *LogicalOperation) ExpressionType() hindley_milner.ExpressionType {
	if !ast.HasNext() {
		return hindley_milner.E_PROXY
	}
	return hindley_milner.E_APPLICATION
}