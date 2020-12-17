package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type LogicalOperation struct {
	 generic_ast.BaseASTNode
	Equality *Equality `@@`
	Op         string      `[ @( "|" "|" | "&" "&" )`
	Next       *LogicalOperation   `  @@ ]`
	ParentNode generic_ast.TraversableNode
}

func (ast *LogicalOperation) ExtractConst() (generic_ast.TraversableNode, bool) {
	if ast.HasNext() {
		return nil, false
	}
	return ast.Equality.ExtractConst()
}

func (ast *LogicalOperation) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *LogicalOperation) OverrideParent(node generic_ast.TraversableNode) {
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

func (ast *LogicalOperation) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{
		ast.Equality,
		generic_ast.MakeTraversableNodeToken(ast, ast.Op, ast.Pos, ast.EndPos),
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

func (ast *LogicalOperation) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	next := ast.Next
	if ast.HasNext() {
		next = mapper(ast, ast.Next, context, false).(*LogicalOperation)
	}
	return mapper(parent, &LogicalOperation{
		BaseASTNode: ast.BaseASTNode,
		Equality:    mapper(ast, ast.Equality, context, false).(*Equality),
		Op:          ast.Op,
		Next:        next,
		ParentNode: parent.(generic_ast.TraversableNode),
	}, context, true)
}

func (ast *LogicalOperation) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(ast, ast.Equality, context)
	if ast.HasNext() {
		mapper(ast, ast.Next, context)
	}
	mapper(parent, ast, context)
}

func (ast *LogicalOperation) Fn() generic_ast.Expression {
	return &BuiltinFunction{
		BaseASTNode: ast.BaseASTNode,
		name: ast.Op,
	}
}

func (ast *LogicalOperation) Body() generic_ast.Expression {
	if !ast.HasNext() {
		return ast.Equality
	}
	return hindley_milner.Batch{
		Exp: []generic_ast.Expression{
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