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

type IRProgram struct {
	generic_ast.BaseASTNode
	Statements []*IRFunction `@@*`
	ParentNode generic_ast.TraversableNode
}

func (ast *IRProgram) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *IRProgram) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *IRProgram) Begin() lexer.Position {
	return ast.Pos
}

func (ast *IRProgram) End() lexer.Position {
	return ast.EndPos
}

func (ast *IRProgram) GetNode() interface{} {
	return ast
}

func (ast *IRProgram) Join(otherBlock *IRProgram) {
	ast.Statements = append(ast.Statements, otherBlock.Statements...)
}

func (ast *IRProgram) Print(c *context.ParsingContext) string {
	statementsList := []string{}
	if c.PrinterConfiguration.SkipStatementIdent {
		c.PrinterConfiguration.SkipStatementIdent = false
	}
	c.BlockPush()
	for _, statement := range ast.Statements {
		statementsList = append(statementsList, statement.Print(c))
	}
	c.BlockPop()
	return utils.PrintASTNode(c, ast, "%s%s\n%s\n%s", strings.Repeat("  ", c.BlockDepth), strings.Join(statementsList, "\n"), strings.Repeat("  ", c.BlockDepth))
}

func (ast *IRProgram) GetChildren() []generic_ast.TraversableNode {
	nodes := []generic_ast.TraversableNode{}
	for _, child := range ast.Statements {
		nodes = append(nodes, child)
	}
	return nodes
}

///////

func (ast *IRProgram) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	mappedStmts := []*IRFunction{}
	for _, stmt := range ast.Statements {
		mappedStmts = append(mappedStmts, mapper(ast, stmt, context, false).(*IRFunction))
	}
	return mapper(parent, &IRProgram{
		BaseASTNode: ast.BaseASTNode,
		Statements:  mappedStmts,
		ParentNode:  parent.(generic_ast.TraversableNode),
	}, context, true)
}

func (ast *IRProgram) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	for _, stmt := range ast.Statements {
		mapper(ast, stmt, context)
	}
	mapper(parent, ast, context)
}

func (ast *IRProgram) GetContents() hindley_milner.Batch {
	exprs := []generic_ast.Expression{}
	for _, stmt := range ast.Statements {
		exprs = append(exprs, stmt)
	}
	return hindley_milner.Batch{
		Exp: exprs,
	}
}

func (ast *IRProgram) Expressions() []generic_ast.Expression {
	return ast.GetContents().Exp
}

func (ast *IRProgram) Body() generic_ast.Expression {
	panic(fmt.Errorf("Batch Body() method cannot be called."))
}

//

func (ast *IRProgram) BuildFlowGraph(builder cfg.CFGBuilder) {
	builder.BuildBlock(ast)
}


/*

type VirtualBlock struct {
	parent generic_ast.TraversableNode
	nodes  []generic_ast.NormalNode
}

type VirtualBlockType string

func CreateVirtualBlock(nodes []generic_ast.NormalNode) *VirtualBlock {
	return &VirtualBlock{
		parent: nil,
		nodes:  nodes,
	}
}

func (n *VirtualBlock) Join(b *VirtualBlock) {
	n.nodes = append(n.nodes, b.nodes...)
}

func (n *VirtualBlock) Print(c *context.ParsingContext) string {
	statementsList := []string{}
	if c.PrinterConfiguration.SkipStatementIdent {
		c.PrinterConfiguration.SkipStatementIdent = false
	}
	c.BlockPush()
	for _, statement := range n.nodes {
		statementsList = append(statementsList, statement.Print(c))
	}
	c.BlockPop()
	return PrintASTNode(c, n, "%s", strings.Join(statementsList, "\n"))
}

func (n *VirtualBlock) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{}
}

func (n *VirtualBlock) GetNode() interface{} {
	return n
}

func (n *VirtualBlock) Parent() generic_ast.TraversableNode {
	return n.parent
}

func (n *VirtualBlock) OverrideParent(node generic_ast.TraversableNode) {
	n.parent = node
}

func (n *VirtualBlock) Begin() lexer.Position {
	dummyOffset := 0
	return lexer.Position{
		Filename: "",
		Offset:   dummyOffset,
		Line:     dummyOffset,
		Column:   dummyOffset,
	}
}

func (n *VirtualBlock) End() lexer.Position {
	return n.Begin()
}

*/
