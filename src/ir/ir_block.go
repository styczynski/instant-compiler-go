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

type IRBlock struct {
	generic_ast.BaseASTNode
	Statements []*IRStatement `"{" @@* "}"`
	ParentNode generic_ast.TraversableNode
	BlockID    int
	Label      string
}

func (ast *IRBlock) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *IRBlock) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *IRBlock) Begin() lexer.Position {
	return ast.Pos
}

func (ast *IRBlock) End() lexer.Position {
	return ast.EndPos
}

func (ast *IRBlock) GetNode() interface{} {
	return ast
}

func (ast *IRBlock) Join(otherBlock *IRBlock) {
	ast.Statements = append(ast.Statements, otherBlock.Statements...)
}

func (ast *IRBlock) Print(c *context.ParsingContext) string {
	statementsList := []string{}
	if c.PrinterConfiguration.SkipStatementIdent {
		c.PrinterConfiguration.SkipStatementIdent = false
	}
	c.BlockPush()
	for _, statement := range ast.Statements {
		statementsList = append(statementsList, statement.Print(c))
	}
	c.BlockPop()
	labelStr := ""
	if len(ast.Label) > 0 {
		labelStr = fmt.Sprintf("block_%d: ; %s", ast.BlockID, ast.Label)
	}
	return utils.PrintASTNode(c, ast, "%s%s\n%s\n%s", strings.Repeat("  ", c.BlockDepth), labelStr, strings.Join(statementsList, "\n"), strings.Repeat("  ", c.BlockDepth))
}

func (ast *IRBlock) GetChildren() []generic_ast.TraversableNode {
	nodes := []generic_ast.TraversableNode{}
	for _, child := range ast.Statements {
		nodes = append(nodes, child)
	}
	return nodes
}

///////

func (ast *IRBlock) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	mappedStmts := []*IRStatement{}
	for _, stmt := range ast.Statements {
		mappedStmts = append(mappedStmts, mapper(ast, stmt, context, false).(*IRStatement))
	}
	return mapper(parent, &IRBlock{
		BaseASTNode: ast.BaseASTNode,
		Statements:  mappedStmts,
		ParentNode:  parent.(generic_ast.TraversableNode),
	}, context, true)
}

func (ast *IRBlock) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	for _, stmt := range ast.Statements {
		mapper(ast, stmt, context)
	}
	mapper(parent, ast, context)
}

func (ast *IRBlock) GetContents() hindley_milner.Batch {
	exprs := []generic_ast.Expression{}
	for _, stmt := range ast.Statements {
		exprs = append(exprs, stmt)
	}
	return hindley_milner.Batch{
		Exp: exprs,
	}
}

func (ast *IRBlock) Expressions() []generic_ast.Expression {
	return ast.GetContents().Exp
}

func (ast *IRBlock) Body() generic_ast.Expression {
	panic(fmt.Errorf("Batch Body() method cannot be called."))
}

//

func (ast *IRBlock) BuildFlowGraph(builder cfg.CFGBuilder) {
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
