package ir

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/parser/utils"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type IRMacroCall struct {
	generic_ast.BaseASTNode
	Type       IRType `@Ident`
	Var        string
	MacroName  string
	Data       map[string]interface{}
	ParentNode generic_ast.TraversableNode

	TargetName *string
}

func (ast *IRMacroCall) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *IRMacroCall) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *IRMacroCall) Begin() lexer.Position {
	return ast.Pos
}

func (ast *IRMacroCall) End() lexer.Position {
	return ast.EndPos
}

func (ast *IRMacroCall) GetNode() interface{} {
	return ast
}

func (ast *IRMacroCall) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{
		generic_ast.MakeTraversableNodeToken(ast, ast.MacroName, ast.Pos, ast.EndPos),
	}
}

func (ast *IRMacroCall) Print(c *context.ParsingContext) string {
	if ast.TargetName != nil {
		return utils.PrintASTNode(c, ast, "%s %s = CallMacro[%s](%v %v)", ast.Type, *ast.TargetName, ast.MacroName, ast.Type, ast.Var)
	}
	return utils.PrintASTNode(c, ast, "CallMacro[%s](%v %v)", ast.MacroName, ast.Type, ast.Var)
}

//

func (ast *IRMacroCall) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	return mapper(parent, &IRMacroCall{
		BaseASTNode: ast.BaseASTNode,
		MacroName:   ast.MacroName,
		Var:         ast.Var,
		Type:        ast.Type,
		TargetName:  ast.TargetName,
		ParentNode:  parent.(generic_ast.TraversableNode),
	}, context, true)
}

func (ast *IRMacroCall) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(parent, ast, context)
}

func (ast *IRMacroCall) Body() generic_ast.Expression {
	return hindley_milner.Batch{}
}

func (ast *IRMacroCall) GetUsedVariables(vars cfg.VariableSet, visitedMap map[generic_ast.TraversableNode]struct{}) cfg.VariableSet {
	if len(ast.Var) > 0 {
		vars.Add(cfg.NewVariable(ast.Var, nil))
	}
	return vars
}

func (ast *IRMacroCall) RenameVariables(substUsed, substDecl cfg.VariableSubstitution) {
	if len(ast.Var) > 0 {
		ast.Var = substDecl.Replace(ast.Var)
	}
	if ast.TargetName != nil {
		newVal := substDecl.Replace(*ast.TargetName)
		ast.TargetName = &newVal
	}
}

func (ast *IRMacroCall) GetAssignedVariables(wantMembers bool, visitedMap map[generic_ast.TraversableNode]struct{}) cfg.VariableSet {
	if ast.TargetName != nil {
		return cfg.NewVariableSet(cfg.NewVariable(*ast.TargetName, nil))
	}
	return cfg.NewVariableSet()
}

func (ast *IRMacroCall) GetDeclaredVariables(visitedMap map[generic_ast.TraversableNode]struct{}) cfg.VariableSet {
	if ast.TargetName != nil {
		return cfg.NewVariableSet(cfg.NewVariable(*ast.TargetName, nil))
	}
	return cfg.NewVariableSet()
}
