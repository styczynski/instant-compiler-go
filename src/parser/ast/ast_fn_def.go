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

func (ast *FnDef) canBeInputType(t hindley_milner.Type) bool {
	return !( t.Eq(CreatePrimitive(T_VOID_ARG)) || t.Eq(CreatePrimitive(T_VOID)) )
}

func (ast *FnDef) Validate(c *context.ParsingContext) generic_ast.NodeError {
	fmt.Printf("Validate! HUJUPIZDO %s\n", ast.Name)
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
		for _, arg := range ast.Arg {
			t, _ := arg.ArgumentType.GetType().Type()
			if !ast.canBeInputType(t) {
				filteredArgsStrs := []string{}
				for _, newArg := range ast.Arg {
					if ast.canBeInputType(t) {
						filteredArgsStrs = append(filteredArgsStrs, newArg.Print(c))
					}
				}
				message := fmt.Sprintf("Functions cannot accept type %s as a parameter. Did you specified it by a mistake? If so please try: %s %s(%s) { ... }",
					t.String(),
					ast.ReturnType.Print(c),
					ast.Name,
					strings.Join(filteredArgsStrs, ", "),)
				return generic_ast.NewNodeError(
					"Type not allowed",
					ast,
					message,
					message)
			} else {
				fmt.Printf("NOPE HUJU %s\n", t.String())
			}
		}
	} else {
		message := fmt.Sprintf("%s(...) function is defined in inner scope. Inner functions are not supported.", ast.Name)
		return generic_ast.NewNodeError(
			"Inner function",
			ast,
			message,
			message)
	}
	return nil
}

//////

func (ast *FnDef) GetDeclarationType() *hindley_milner.Scheme {
	//if len(ast.Arg) == 0 {
	//	hindley_milner.NewScheme(hindley_milner.TypeVarSet{hindley_milner.TVar('a')})
	//	return hindley_milner.NamesWithTypesFromMap([]string{""}, map[string]*hindley_milner.Scheme{
	//		"void": hindley_milner.NewScheme(nil, CreatePrimitive(T_VOID)),
	//	})
	//}

	argCount := int16(len(ast.Arg))
	if argCount == 0 {
		argCount = 1
	}
	signature := []hindley_milner.Type{}
	vars := []hindley_milner.TypeVariable{}
	for i:=int16(0); i<argCount+1; i++ {
		signature = append(signature, hindley_milner.TVar(i))
		vars = append(vars, hindley_milner.TVar(i))
	}
	s := hindley_milner.NewScheme(hindley_milner.TypeVarSet(vars), hindley_milner.NewFnType(signature...))
	return s
}

func (ast *FnDef) Args() hindley_milner.NameGroup {
	if len(ast.Arg) == 0 {
		return hindley_milner.NamesWithTypesFromMap([]string{""}, map[string]*hindley_milner.Scheme{
			"void": hindley_milner.NewScheme(nil, CreatePrimitive(T_VOID_ARG)),
		})
	}
	argsTypes := map[string]*hindley_milner.Scheme{}
	names := []string{}
	for _, arg := range ast.Arg {
		argsTypes[arg.Name] = arg.ArgumentType.GetType()
		names = append(names, arg.Name)
	}
	return hindley_milner.NamesWithTypesFromMap(names, argsTypes)
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