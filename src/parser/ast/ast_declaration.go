package ast

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Declaration struct {
	generic_ast.BaseASTNode
	DeclarationType Type               `@@`
	Items           []*DeclarationItem `( @@ ( "," @@ )* ) ";"`
	ParentNode      generic_ast.TraversableNode
}

func (ast *Declaration) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *Declaration) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *Declaration) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Declaration) End() lexer.Position {
	return ast.EndPos
}

func (ast *Declaration) GetNode() interface{} {
	return ast
}

func (ast *Declaration) GetChildren() []generic_ast.TraversableNode {
	nodes := make([]generic_ast.TraversableNode, len(ast.Items)+1)
	nodes = append(nodes, &ast.DeclarationType)
	for _, child := range ast.Items {
		nodes = append(nodes, child)
	}
	return nodes
}

func (ast *Declaration) Print(c *context.ParsingContext) string {
	declarationItemsList := []string{}
	for _, item := range ast.Items {
		declarationItemsList = append(declarationItemsList, item.Print(c))
	}
	return printNode(c, ast, "%s %s", ast.DeclarationType.Print(c), strings.Join(declarationItemsList, ", "))
}

func (ast *Declaration) canBeInputType(t hindley_milner.Type) bool {
	return !(t.Eq(CreatePrimitive(T_VOID_ARG)) || t.Eq(CreatePrimitive(T_VOID)))
}

func (ast *Declaration) Validate(c *context.ParsingContext) generic_ast.NodeError {
	t, _ := ast.DeclarationType.GetType(nil).Type()
	if !ast.canBeInputType(t) {
		message := fmt.Sprintf("Declarations cannot set variable with type %s. Change the type to other possible alternatives.", t.String())
		return generic_ast.NewNodeError(
			"Type not allowed",
			ast,
			message,
			message)
	}
	return nil
}

//////

func (ast *Declaration) Body() generic_ast.Expression {
	return ast
}

func (ast *Declaration) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	return mapper(parent, &Declaration{
		BaseASTNode:     ast.BaseASTNode,
		DeclarationType: ast.DeclarationType,
		Items:           ast.Items,
		ParentNode:      parent.(generic_ast.TraversableNode),
	}, context, true).(*Declaration)
}

func (ast *Declaration) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	for _, item := range ast.Items {
		mapper(ast, item, context)
	}
	mapper(parent, ast, context)
}

func (ast *Declaration) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_DECLARATION
}

func (ast *Declaration) Var(c hindley_milner.InferContext) *hindley_milner.NameGroup {
	names := []string{}
	types := map[string]*hindley_milner.Scheme{}
	for _, item := range ast.Items {
		names = append(names, item.Name)
		types[item.Name] = ast.DeclarationType.GetType(c)
	}
	return hindley_milner.NamesWithTypes(names, types)
}

func (ast *Declaration) Def(c hindley_milner.InferContext) generic_ast.Expression {
	defs := []generic_ast.Expression{}
	for _, item := range ast.Items {
		defs = append(defs, item)
	}
	return hindley_milner.Batch{
		Exp: defs,
	}
}

func (ast *Declaration) RemoveVariableAssignment(variableNames map[string]struct{}) generic_ast.NormalNode {
	newDecls := []*DeclarationItem{}
	for _, declItem := range ast.Items {
		if _, ok := variableNames[declItem.Name]; !ok {
			newDecls = append(newDecls, declItem)
		}
	}
	ast.Items = newDecls
	if len(ast.Items) == 0 {
		return nil
	}
	return ast
}
