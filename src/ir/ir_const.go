package ir

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/parser/utils"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type IRConst struct {
	generic_ast.BaseASTNode
	Type       IRType `@Ident`
	TargetName string `@Ident "="`
	Value      int64
	ParentNode generic_ast.TraversableNode
}

func (ast *IRConst) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *IRConst) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *IRConst) Begin() lexer.Position {
	return ast.Pos
}

func (ast *IRConst) End() lexer.Position {
	return ast.EndPos
}

func (ast *IRConst) GetNode() interface{} {
	return ast
}

func (ast *IRConst) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{
		generic_ast.MakeTraversableNodeToken(ast, ast.TargetName, ast.Pos, ast.EndPos),
	}
}

func (ast *IRConst) Print(c *context.ParsingContext) string {
	return utils.PrintASTNode(c, ast, "%s %s = Const(%v)", ast.Type, ast.TargetName, ast.Value)
}

//

func (ast *IRConst) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	return mapper(parent, &IRConst{
		BaseASTNode: ast.BaseASTNode,
		TargetName:  ast.TargetName,
		Value:       ast.Value,
		ParentNode:  parent.(generic_ast.TraversableNode),
	}, context, true)
}

func (ast *IRConst) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(parent, ast, context)
}

func (ast *IRConst) Body() generic_ast.Expression {
	return hindley_milner.Batch{}
}

func (ast *IRConst) GetAssignedVariables(wantMembers bool, visitedMap map[generic_ast.TraversableNode]struct{}) cfg.VariableSet {
	return cfg.NewVariableSet(cfg.NewVariable(ast.TargetName, nil))
}

func (ast *IRConst) GetUsedVariables(vars cfg.VariableSet, visitedMap map[generic_ast.TraversableNode]struct{}) cfg.VariableSet {
	return vars
}

func (ast *IRConst) GetDeclaredVariables(visitedMap map[generic_ast.TraversableNode]struct{}) cfg.VariableSet {
	return cfg.NewVariableSet(cfg.NewVariable(ast.TargetName, nil))
}

func (ast *IRConst) RenameVariables(substUsed, substDecl cfg.VariableSubstitution) {
	ast.TargetName = substDecl.Replace(ast.TargetName)
}
