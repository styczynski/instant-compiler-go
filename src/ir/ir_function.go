package ir

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
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
	nodes := []generic_ast.TraversableNode{}
	for _, child := range ast.FunctionBody {
		nodes = append(nodes, child)
	}
	return nodes
}

///////

func (ast *IRFunction) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	mappedStmts := []*IRBlock{}
	for _, stmt := range ast.FunctionBody {
		mappedStmts = append(mappedStmts, mapper(ast, stmt, context, false).(*IRBlock))
	}
	return mapper(parent, &IRFunction{
		BaseASTNode:  ast.BaseASTNode,
		ReturnType:   ast.ReturnType,
		Name:         ast.Name,
		Args:         ast.Args,
		ArgsTypes:    ast.ArgsTypes,
		FunctionBody: mappedStmts,
		ParentNode:   parent.(generic_ast.TraversableNode),
	}, context, true)
}

func (ast *IRFunction) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	for _, stmt := range ast.FunctionBody {
		mapper(ast, stmt, context)
	}
	mapper(parent, ast, context)
}

func (ast *IRFunction) GetContents() hindley_milner.Batch {
	exprs := []generic_ast.Expression{}
	for _, stmt := range ast.FunctionBody {
		exprs = append(exprs, stmt)
	}
	return hindley_milner.Batch{
		Exp: exprs,
	}
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

func (ast *IRFunction) Body() generic_ast.Expression {
	exps := []generic_ast.Expression{}
	for _, stmt := range ast.FunctionBody {
		exps = append(exps, stmt)
	}
	return hindley_milner.Batch{Exp: exps}
}

func (ast *IRFunction) BuildFlowGraph(builder cfg.CFGBuilder) {
	builder.BuildBlock(ast)
}
