package ast

import (
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
	return mapper(parent, &If{
		BaseASTNode: ast.BaseASTNode,
		Condition: mapper(ast, ast.Condition, context).(*Expression),
		Then: mapper(ast, ast.Then, context).(*Statement),
		Else: mapper(ast, ast.Else, context).(*Statement),
		ParentNode: parent.(generic_ast.TraversableNode),
	}, context)
}

func (ast *If) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) {
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
	}}
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

//

func (ast *If) BuildFlowGraph(builder cfg.CFGBuilder) {
	builder.AddSucc(ast)

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