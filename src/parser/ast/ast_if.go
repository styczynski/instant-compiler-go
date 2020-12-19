package ast

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type If struct {
	generic_ast.BaseASTNode
	Condition *Expression `"if" "(" @@ ")"`
	Then *Statement `@@`
	Else *Statement `( "else" @@ )?`
	ParentNode generic_ast.TraversableNode
}

func (ast *If) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *If) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *If) Begin() lexer.Position {
	return ast.Pos
}

func (ast *If) End() lexer.Position {
	return ast.EndPos
}

func (ast *If) GetNode() interface{} {
	return ast
}

func (ast *If) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{
		ast.Condition,
		ast.Then,
		ast.Else,
	}
}

func (ast *If) HasElseBlock() bool {
	return ast.Else != nil
}

func (ast *If) Print(c *context.ParsingContext) string {
	if ast.HasElseBlock(){
		return printNode(c, ast, "if (%s) %s else %s", ast.Condition.Print(c), makeBlockFromStatement(ast.Then).Print(c), makeBlockFromStatement(ast.Else).Print(c))
	}
	return printNode(c, ast, "if (%s) %s", ast.Condition.Print(c), ast.Then.Print(c))
}

///

func (ast *If) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	if ast.HasElseBlock() {
		return mapper(parent, &If{
			BaseASTNode: ast.BaseASTNode,
			Condition:   mapper(ast, ast.Condition, context, false).(*Expression),
			Then:        mapper(ast, ast.Then, context, false).(*Statement),
			Else:        mapper(ast, ast.Else, context, false).(*Statement),
			ParentNode:  parent.(generic_ast.TraversableNode),
		}, context, true)
	} else {
		return mapper(parent, &If{
			BaseASTNode: ast.BaseASTNode,
			Condition:   mapper(ast, ast.Condition, context, false).(*Expression),
			Then:        mapper(ast, ast.Then, context, false).(*Statement),
			ParentNode:  parent.(generic_ast.TraversableNode),
		}, context, true)
	}
}

func (ast *If) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(ast, ast.Condition, context)
	mapper(ast, ast.Then, context)
	if ast.HasElseBlock() {
		mapper(ast, ast.Else, context)
	}
	mapper(parent, ast, context)
}

func (ast *If) Fn() generic_ast.Expression {
	//return &BuiltinFunction{
	//	BaseASTNode: ast.BaseASTNode,
	//	name: "if",
	//}
	return &hindley_milner.EmbeddedTypeExpr{GetType: func() *hindley_milner.Scheme {
		return hindley_milner.NewScheme(
			hindley_milner.TypeVarSet{hindley_milner.TVar('a'), hindley_milner.TVar('b')},
			hindley_milner.NewFnType(CreatePrimitive(T_BOOL), hindley_milner.TVar('a'), hindley_milner.TVar('b'), CreatePrimitive(T_VOID)))
	}, Source: ast,}
}

func (ast *If) Body() generic_ast.Expression {
	args := []generic_ast.Expression{
		ast.Condition,
		ast.Then,
	}
	if ast.HasElseBlock() {
		args = append(args, ast.Else)
	} else {
		args = append(args, ast.Then)
	}
	return hindley_milner.Batch{
		Exp: args,
	}
}

func (ast *If) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_APPLICATION
}

func (ast *If) Validate(c *context.ParsingContext) generic_ast.NodeError {
	if ast.Then != nil {
		if ast.Then.IsDeclaration() {
			message := fmt.Sprintf("Declaration as a non-block expression inside if statement is forbidden. Please use { } brackets and create the definition there.")
			return generic_ast.NewNodeError(
				"Declaration not allowed",
				ast,
				message,
				message)
		}
	}
	if ast.Else != nil {
		if ast.Else.IsDeclaration() {
			message := fmt.Sprintf("Declaration as a non-block expression inside if statement else block is forbidden. Please use { } brackets and create the definition there.")
			return generic_ast.NewNodeError(
				"Declaration not allowed",
				ast,
				message,
				message)
		}
	}
	return nil
}

//

func (ast *If) BuildFlowGraph(builder cfg.CFGBuilder) {

	skipThen := false
	skipElse := false
	if constCond, hasConstCond := ast.Condition.ExtractConst(); hasConstCond {
		if *(constCond.(*Primary).Bool) {
			skipElse = true
		} else {
			skipThen = true
		}
	}

	if skipThen || skipElse {
		if skipThen {
			builder.BuildNode(ast.Else)
		} else if skipElse {
			builder.BuildNode(ast.Then)
		}
		return
	}

	builder.AddBlockSuccesor(ast)

	builder.UpdatePrev([]generic_ast.NormalNode{ast})
	builder.BuildNode(ast.Then)

	ctrlExits := builder.GetPrev() // aggregate of builder.prev from each condition
	if ast.HasElseBlock() {
		builder.UpdatePrev([]generic_ast.NormalNode{ast})
		builder.BuildNode(ast.Else)
		ctrlExits = append(ctrlExits, builder.GetPrev()...)
	} else {
		ctrlExits = append(ctrlExits, ast)
	}
	builder.UpdatePrev(ctrlExits)
}

func (ast *If) GetUsedVariables(vars cfg.VariableSet, visitedMap map[generic_ast.TraversableNode]struct{}) cfg.VariableSet {
	return cfg.GetAllUsagesVariables(ast.Condition, visitedMap)
}