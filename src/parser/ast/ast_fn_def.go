package ast

import (
	"strings"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type FnDef struct {
	BaseASTNode
	ReturnType Type `@@`
	Name string `@Ident`
	Arg []*Arg `"(" (@@ ( "," @@ )*)? ")"`
	FunctionBody *Block `@@`
	ParentNode TraversableNode
}

func (ast *FnDef) Parent() TraversableNode {
	return ast.ParentNode
}

func (ast *FnDef) OverrideParent(node TraversableNode) {
	ast.ParentNode = node
}

func (ast *FnDef) Begin() lexer.Position {
	return ast.Pos
}

func (ast *FnDef) End() lexer.Position {
	return ast.EndPos
}

func (ast *FnDef) GetNode() interface{} {
	return ast
}

func (ast *FnDef) GetChildren() []TraversableNode {
	nodes := make([]TraversableNode, len(ast.Arg) + 3)
	nodes = append(nodes, &ast.ReturnType)
	nodes = append(nodes, MakeTraversableNodeToken(ast, ast.Name, ast.Pos, ast.EndPos))

	for _, child := range ast.Arg {
		nodes = append(nodes, child)
	}
	nodes = append(nodes, ast.FunctionBody)

	return nodes
}

func (ast *FnDef) Print(c *context.ParsingContext) string {
	argsList := []string{}
	for _, arg := range ast.Arg {
		argsList = append(argsList, arg.Print(c))
	}

	return printNode(c, ast, "%s %s(%s) %s",
		ast.ReturnType.Print(c),
		ast.Name,
		strings.Join(argsList, ", "),
		ast.FunctionBody.Print(c))
}

//////

func (ast *FnDef) Args() hindley_milner.NameGroup {
	if len(ast.Arg) == 0 {
		return hindley_milner.NamesWithTypesFromMap(map[string]*hindley_milner.Scheme{
			"void": hindley_milner.NewScheme(nil, CreatePrimitive(T_VOID)),
		})
	}
	argsTypes := map[string]*hindley_milner.Scheme{}
	for _, arg := range ast.Arg {
		argsTypes[arg.Name] = arg.ArgumentType.GetType()
	}
	return hindley_milner.NamesWithTypesFromMap(argsTypes)
}

func (ast *FnDef) Var() hindley_milner.NameGroup {
	return hindley_milner.Name(ast.Name)
}

func (ast *FnDef) Body() hindley_milner.Expression {
	return ast.FunctionBody
}

func (ast *FnDef) ExpressionType() hindley_milner.ExpressionType { return hindley_milner.E_FUNCTION_DECLARATION }

func (ast *FnDef) DefaultType() *hindley_milner.Scheme {
	return ast.ReturnType.GetType()
}

func (ast *FnDef) Map(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	return mapper(parent, &FnDef{
		BaseASTNode:  ast.BaseASTNode,
		ReturnType:   ast.ReturnType,
		Name:         ast.Name,
		Arg:          ast.Arg,
		FunctionBody: mapper(ast, ast.FunctionBody).(*Block),
		ParentNode: parent.(TraversableNode),
	})
}
func (ast *FnDef) Visit(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) {
	mapper(ast, ast.FunctionBody)
	mapper(parent, ast)
}