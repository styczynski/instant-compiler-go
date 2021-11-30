package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type UnaryStatement struct {
	generic_ast.BaseASTNode
	TargetName *string `@Ident`
	Operation  string  `@( "+" "+" | "-" "-" ) ";"`
	ParentNode generic_ast.TraversableNode
}

func (ast *UnaryStatement) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *UnaryStatement) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *UnaryStatement) Begin() lexer.Position {
	return ast.Pos
}

func (ast *UnaryStatement) End() lexer.Position {
	return ast.EndPos
}

func (ast *UnaryStatement) GetNode() interface{} {
	return ast
}

func (ast *UnaryStatement) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{
		generic_ast.MakeTraversableNodeToken(ast, *ast.TargetName, ast.Pos, ast.EndPos),
		generic_ast.MakeTraversableNodeToken(ast, ast.Operation, ast.Pos, ast.EndPos),
	}
}

func (ast *UnaryStatement) Print(c *context.ParsingContext) string {
	return printNode(c, ast, "%s%s;", *ast.TargetName, ast.Operation)
}

///

func (ast *UnaryStatement) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	return mapper(parent, &UnaryStatement{
		BaseASTNode: ast.BaseASTNode,
		TargetName:  ast.TargetName,
		Operation:   ast.Operation,
		ParentNode:  parent.(generic_ast.TraversableNode),
	}, context, true)
}

func (ast *UnaryStatement) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(parent, ast, context)
}

func (ast *UnaryStatement) Fn(c hindley_milner.InferContext) generic_ast.Expression {
	return &BuiltinFunction{
		BaseASTNode: ast.BaseASTNode,
		name:        ast.Operation,
	}
}

func (ast *UnaryStatement) Body() generic_ast.Expression {
	return hindley_milner.Batch{
		Exp: []generic_ast.Expression{
			&VarName{
				BaseASTNode: ast.BaseASTNode,
				name:        *ast.TargetName,
			},
		},
	}
}

func (ast *UnaryStatement) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_APPLICATION
}

///

func (ast *UnaryStatement) GetAssignedVariables(wantMembers bool, visitedMap map[generic_ast.TraversableNode]struct{}) cfg.VariableSet {
	return cfg.NewVariableSet(cfg.NewVariable(*ast.TargetName, nil))
}

func (ast *UnaryStatement) GetUsedVariables(vars cfg.VariableSet, visitedMap map[generic_ast.TraversableNode]struct{}) cfg.VariableSet {
	vars.Add(cfg.NewVariable(*ast.TargetName, nil))
	return vars
}

func (ast *UnaryStatement) RenameVariables(subst cfg.VariableSubstitution, visitedMap map[generic_ast.TraversableNode]struct{}) {
	v := subst.Replace(*ast.TargetName)
	ast.TargetName = &v
}
