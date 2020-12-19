package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Assignment struct {
	generic_ast.BaseASTNode
	TargetName string `@Ident`
	Value *Expression `"=" @@ ";"`
	ParentNode generic_ast.TraversableNode
}

func (ast *Assignment) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *Assignment) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *Assignment) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Assignment) End() lexer.Position {
	return ast.EndPos
}

func (ast *Assignment) GetNode() interface{} {
	return ast
}

func (ast *Assignment) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{
		generic_ast.MakeTraversableNodeToken(ast, ast.TargetName, ast.Pos, ast.EndPos),
		ast.Value,
	}
}

func (ast *Assignment) Print(c *context.ParsingContext) string {
	return printNode(c, ast, "%s = %s;", ast.TargetName, ast.Value.Print(c))
}

//

func (ast *Assignment) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	return mapper(parent, &Assignment{
		BaseASTNode: ast.BaseASTNode,
		Value:    mapper(ast, ast.Value, context, false).(*Expression),
		TargetName: ast.TargetName,
		ParentNode: parent.(generic_ast.TraversableNode),
	}, context, true)
}

func (ast *Assignment) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(ast, ast.Value, context)
	mapper(parent, ast, context)
}

func (ast *Assignment) Fn() generic_ast.Expression {
	//return &BuiltinFunction{
	//	BaseASTNode: ast.BaseASTNode,
	//	name: "=",
	//}
	return &hindley_milner.EmbeddedTypeExpr{GetType: func() *hindley_milner.Scheme {
		return hindley_milner.NewScheme(
			hindley_milner.TypeVarSet{hindley_milner.TVar('a')},
			hindley_milner.NewFnType(hindley_milner.TVar('a'), hindley_milner.TVar('a'), hindley_milner.TVar('a')))
	}}
}

func (ast *Assignment) Body() generic_ast.Expression {
	return hindley_milner.Batch{
		Exp: []generic_ast.Expression{
			&VarName{
				BaseASTNode: ast.BaseASTNode,
				name: ast.TargetName,
			},
			ast.Value,
		},
	}
}

func (ast *Assignment) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_TYPE_EQUALITY
}

// Validate here this shit

func (ast *Assignment) GetAssignedVariables(wantMembers bool, visitedMap map[generic_ast.TraversableNode]struct{}) cfg.VariableSet {
	return cfg.NewVariableSet(cfg.NewVariable(ast.TargetName, ast.Value))
}

func (ast *Assignment) RenameVariables(subst cfg.VariableSubstitution) {
	ast.TargetName = subst.Replace(ast.TargetName)
}

//func (ast *Assignment) GetAssignedVariables(wantMembers bool) cfg.VariableSet {
//	if ast.HasIndexingExpr() {
//		if wantMembers {
//			return cfg.GetAllAssignedVariables(ast.Primary, wantMembers)
//		} else {
//			return cfg.NewVariableSet()
//		}
//	}
//	return cfg.GetAllVariables(ast.Primary)
//}