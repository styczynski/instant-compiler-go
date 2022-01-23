package ir

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/parser/utils"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type IRCall struct {
	generic_ast.BaseASTNode
	Type           IRType `@Ident`
	TargetName     string `@Ident "="`
	CallTarget     string
	CallTargetType IRType
	ArgumentsTypes []IRType
	Arguments      []string
	IsBuiltin      bool
	ParentNode     generic_ast.TraversableNode
}

func (ast *IRCall) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *IRCall) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *IRCall) Begin() lexer.Position {
	return ast.Pos
}

func (ast *IRCall) End() lexer.Position {
	return ast.EndPos
}

func (ast *IRCall) GetNode() interface{} {
	return ast
}

func (ast *IRCall) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{
		generic_ast.MakeTraversableNodeToken(ast, ast.TargetName, ast.Pos, ast.EndPos),
		generic_ast.MakeTraversableNodeToken(ast, ast.CallTarget, ast.Pos, ast.EndPos),
	}
}

func (ast *IRCall) Print(c *context.ParsingContext) string {
	descr := []string{}
	for i, arg := range ast.Arguments {
		descr = append(descr, fmt.Sprintf("%s %s", string(ast.ArgumentsTypes[i]), arg))
	}
	return utils.PrintASTNode(c, ast, "%s %s = Call(%s %s) (%s)", ast.Type, ast.TargetName, ast.CallTargetType, ast.CallTarget, strings.Join(descr, ","))
}

//

func (ast *IRCall) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	return mapper(parent, &IRCall{
		BaseASTNode:    ast.BaseASTNode,
		CallTarget:     ast.CallTarget,
		CallTargetType: ast.CallTargetType,
		Arguments:      ast.Arguments,
		TargetName:     ast.TargetName,
		ParentNode:     parent.(generic_ast.TraversableNode),
	}, context, true)
}

func (ast *IRCall) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(parent, ast, context)
}

func (ast *IRCall) Body() generic_ast.Expression {
	return hindley_milner.Batch{}
}

func (ast *IRCall) GetAssignedVariables(wantMembers bool, visitedMap map[generic_ast.TraversableNode]struct{}) cfg.VariableSet {
	return cfg.NewVariableSet(cfg.NewVariable(ast.TargetName, nil))
}

func (ast *IRCall) GetUsedVariables(vars cfg.VariableSet, visitedMap map[generic_ast.TraversableNode]struct{}) cfg.VariableSet {
	for _, arg := range ast.Arguments {
		vars.Add(cfg.NewVariable(arg, nil))
	}
	vars.Add(cfg.NewVariable(ast.CallTarget, nil))
	return vars
}

func (ast *IRCall) GetDeclaredVariables(visitedMap map[generic_ast.TraversableNode]struct{}) cfg.VariableSet {
	return cfg.NewVariableSet(cfg.NewVariable(ast.TargetName, nil))
}

func (ast *IRCall) RenameVariables(substUsed, substDecl cfg.VariableSubstitution) {
	for i, arg := range ast.Arguments {
		ast.Arguments[i] = substUsed.Replace(arg)
	}
	ast.CallTarget = substUsed.Replace(ast.CallTarget)
	ast.TargetName = substDecl.Replace(ast.TargetName)
}
