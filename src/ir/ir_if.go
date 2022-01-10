package ir

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/parser/utils"
)

type IRIf struct {
	generic_ast.BaseASTNode
	ConditionType IRType `"If" @@`
	Condition     string `@Ident`
	BlockThen     int    `"jump" "to" @Int`
	BlockElse     int    `"else" @Int`
	ParentNode    generic_ast.TraversableNode
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

func (ast *IRIf) Print(c *context.ParsingContext) string {
	return utils.PrintASTNode(c, ast, "If %s %s jump to block_%d else block_%d", ast.ConditionType, ast.Condition, ast.BlockThen, ast.BlockElse)
}

func (ast *IRIf) GetUsedVariables(vars cfg.VariableSet, visitedMap map[generic_ast.TraversableNode]struct{}) cfg.VariableSet {
	vars.Add(cfg.NewVariable(ast.Condition, nil))
	return vars
}

func (ast *IRIf) RenameVariables(substUsed, substDecl cfg.VariableSubstitution) {
	ast.Condition = substUsed.Replace(ast.Condition)
}
