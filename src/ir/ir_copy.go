package ir

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/parser/utils"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type IRCopy struct {
	generic_ast.BaseASTNode
	Type       IRType `@Ident`
	TargetName string `@Ident "="`
	Var        string
	ParentNode generic_ast.TraversableNode
}

func (ast *IRCopy) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *IRCopy) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *IRCopy) Begin() lexer.Position {
	return ast.Pos
}

func (ast *IRCopy) End() lexer.Position {
	return ast.EndPos
}

func (ast *IRCopy) GetNode() interface{} {
	return ast
}

func (ast *IRCopy) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{
		generic_ast.MakeTraversableNodeToken(ast, ast.TargetName, ast.Pos, ast.EndPos),
		generic_ast.MakeTraversableNodeToken(ast, ast.Var, ast.Pos, ast.EndPos),
	}
}

func (ast *IRCopy) Print(c *context.ParsingContext) string {
	return utils.PrintASTNode(c, ast, "%s %s = Copy(%s)", ast.Type, ast.TargetName, ast.Var)
}

//

func (ast *IRCopy) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	return mapper(parent, &IRCopy{
		BaseASTNode: ast.BaseASTNode,
		Var:         ast.Var,
		TargetName:  ast.TargetName,
		ParentNode:  parent.(generic_ast.TraversableNode),
	}, context, true)
}

func (ast *IRCopy) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(parent, ast, context)
}

func (ast *IRCopy) Body() generic_ast.Expression {
	return hindley_milner.Batch{}
}

func (ast *IRCopy) GetDeclaredVariables(visitedMap map[generic_ast.TraversableNode]struct{}) cfg.VariableSet {
	return cfg.NewVariableSet(cfg.NewVariable(ast.TargetName, nil))
}

func (ast *IRCopy) GetAssignedVariables(wantMembers bool, visitedMap map[generic_ast.TraversableNode]struct{}) cfg.VariableSet {
	return cfg.NewVariableSet(cfg.NewVariable(ast.TargetName, nil))
}

func (ast *IRCopy) GetUsedVariables(vars cfg.VariableSet, visitedMap map[generic_ast.TraversableNode]struct{}) cfg.VariableSet {
	vars.Add(cfg.NewVariable(ast.Var, nil))
	return vars
}

func (ast *IRCopy) RenameVariables(substUsed, substDecl cfg.VariableSubstitution) {
	ast.Var = substUsed.Replace(ast.Var)
	ast.TargetName = substDecl.Replace(ast.TargetName)
}
