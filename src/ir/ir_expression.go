package ir

import (
	"strings"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/parser/utils"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type IRExpression struct {
	generic_ast.BaseASTNode
	Type           IRType   `@Ident`
	TargetName     string   `@Ident "="`
	Operation      string   `@Ident ( "("`
	ArgumentsTypes []IRType `(@@ ("," @@)*)? ")" )?`
	Arguments      []string `(@@ ("," @@)*)? ")" )?`
	ParentNode     generic_ast.TraversableNode
}

func (ast *IRExpression) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *IRExpression) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *IRExpression) Begin() lexer.Position {
	return ast.Pos
}

func (ast *IRExpression) End() lexer.Position {
	return ast.EndPos
}

func (ast *IRExpression) GetNode() interface{} {
	return ast
}

func (ast *IRExpression) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{
		generic_ast.MakeTraversableNodeToken(ast, ast.TargetName, ast.Pos, ast.EndPos),
		generic_ast.MakeTraversableNodeToken(ast, ast.Operation, ast.Pos, ast.EndPos),
	}
}

func (ast *IRExpression) Print(c *context.ParsingContext) string {
	args := []string{}
	argsTypes := []string{}
	for i, arg := range ast.Arguments {
		args = append(args, arg)
		argsTypes = append(argsTypes, string(ast.ArgumentsTypes[i]))
	}
	return utils.PrintASTNode(c, ast, "%s %s = (%s) %s(%s)", ast.Type, ast.TargetName, strings.Join(argsTypes, ","), ast.Operation, strings.Join(args, ","))
}

//

func (ast *IRExpression) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	return mapper(parent, &IRExpression{
		BaseASTNode: ast.BaseASTNode,
		Operation:   ast.Operation,
		Arguments:   ast.Arguments,
		TargetName:  ast.TargetName,
		ParentNode:  parent.(generic_ast.TraversableNode),
	}, context, true)
}

func (ast *IRExpression) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(parent, ast, context)
}

func (ast *IRExpression) Body() generic_ast.Expression {
	return hindley_milner.Batch{}
}

func (ast *IRExpression) GetAssignedVariables(wantMembers bool, visitedMap map[generic_ast.TraversableNode]struct{}) cfg.VariableSet {
	return cfg.NewVariableSet(cfg.NewVariable(ast.TargetName, nil))
}

func (ast *IRExpression) GetUsedVariables(vars cfg.VariableSet, visitedMap map[generic_ast.TraversableNode]struct{}) cfg.VariableSet {
	for _, arg := range ast.Arguments {
		vars.Add(cfg.NewVariable(arg, nil))
	}
	return vars
}

func (ast *IRExpression) RenameVariables(substUsed, substDecl cfg.VariableSubstitution) {
	for i, arg := range ast.Arguments {
		ast.Arguments[i] = substUsed.Replace(arg)
	}
	ast.TargetName = substDecl.Replace(ast.TargetName)
}
