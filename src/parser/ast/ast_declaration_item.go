package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type DeclarationItem struct {
	generic_ast.BaseASTNode
	Name string `@Ident`
	Initializer *Expression `( "=" @@ )?`
	ParentNode generic_ast.TraversableNode
}

func (ast *DeclarationItem) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *DeclarationItem) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *DeclarationItem) Begin() lexer.Position {
	return ast.Pos
}

func (ast *DeclarationItem) End() lexer.Position {
	return ast.EndPos
}

func (ast *DeclarationItem) GetNode() interface{} {
	return ast
}

func (ast *DeclarationItem) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{
		generic_ast.MakeTraversableNodeToken(ast, ast.Name, ast.Pos, ast.EndPos),
		ast.Initializer,
	}
}

func (ast *DeclarationItem) HasInitializer() bool {
	return ast.Initializer != nil
}

func (ast *DeclarationItem) Print(c *context.ParsingContext) string {
	if ast.HasInitializer() {
		return printNode(c, ast, "%s = %s", ast.Name, ast.Initializer.Print(c))
	}
	return printNode(c, ast, "%s", ast.Name)
}

/////

func (ast *DeclarationItem) Body() generic_ast.Expression {
	if !ast.HasInitializer() {
		return hindley_milner.Batch{Exp: []generic_ast.Expression{}}
	}
	return ast.Initializer
}

func (ast *DeclarationItem) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	return mapper(parent, &DeclarationItem{
		BaseASTNode: ast.BaseASTNode,
		Name:        ast.Name,
		Initializer: mapper(ast, ast.Initializer, context, false).(*Expression),
		ParentNode: parent.(generic_ast.TraversableNode),
	}, context, true)
}

func (ast *DeclarationItem) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	if ast.HasInitializer() {
		mapper(ast, ast.Initializer, context)
	}
	mapper(parent, ast, context)
}

func (ast *DeclarationItem) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_PROXY
}

//

func (ast *DeclarationItem) GetDeclaredVariables() cfg.VariableSet {
	return cfg.NewVariableSet(cfg.NewVariable(ast.Name, ast.Initializer))
}

func (ast *DeclarationItem) GetUsedVariables(vars cfg.VariableSet) cfg.VariableSet {
	if !ast.HasInitializer() {
		return cfg.NewVariableSet()
	}
	return cfg.GetAllUsagesVariables(ast.Initializer)
}


func (ast *DeclarationItem) RenameVariables(subst cfg.VariableSubstitution) {
	ast.Name = subst.Replace(ast.Name)
}