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

type IRIf struct {
	generic_ast.BaseASTNode
	ConditionType IRType `"If" @@`
	Condition     string `@Ident`
	BlockThen     int    `"jump" "to" @Int`
	BlockElse     int    `"else" @Int`
	ParentNode    generic_ast.TraversableNode
	Negated       bool
}

func (ast *IRIf) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *IRIf) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *IRIf) Begin() lexer.Position {
	return ast.Pos
}

func (ast *IRIf) End() lexer.Position {
	return ast.EndPos
}

func (ast *IRIf) GetNode() interface{} {
	return ast
}

func (ast *IRIf) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{
		generic_ast.MakeTraversableNodeValue(ast, ast.Condition, "ident", ast.Pos, ast.EndPos),
	}
}

func (ast *IRIf) HasElseBlock() bool {
	return ast.BlockElse > -1
}

func (ast *IRIf) Print(c *context.ParsingContext) string {
	ifPostfix := " "
	if ast.Negated {
		ifPostfix = " not "
	}
	if !ast.HasElseBlock() {
		return utils.PrintASTNode(c, ast, "If%s%s %s jump to block%d else continue", ifPostfix, ast.ConditionType, ast.Condition, ast.BlockThen)
	}
	return utils.PrintASTNode(c, ast, "If%s%s %s jump to block%d else block%d", ifPostfix, ast.ConditionType, ast.Condition, ast.BlockThen, ast.BlockElse)
}

func (ast *IRIf) GetUsedVariables(vars cfg.VariableSet, visitedMap map[generic_ast.TraversableNode]struct{}) cfg.VariableSet {
	vars.Add(cfg.NewVariable(ast.Condition, nil))
	return vars
}

func (ast *IRIf) RenameVariables(substUsed, substDecl cfg.VariableSubstitution) {
	ast.Condition = substUsed.Replace(ast.Condition)
}

func (ast *IRIf) BuildFlowGraph(builder cfg.CFGBuilder) {
	fmt.Printf("IRIF BUILD FLOW\n")

	builder.AddBlockSuccesor(ast)

	builder.UpdatePrev([]generic_ast.NormalNode{ast})
	builder.BuildNode(ast.ParentNode.(*IRStatement).LookupBlock(ast.BlockThen))

	ctrlExits := builder.GetPrev() // aggregate of builder.prev from each condition

	if ast.HasElseBlock() {
		builder.UpdatePrev([]generic_ast.NormalNode{ast})
		builder.BuildNode(ast.ParentNode.(*IRStatement).LookupBlock(ast.BlockElse))
		ctrlExits = append(ctrlExits, builder.GetPrev()...)
	} else {
		ctrlExits = append(ctrlExits, ast)
	}
	builder.UpdatePrev(ctrlExits)
}

func (ast *IRIf) Body() generic_ast.Expression {
	return hindley_milner.Batch{}
}

func (ast *IRIf) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	return mapper(parent, &IRIf{
		BaseASTNode:   ast.BaseASTNode,
		ConditionType: ast.ConditionType,
		Condition:     ast.Condition,
		BlockThen:     ast.BlockThen,
		BlockElse:     ast.BlockElse,
		Negated:       ast.Negated,
		ParentNode:    parent.(generic_ast.TraversableNode),
	}, context, true)
}

func (ast *IRIf) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(parent, ast, context)
}
