package ast

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/errors"
	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type TopDef struct {
	generic_ast.BaseASTNode
	Class      *Class `@@`
	Function   *FnDef `| @@`
	ParentNode generic_ast.TraversableNode
}

func (ast *TopDef) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *TopDef) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *TopDef) GetDefinedIdentifiers(c hindley_milner.InferContext, pre bool) *hindley_milner.NameGroup {
	if ast.IsFunction() {
		return hindley_milner.NameWithType(ast.Function.Name, ast.Function.GetDeclarationType())
	} else if ast.IsClass() {
		return ast.Class.GetDeclarationIdentifiers()
	}
	return hindley_milner.EmptyNameGroup()
}

func (ast *TopDef) Begin() lexer.Position {
	return ast.Pos
}

func (ast *TopDef) End() lexer.Position {
	return ast.EndPos
}

func (ast *TopDef) GetNode() interface{} {
	return ast
}

func (ast *TopDef) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{
		ast.Function,
	}
}

func (ast *TopDef) Print(c *context.ParsingContext) string {
	if ast.IsFunction() {
		return printNode(c, ast, "%s", ast.Function.Print(c))
	} else if ast.IsClass() {
		return printNode(c, ast, "%s", ast.Class.Print(c))
	}
	panic("Invalid TopDef type.")
}

func (ast *TopDef) IsClass() bool {
	return ast.Class != nil
}

func (ast *TopDef) IsFunction() bool {
	return ast.Function != nil
}

///

func (ast *TopDef) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_OPAQUE_BLOCK
}

func (ast *TopDef) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	if ast.IsFunction() {
		return mapper(parent, &TopDef{
			BaseASTNode: ast.BaseASTNode,
			Function:    mapper(ast, ast.Function, context, false).(*FnDef),
			ParentNode:  parent.(generic_ast.TraversableNode),
		}, context, true)
	} else if ast.IsClass() {
		return mapper(parent, &TopDef{
			BaseASTNode: ast.BaseASTNode,
			Class:       mapper(ast, ast.Class, context, false).(*Class),
			ParentNode:  parent.(generic_ast.TraversableNode),
		}, context, true)
	} else {
		panic("Invalid TopDef type.")
	}
}

func (ast *TopDef) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	if ast.IsFunction() {
		mapper(ast, ast.Function, context)
	} else if ast.IsClass() {
		mapper(ast, ast.Class, context)
	} else {
		panic("Invalid TopDef type.")
	}
	mapper(parent, ast, context)
}

func (ast *TopDef) GetContents() hindley_milner.Batch {
	if ast.IsFunction() {
		return hindley_milner.Batch{Exp: ast.Expressions()}
	} else if ast.IsClass() {
		return hindley_milner.Batch{Exp: ast.Expressions()}
	} else {
		panic("Invalid TopDef type.")
	}
}

func (ast *TopDef) Expressions() []generic_ast.Expression {
	if ast.IsFunction() {
		return []generic_ast.Expression{ast.Function}
	} else if ast.IsClass() {
		return []generic_ast.Expression{ast.Class}
	} else {
		panic("Invalid TopDef type.")
	}
}

func (ast *TopDef) Body() generic_ast.Expression {
	panic(fmt.Errorf("Batch Body() method cannot be called."))
}

//

func (ast *TopDef) BuildFlowGraph(builder cfg.CFGBuilder) {
	if ast.IsFunction() {
		builder.BuildNode(ast.Function)
	}
}

func (ast *TopDef) OnFlowAnalysis(flow cfg.FlowAnalysis) error {
	if ast.IsFunction() {
		// Validate flow graph
		retType, _ := ast.Function.ReturnType.GetType(nil).Type()
		if !retType.Eq(CreatePrimitive(T_VOID)) {

			gateways := flow.Graph().GetAllEndGateways()
			if len(flow.Graph().Blocks()) <= 3 {
				return errors.CreateLocalizedError("Missing return", "Non-void function is missing return.", ast)
			}
			for _, gateway := range gateways {
				if returnNode, isRet := gateway.(*Return); isRet {
					if !returnNode.HasExpression() {
						// Fail
						return errors.CreateLocalizedError("Invalid return", "Non-void function cannot use return statement without expression.", returnNode)
					}
				} else {
					// Fail
					return errors.CreateLocalizedError("Missing return", "Non-void function should have return in each possible branch.", gateway)
				}
			}
		}
	}
	return nil
}

/*
graph := flow.Graph()
	output := []generic_ast.NormalNode{}
	for _, stmt := range flow.input {
		block := graph.blocks[stmt]
		if block != nil {
			output = append(output, stmt)
		}
	}
	return output
*/

func (ast *TopDef) AfterFlowAnalysis(flow cfg.FlowAnalysis) {
	// if ast.IsFunction() {
	// 	output := []*Statement{}
	// 	graph := flow.Graph()
	// 	for _, stmt := range ast.Function.FunctionBody.Statements {
	// 		if graph.Exists(stmt.GetChildren()[0].(generic_ast.NormalNode)) {
	// 			output = append(output, stmt)
	// 		}
	// 	}
	// 	ast.Function.FunctionBody.Statements = output
	// }
}
