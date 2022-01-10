package ir

import (
	"strings"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/parser/utils"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type IRStatement struct {
	generic_ast.BaseASTNode
	Empty      *IREmpty      `@@`
	Exit       *IRExit       `| @@`
	If         *IRIf         `| @@`
	Phi        *IRPhi        `| @@`
	Expression *IRExpression `| @@`
	ParentNode generic_ast.TraversableNode
}

func WrapIREmpty() *IRStatement {
	return &IRStatement{
		Empty: &IREmpty{},
	}
}

func WrapIRPhi(ast *IRPhi) *IRStatement {
	return &IRStatement{
		Phi: ast,
	}
}

func WrapIRIf(ast *IRIf) *IRStatement {
	return &IRStatement{
		If: ast,
	}
}

func WrapIRExit(ast *IRExit) *IRStatement {
	return &IRStatement{
		Exit: ast,
	}
}

func WrapIRExpression(ast *IRExpression) *IRStatement {
	return &IRStatement{
		Expression: ast,
	}
}

func (ast *IRStatement) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *IRStatement) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *IRStatement) Begin() lexer.Position {
	return ast.Pos
}

func (ast *IRStatement) End() lexer.Position {
	return ast.EndPos
}

func (ast *IRStatement) GetNode() interface{} {
	return ast
}

func (ast *IRStatement) IsEmpty() bool {
	return (ast.Empty != nil || (!ast.IsExit() &&
		!ast.IsIf() &&
		!ast.IsPhi() &&
		!ast.IsExpression()))
}

func (ast *IRStatement) IsExit() bool {
	return ast.Exit != nil
}

func (ast *IRStatement) IsIf() bool {
	return ast.If != nil
}

func (ast *IRStatement) IsPhi() bool {
	return ast.Phi != nil
}

func (ast *IRStatement) IsExpression() bool {
	return ast.Expression != nil
}

func (ast *IRStatement) GetChildren() []generic_ast.TraversableNode {
	if ast.IsEmpty() {
		return []generic_ast.TraversableNode{}
		//return []generic_ast.TraversableNode{ generic_ast.MakeTraversableNodeToken(ast, *ast.Empty, ast.Pos, ast.EndPos) }
	} else if ast.IsExit() {
		return []generic_ast.TraversableNode{ast.Exit}
	} else if ast.IsIf() {
		return []generic_ast.TraversableNode{ast.If}
	} else if ast.IsPhi() {
		return []generic_ast.TraversableNode{ast.Phi}
	} else if ast.IsExpression() {
		return []generic_ast.TraversableNode{ast.Expression}
	}
	return []generic_ast.TraversableNode{}
}

func (ast *IRStatement) formatIRStatementInstruction(irstatement string, c *context.ParsingContext) string {
	if c.PrinterConfiguration.SkipStatementIdent {
		c.PrinterConfiguration.SkipStatementIdent = false
		return irstatement
	}
	return utils.PrintASTNode(c, ast, "%s%s", strings.Repeat("  ", c.BlockDepth), irstatement)
}

func (ast *IRStatement) Print(c *context.ParsingContext) string {
	ret := "UNKNOWN"
	propagateSkipIRStatementIdent := false
	if ast.IsEmpty() {
		ret = ";"
	} else if ast.IsExit() {
		ret = ast.Exit.Print(c)
	} else if ast.IsIf() {
		ret = ast.If.Print(c)
	} else if ast.IsPhi() {
		ret = ast.Phi.Print(c)
	} else if ast.IsExpression() {
		ret = utils.PrintASTNode(c, ast, "%s;", ast.Expression.Print(c))
	}
	if propagateSkipIRStatementIdent {
		c.PrinterConfiguration.SkipStatementIdent = true
		ret = ast.formatIRStatementInstruction(ret, c)
	} else {
		ret = ast.formatIRStatementInstruction(ret, c)
	}
	return ret
}

//////

func (ast *IRStatement) Body() generic_ast.Expression {
	if ast.IsEmpty() {
		return &IREmpty{}
		//return hindley_milner.Batch{Exp: []generic_ast.Expression{}}
	} else if ast.IsExit() {
		return ast.Exit
	} else if ast.IsExpression() {
		return ast.Expression
	}
	ast.BaseASTNode = generic_ast.BaseASTNode{}
	ast.ParentNode = nil
	//fmt.Printf("Failed for node %s\n", repr.String(ast))
	panic("Invalid IRStatement type")
}

func feedExpressionIntoIRStatement(node interface{}, base generic_ast.BaseASTNode) *IRStatement {
	if exitStmt, ok := node.(*IRExit); ok {
		return &IRStatement{
			BaseASTNode: base,
			Exit:        exitStmt,
		}
	} else if ifStmt, ok := node.(*IRIf); ok {
		return &IRStatement{
			BaseASTNode: base,
			If:          ifStmt,
		}
	} else if expr, ok := node.(*IRExpression); ok {
		return &IRStatement{
			BaseASTNode: base,
			Expression:  expr,
		}
	}
	panic("Invalid irstatement type")
}

func (ast *IRStatement) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	mappedIRStatement := feedExpressionIntoIRStatement(mapper(ast, ast.Body(), context, false), ast.BaseASTNode)
	return mapper(parent, mappedIRStatement, context, true)
}

func (ast *IRStatement) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(ast, ast.Body(), context)
	mapper(parent, ast, context)
}

func (ast *IRStatement) ExpressionType() hindley_milner.ExpressionType { return hindley_milner.E_PROXY }

//

func (ast *IRStatement) BuildFlowGraph(builder cfg.CFGBuilder) {
	if ast.IsEmpty() {
		// Do nothing
	} else if ast.IsExit() {
		builder.BuildNode(ast.Exit)
	} else if ast.IsIf() {
		builder.BuildNode(ast.If)
	} else if ast.IsPhi() {
		builder.BuildNode(ast.Phi)
	} else if ast.IsExpression() {
		builder.BuildNode(ast.Expression)
	}
}
