package ir

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/parser/utils"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type IRExit struct {
	generic_ast.BaseASTNode
	Type       IRType  `"Exit" @Ident`
	Value      *string `"(" (@Ident)? ")"`
	ParentNode generic_ast.TraversableNode
}

func (ast *IRExit) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *IRExit) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *IRExit) Begin() lexer.Position {
	return ast.Pos
}

func (ast *IRExit) End() lexer.Position {
	return ast.EndPos
}

func (ast *IRExit) HasValue() bool {
	return ast.Value != nil
}

func (ast *IRExit) GetNode() interface{} {
	return ast
}

func (ast *IRExit) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{}
}

func (ast *IRExit) Print(c *context.ParsingContext) string {
	if ast.HasValue() {
		return utils.PrintASTNode(c, ast, "Exit %s (%s)", ast.Type, *ast.Value)
	}
	return utils.PrintASTNode(c, ast, "Exit %s ()", ast.Type)
}

///

func (ast *IRExit) Body() generic_ast.Expression {
	return hindley_milner.Batch{Exp: []generic_ast.Expression{}}
}

func (ast *IRExit) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	if ast.HasValue() {
		return mapper(parent, &IRExit{
			BaseASTNode: ast.BaseASTNode,
			Value:       ast.Value,
			ParentNode:  parent.(generic_ast.TraversableNode),
		}, context, true)
	}
	return mapper(parent, &IRExit{
		BaseASTNode: ast.BaseASTNode,
		ParentNode:  parent.(generic_ast.TraversableNode),
	}, context, true)
}

func (ast *IRExit) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(parent, ast, context)
}

//

func (ast *IRExit) BuildFlowGraph(builder cfg.CFGBuilder) {
	builder.AddBlockSuccesor(ast)
	builder.UpdatePrev([]generic_ast.NormalNode{ast})
	builder.AddBlockSuccesor(builder.Exit())
	builder.UpdatePrev(nil)
}

func (ast *IRExit) GetUsedVariables(vars cfg.VariableSet, visitedMap map[generic_ast.TraversableNode]struct{}) cfg.VariableSet {
	if ast.HasValue() {
		vars.Add(cfg.NewVariable(*ast.Value, nil))
	}
	return vars
}

func (ast *IRExit) RenameVariables(substUsed, substDecl cfg.VariableSubstitution) {
	if ast.HasValue() {
		newVal := substUsed.Replace(*ast.Value)
		ast.Value = &newVal
	}
}
