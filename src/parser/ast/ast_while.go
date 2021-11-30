package ast

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type While struct {
	generic_ast.BaseASTNode
	Condition  *Expression `"while" "(" @@ ")"`
	Do         *Statement  `@@`
	ParentNode generic_ast.TraversableNode
}

func (ast *While) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *While) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *While) Begin() lexer.Position {
	return ast.Pos
}

func (ast *While) End() lexer.Position {
	return ast.EndPos
}

func (ast *While) GetNode() interface{} {
	return ast
}

func (ast *While) Print(c *context.ParsingContext) string {
	c.PrinterConfiguration.SkipStatementIdent = true
	body := ast.Do.Print(c)
	return printNode(c, ast, "while (%s) %s", ast.Condition.Print(c), body)
}

func (ast *While) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{
		ast.Condition,
		ast.Do,
	}
}

func (ast *While) Validate(c *context.ParsingContext) generic_ast.NodeError {
	if ast.Do != nil {
		if ast.Do.IsDeclaration() {
			message := fmt.Sprintf("Declaration as a non-block expression inside while statement is forbidden. Please use { } brackets and create the definition there.")
			return generic_ast.NewNodeError(
				"Declaration not allowed",
				ast,
				message,
				message)
		}
	}
	return nil
}

///

func (ast *While) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	return mapper(parent, &While{
		BaseASTNode: ast.BaseASTNode,
		Condition:   mapper(ast, ast.Condition, context, false).(*Expression),
		Do:          mapper(ast, ast.Do, context, false).(*Statement),
		ParentNode:  parent.(generic_ast.TraversableNode),
	}, context, true)
}

func (ast *While) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(ast, ast.Condition, context)
	mapper(ast, ast.Do, context)
	mapper(parent, ast, context)
}

func (ast *While) Fn(c hindley_milner.InferContext) generic_ast.Expression {
	return &hindley_milner.EmbeddedTypeExpr{GetType: func() *hindley_milner.Scheme {
		return hindley_milner.Concreate(
			hindley_milner.NewFnType(CreatePrimitive(T_BOOL), CreatePrimitive(T_VOID), CreatePrimitive(T_VOID)))
	}, Source: ast}
}

func (ast *While) Body() generic_ast.Expression {
	args := []generic_ast.Expression{
		ast.Condition,
		ast.Do,
	}
	return hindley_milner.Batch{
		Exp: args,
	}
}

func (ast *While) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_APPLICATION
}

///

func (ast *While) BuildFlowGraph(builder cfg.CFGBuilder) {
	// flows as such (range same w/o init & post):
	// previous -> [ init -> ] for -> body -> [ post -> ] for -> next

	skipJumpOut := false
	skipBody := false
	if constCond, hasConstCond := ast.Condition.ExtractConst(); hasConstCond {
		if *(constCond.(*Primary).Bool) {
			skipJumpOut = true
		} else {
			skipBody = true
		}
	}

	if skipBody {
		return
	}

	var post generic_ast.NormalNode = ast

	builder.AddBlockSuccesor(ast)

	builder.UpdatePrev([]generic_ast.NormalNode{ast})
	builder.BuildNode(ast.Do)

	builder.AddBlockSuccesor(post)

	ctrlExits := []generic_ast.NormalNode{ast}

	// handle any branches; if no label or for me: handle and remove from branches.
	for i := 0; i < len(builder.Branches()); i++ {
		//br := builder.Branches()[i]
		// Deal with continue/break here if such thing will be implemented
	}

	if skipJumpOut {
		builder.UpdatePrev([]generic_ast.NormalNode{})
	} else {
		builder.UpdatePrev(ctrlExits) // for stmt and any appropriate break statements
	}
}

func (ast *While) GetUsedVariables(vars cfg.VariableSet, visitedMap map[generic_ast.TraversableNode]struct{}) cfg.VariableSet {
	return cfg.GetAllVariables(ast.Condition, visitedMap)
}
