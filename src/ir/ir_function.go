package ir

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/parser/utils"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type IRFunction struct {
	generic_ast.BaseASTNode
	ReturnType   IRType `"Function" @Ident`
	Name         string
	Args         []string
	ArgsTypes    []IRType
	FunctionBody []*IRBlock
	ParentNode   generic_ast.TraversableNode
}

func (ast *IRFunction) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *IRFunction) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *IRFunction) Begin() lexer.Position {
	return ast.Pos
}

func (ast *IRFunction) End() lexer.Position {
	return ast.EndPos
}

func (ast *IRFunction) GetNode() interface{} {
	return ast
}

func (ast *IRFunction) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{}
}

func (ast *IRFunction) Print(c *context.ParsingContext) string {
	argsStrs := []string{}
	for i, argName := range ast.Args {
		argsStrs = append(argsStrs, fmt.Sprintf("%s %s", ast.ArgsTypes[i], argName))
	}
	blocksStrs := []string{}
	c.BlockPush()
	for _, block := range ast.FunctionBody {
		blocksStrs = append(blocksStrs, block.Print(c))
	}
	c.BlockPop()
	return utils.PrintASTNode(c, ast, "Function %s %s(%s)\n{\n%s\n}\n", ast.ReturnType, ast.Name, strings.Join(argsStrs, ", "), strings.Join(blocksStrs, "\n"))
}

///

func (ast *IRFunction) Body() generic_ast.Expression {
	return hindley_milner.Batch{Exp: []generic_ast.Expression{}}
}

func (ast *IRFunction) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	return mapper(parent, &IRFunction{
		BaseASTNode: ast.BaseASTNode,
		ParentNode:  parent.(generic_ast.TraversableNode),
	}, context, true)
}

func (ast *IRFunction) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(parent, ast, context)
}
