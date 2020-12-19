package ast

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/repr"

	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Statement struct {
	generic_ast.BaseASTNode
	Empty *string `";"`
	BlockStatement *Block `| @@`
	Declaration *Declaration `| @@`
	Assignment *Assignment `| @@`
	UnaryStatement *UnaryStatement `| @@`
	Return *Return `| @@`
	If *If `| @@`
	While *While `| @@`
	For *For `| @@`
	Expression *Expression `| @@ ";"`
	ParentNode generic_ast.TraversableNode
}

func (ast *Statement) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *Statement) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *Statement) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Statement) End() lexer.Position {
	return ast.EndPos
}

func (ast *Statement) GetNode() interface{} {
	return ast
}

func (ast *Statement) IsEmpty() bool {
	return (
		ast.Empty != nil || (
		!ast.IsBlockStatement() &&
		!ast.IsDeclaration() &&
		!ast.IsAssignment() &&
		!ast.IsUnaryStatement() &&
		!ast.IsReturn() &&
		!ast.IsIf() &&
		!ast.IsWhile() &&
		!ast.IsFor() &&
		!ast.IsExpression()))
}

func (ast *Statement) IsBlockStatement() bool {
	return ast.BlockStatement != nil
}

func (ast *Statement) IsDeclaration() bool {
	return ast.Declaration != nil
}

func (ast *Statement) IsAssignment() bool {
	return ast.Assignment != nil
}

func (ast *Statement) IsUnaryStatement() bool {
	return ast.UnaryStatement != nil
}

func (ast *Statement) IsReturn() bool {
	return ast.Return != nil
}

func (ast *Statement) IsIf() bool {
	return ast.If != nil
}

func (ast *Statement) IsWhile() bool {
	return ast.While != nil
}

func (ast *Statement) IsFor() bool {
	return ast.For != nil
}

func (ast *Statement) IsExpression() bool {
	return ast.Expression != nil
}

func (ast *Statement) GetChildren() []generic_ast.TraversableNode {
	if ast.IsEmpty() {
		return []generic_ast.TraversableNode{}
		//return []generic_ast.TraversableNode{ generic_ast.MakeTraversableNodeToken(ast, *ast.Empty, ast.Pos, ast.EndPos) }
	} else if ast.IsBlockStatement() {
		return []generic_ast.TraversableNode{ ast.BlockStatement }
	} else if ast.IsDeclaration() {
		return []generic_ast.TraversableNode{ ast.Declaration }
	} else if ast.IsAssignment() {
		return []generic_ast.TraversableNode{ ast.Assignment }
	} else if ast.IsUnaryStatement() {
		return []generic_ast.TraversableNode{ ast.UnaryStatement }
	} else if ast.IsReturn() {
		return []generic_ast.TraversableNode{ ast.Return }
	} else if ast.IsIf() {
		return []generic_ast.TraversableNode{ ast.If }
	} else if ast.IsWhile() {
		return []generic_ast.TraversableNode{ ast.While }
	} else if ast.IsFor() {
		return []generic_ast.TraversableNode{ ast.For }
	} else if ast.IsExpression() {
		return []generic_ast.TraversableNode{ ast.Expression }
	}
	return []generic_ast.TraversableNode{}
}

func (ast *Statement) formatStatementInstruction(statement string, c *context.ParsingContext) string {
	if c.PrinterConfiguration.SkipStatementIdent {
		c.PrinterConfiguration.SkipStatementIdent = false
		return statement
	}
	return printNode(c, ast, "%s%s", strings.Repeat("  ", c.BlockDepth), statement)
}

func (ast *Statement) Print(c *context.ParsingContext) string {
	ret := "UNKNOWN"
	propagateSkipStatementIdent := false
	if ast.IsEmpty() {
		ret = ";"
	} else if ast.IsBlockStatement() {
		if c.PrinterConfiguration.SkipStatementIdent {
			propagateSkipStatementIdent = true
		}
		ret = ast.BlockStatement.Print(c)
	} else if ast.IsDeclaration() {
		ret = ast.Declaration.Print(c)
	} else if ast.IsAssignment() {
		ret = ast.Assignment.Print(c)
	} else if ast.IsUnaryStatement() {
		ret = ast.UnaryStatement.Print(c)
	} else if ast.IsReturn() {
		ret = ast.Return.Print(c)
	} else if ast.IsIf() {
		ret =  ast.If.Print(c)
	} else if ast.IsWhile() {
		ret = ast.While.Print(c)
	} else if ast.IsFor() {
		ret = ast.For.Print(c)
	} else if ast.IsExpression() {
		ret = printNode(c, ast, "%s;", ast.Expression.Print(c))
	}
	if propagateSkipStatementIdent {
		c.PrinterConfiguration.SkipStatementIdent = true
		ret = ast.formatStatementInstruction(ret, c)
	} else {
		ret = ast.formatStatementInstruction(ret, c)
	}
	return ret
}

//////

func (ast *Statement) Body() generic_ast.Expression {
	if ast.IsEmpty() {
		return &Empty{}
		//return hindley_milner.Batch{Exp: []generic_ast.Expression{}}
	} else if ast.IsBlockStatement() {
		return ast.BlockStatement
	} else if ast.IsDeclaration() {
		return ast.Declaration
	} else if ast.IsAssignment() {
		return ast.Assignment
	} else if ast.IsUnaryStatement() {
		return ast.UnaryStatement
	} else if ast.IsReturn() {
		return ast.Return
	} else if ast.IsIf() {
		return ast.If
	} else if ast.IsWhile() {
		return ast.While
	} else if ast.IsFor() {
		return ast.For
	} else if ast.IsExpression() {
		return ast.Expression
	}
	ast.BaseASTNode = generic_ast.BaseASTNode{}
	ast.ParentNode = nil
	fmt.Printf("Failed for node %s\n", repr.String(ast))
	panic("Invalid Statement type")
}

func feedExpressionIntoStatement(node interface{}, base generic_ast.BaseASTNode) *Statement {
	if block, ok := node.(*Block); ok {
		return &Statement{
			BaseASTNode:    base,
			BlockStatement: block,
		}
	} else if decl, ok := node.(*Declaration); ok {
		return &Statement{
			BaseASTNode:    base,
			Declaration: decl,
		}
	} else if assignment, ok := node.(*Assignment); ok {
		return &Statement{
			BaseASTNode:    base,
			Assignment: assignment,
		}
	} else if unaryStmt, ok := node.(*UnaryStatement); ok {
		return &Statement{
			BaseASTNode:    base,
			UnaryStatement: unaryStmt,
		}
	} else if returnStmt, ok := node.(*Return); ok {
		return &Statement{
			BaseASTNode:    base,
			Return: returnStmt,
		}
	} else if ifStmt, ok := node.(*If); ok {
		return &Statement{
			BaseASTNode:    base,
			If: ifStmt,
		}
	} else if whileStmt, ok := node.(*While); ok {
		return &Statement{
			BaseASTNode:    base,
			While: whileStmt,
		}
	} else if forStmt, ok := node.(*For); ok {
		return &Statement{
			BaseASTNode:    base,
			For: forStmt,
		}
	} else if expr, ok := node.(*Expression); ok {
		return &Statement{
			BaseASTNode:    base,
			Expression: expr,
		}
	}
	panic("Invalid statement type")
}

func (ast *Statement) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	mappedStatement := feedExpressionIntoStatement(mapper(ast, ast.Body(), context, false), ast.BaseASTNode)
	return mapper(parent, mappedStatement, context, true)
}

func (ast *Statement) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(ast, ast.Body(), context)
	mapper(parent, ast, context)
}

func (ast *Statement) ExpressionType() hindley_milner.ExpressionType { return hindley_milner.E_PROXY }

//

func (ast *Statement) BuildFlowGraph(builder cfg.CFGBuilder) {
	if ast.IsEmpty() {
		// Do nothing
	} else if ast.IsBlockStatement() {
		builder.BuildNode(ast.BlockStatement)
	} else if ast.IsDeclaration() {
		builder.BuildNode(ast.Declaration)
	} else if ast.IsAssignment() {
		builder.BuildNode(ast.Assignment)
	} else if ast.IsUnaryStatement() {
		builder.BuildNode(ast.UnaryStatement)
	} else if ast.IsReturn() {
		builder.BuildNode(ast.Return)
	} else if ast.IsIf() {
		builder.BuildNode(ast.If)
	} else if ast.IsWhile() {
		builder.BuildNode(ast.While)
	} else if ast.IsFor() {
		builder.BuildNode(ast.For)
	} else if ast.IsExpression() {
		builder.BuildNode(ast.Expression)
	}
}
