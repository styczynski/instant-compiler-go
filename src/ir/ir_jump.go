package ir

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/parser/utils"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type IRJump struct {
	generic_ast.BaseASTNode
	BlockTarget int `"Jump" "to" @Int`
	ParentNode  generic_ast.TraversableNode
}

func (ast *IRJump) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *IRJump) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *IRJump) Begin() lexer.Position {
	return ast.Pos
}

func (ast *IRJump) End() lexer.Position {
	return ast.EndPos
}

func (ast *IRJump) GetNode() interface{} {
	return ast
}

func (ast *IRJump) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{}
}

func (ast *IRJump) Print(c *context.ParsingContext) string {
	return utils.PrintASTNode(c, ast, "Jump to block%d", ast.BlockTarget)
}

func (ast *IRJump) BuildFlowGraph(builder cfg.CFGBuilder) {
	fmt.Printf("IRIF BUILD FLOW\n")

	builder.AddBlockSuccesor(ast)

	builder.UpdatePrev([]generic_ast.NormalNode{ast})
	builder.BuildNode(ast.ParentNode.(*IRStatement).LookupBlock(ast.BlockTarget))

}

func (ast *IRJump) GetUsedVariables(vars cfg.VariableSet, visitedMap map[generic_ast.TraversableNode]struct{}) cfg.VariableSet {
	return vars
}

func (ast *IRJump) RenameVariables(substUsed, substDecl cfg.VariableSubstitution) {
	// No-op
}

func (ast *IRJump) Body() generic_ast.Expression {
	return hindley_milner.Batch{}
}

func (ast *IRJump) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	return mapper(parent, &IRJump{
		BaseASTNode: ast.BaseASTNode,
		BlockTarget: ast.BlockTarget,
		ParentNode:  parent.(generic_ast.TraversableNode),
	}, context, true)
}

func (ast *IRJump) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(parent, ast, context)
}
