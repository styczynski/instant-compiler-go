package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Return struct {
	generic_ast.BaseASTNode
	Expression *Expression `"return" (@@)? ";"`
	ParentNode generic_ast.TraversableNode
}

func (ast *Return) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *Return) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *Return) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Return) End() lexer.Position {
	return ast.EndPos
}

func (ast *Return) HasExpression() bool {
	return ast.Expression != nil
}

func (ast *Return) GetNode() interface{} {
	return ast
}

func (ast *Return) GetChildren() []generic_ast.TraversableNode {
	if ast.HasExpression() {
		return []generic_ast.TraversableNode{
			ast.Expression,
		}
	}
	return []generic_ast.TraversableNode{}
}

func (ast *Return) Print(c *context.ParsingContext) string {
	if ast.HasExpression() {
		return printNode(c, ast, "return %s;", ast.Expression.Print(c))
	}
	return printNode(c, ast, "return;")
}

///

func (ast *Return) Body() generic_ast.Expression {
	if ast.HasExpression() {
		return ast.Expression
	}
	return hindley_milner.Batch{Exp: []generic_ast.Expression{}}
}

func (ast *Return) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	if ast.HasExpression() {
		return mapper(parent, &Return{
			BaseASTNode: ast.BaseASTNode,
			Expression:  mapper(ast, ast.Expression, context, false).(*Expression),
			ParentNode:  parent.(generic_ast.TraversableNode),
		}, context, true)
	}
	return mapper(parent, &Return{
		BaseASTNode: ast.BaseASTNode,
		ParentNode:  parent.(generic_ast.TraversableNode),
	}, context, true)
}

func (ast *Return) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	if ast.HasExpression() {
		mapper(ast, ast.Expression, context)
	}
	mapper(parent, ast, context)
}

func (ast *Return) ExpressionType() hindley_milner.ExpressionType { return hindley_milner.E_RETURN }

//

func (ast *Return) BuildFlowGraph(builder cfg.CFGBuilder) {
	builder.AddBlockSuccesor(ast)
	builder.UpdatePrev([]generic_ast.NormalNode{ast})
	builder.AddBlockSuccesor(builder.Exit())
	builder.UpdatePrev(nil)
}
