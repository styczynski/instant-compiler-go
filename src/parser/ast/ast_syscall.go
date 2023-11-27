package ast

import (
	"strings"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Syscall struct {
	generic_ast.BaseASTNode
	ReturnType   Type          `"syscall" @@`
	Target       string        `@Ident`
	AppToken     string        `( @"("`
	Arguments    []*Expression `(@@ ("," @@)*)? ")" )`
	ParentNode   generic_ast.TraversableNode
	ResolvedType hindley_milner.Type
}

func (ast *Syscall) OnTypeReturned(t hindley_milner.Type) {
	ast.ResolvedType = t
}

func (ast *Syscall) ExtractConst() (generic_ast.TraversableNode, bool) {
	return nil, false
}

func (ast *Syscall) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *Syscall) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *Syscall) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Syscall) End() lexer.Position {
	return ast.EndPos
}

func (ast *Syscall) GetNode() interface{} {
	return ast
}

func (ast *Syscall) GetChildren() []generic_ast.TraversableNode {
	nodes := make([]generic_ast.TraversableNode, len(ast.Arguments)+1)
	for _, child := range ast.Arguments {
		nodes = append(nodes, child)
	}
	return nodes
}

func (ast *Syscall) Print(c *context.ParsingContext) string {
	args := []string{}
	for _, argument := range ast.Arguments {
		args = append(args, argument.Print(c))
	}
	return printNode(c, ast, "syscall %s(%s)", ast.Target, strings.Join(args, ", "))
}

////

func (ast *Syscall) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	args := []*Expression{}
	for _, arg := range ast.Arguments {
		args = append(args, mapper(ast, arg, context, false).(*Expression))
	}
	return mapper(parent, &Syscall{
		BaseASTNode: ast.BaseASTNode,
		Target:      ast.Target,
		Arguments:   args,
		ParentNode:  parent.(generic_ast.TraversableNode),
	}, context, true)
}

func (ast *Syscall) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	for _, arg := range ast.Arguments {
		mapper(ast, arg, context)
	}
	mapper(parent, ast, context)
}

func (ast *Syscall) Body() generic_ast.Expression {
	return hindley_milner.Batch{}
}

func (ast *Syscall) EmbeddedType(c hindley_milner.InferContext) *hindley_milner.Scheme {
	return ast.ReturnType.GetType(nil)
	//return hindley_milner.NewScheme(nil, CreatePrimitive(T_VOID))
}

func (ast *Syscall) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_TYPE
}

///

func (ast *Syscall) RenameVariables(subst cfg.VariableSubstitution) {
	// if ast.IsApplication() {
	// 	v := subst.Replace(*ast.Target)
	// 	ast.Target = &v
	// }
}

func (ast *Syscall) GetUsedVariables(vars cfg.VariableSet, visitedMap map[generic_ast.TraversableNode]struct{}) cfg.VariableSet {
	vars.Add(cfg.NewVariable(ast.Target, nil))
	for _, arg := range ast.Arguments {
		vars.Insert(cfg.GetAllUsagesVariables(arg, map[generic_ast.TraversableNode]struct{}{}))
	}
	return vars
}
