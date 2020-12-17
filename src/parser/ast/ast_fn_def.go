package ast

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type FnDef struct {
	generic_ast.BaseASTNode
	ReturnType Type `@@`
	Name string `@Ident`
	Arg []*Arg `"(" (@@ ( "," @@ )*)? ")"`
	FunctionBody *Block `@@`
	ParentNode generic_ast.TraversableNode
}

func (ast *FnDef) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *FnDef) OverrideParent(node generic_ast.TraversableNode) {
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

func (ast *FnDef) GetChildren() []generic_ast.TraversableNode {
	nodes := make([]generic_ast.TraversableNode, len(ast.Arg) + 3)
	nodes = append(nodes, &ast.ReturnType)
	nodes = append(nodes, generic_ast.MakeTraversableNodeToken(ast, ast.Name, ast.Pos, ast.EndPos))

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

/////

func (ast *FnDef) Validate() generic_ast.NodeError {
	if _, ok := ast.Parent().(*TopDef); ok {
		if ast.Name == "main" {
			returnedType, _ := ast.ReturnType.GetType().Type()
			expectedType := CreatePrimitive(T_INT)
			if !returnedType.Eq(expectedType) {
				message := fmt.Sprintf("main() function has return type of %v. Expected it to be %v.", returnedType, expectedType)
				return generic_ast.NewNodeError(
					"Invalid main()",
					 ast,
					 message,
					 message)
			}
			if len(ast.Arg) != 0 {
				message := fmt.Sprintf("main() function cannot accept any arguments, but it has %d parameter/-s in definiton.", len(ast.Arg))
				return generic_ast.NewNodeError(
					"Invalid main()",
					ast,
					message,
					message)
			}
		}
	}
	return nil
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

func (ast *FnDef) Body() generic_ast.Expression {
	return ast.FunctionBody
}

func (ast *FnDef) ExpressionType() hindley_milner.ExpressionType { return hindley_milner.E_FUNCTION_DECLARATION }

func (ast *FnDef) DefaultType() *hindley_milner.Scheme {
	return ast.ReturnType.GetType()
}

func (ast *FnDef) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	return mapper(parent, &FnDef{
		BaseASTNode:  ast.BaseASTNode,
		ReturnType:   ast.ReturnType,
		Name:         ast.Name,
		Arg:          ast.Arg,
		FunctionBody: mapper(ast, ast.FunctionBody, context, false).(*Block),
		ParentNode: parent.(generic_ast.TraversableNode),
	}, context, true)
}
func (ast *FnDef) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(ast, ast.FunctionBody, context)
	mapper(parent, ast, context)
}

//

func (ast *FnDef) BuildFlowGraph(builder cfg.CFGBuilder) {
	builder.BuildNode(ast.FunctionBody)
}