package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Equality struct {
	generic_ast.BaseASTNode
	Comparison *Comparison `@@`
	Op         string      `[ @( "!" "=" | "=" "=" )`
	Next       *Equality   `  @@ ]`
	ParentNode generic_ast.TraversableNode
}

func (ast *Equality) ExtractConst() (generic_ast.TraversableNode, bool) {
	if ast.HasNext() {
		return nil, false
	}
	return ast.Comparison.ExtractConst()
}

func (ast *Equality) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *Equality) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *Equality) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Equality) End() lexer.Position {
	return ast.EndPos
}

func (ast *Equality) GetNode() interface{} {
	return ast
}

func (ast *Equality) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{
		ast.Comparison,
		generic_ast.MakeTraversableNodeToken(ast, ast.Op, ast.Pos, ast.EndPos),
		ast.Next,
	}
}

func (ast *Equality) HasNext() bool {
	return ast.Next != nil
}

func (ast *Equality) Print(c *context.ParsingContext) string {
	if ast.HasNext() {
		return printBinaryOperation(c, ast, ast.Comparison.Print(c), ast.Op, ast.Next.Print(c))
	}
	return ast.Comparison.Print(c)
}

///

/////

func (ast *Equality) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	next := ast.Next
	if ast.HasNext() {
		next = mapper(ast, ast.Next, context, false).(*Equality)
	}
	return mapper(parent, &Equality{
		BaseASTNode: ast.BaseASTNode,
		Comparison:  mapper(ast, ast.Comparison, context, false).(*Comparison),
		Op:          ast.Op,
		Next:        next,
		ParentNode:  parent.(generic_ast.TraversableNode),
	}, context, true)
}

func (ast *Equality) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(ast, ast.Comparison, context)
	if ast.HasNext() {
		mapper(ast, ast.Next, context)
	}
	mapper(parent, ast, context)
}

func (ast *Equality) Fn(c hindley_milner.InferContext) generic_ast.Expression {
	return &BuiltinFunction{
		BaseASTNode: ast.BaseASTNode,
		name:        ast.Op,
	}
}

func (ast *Equality) Body() generic_ast.Expression {
	if !ast.HasNext() {
		return ast.Comparison
	}
	return hindley_milner.Batch{
		Exp: []generic_ast.Expression{
			ast.Comparison,
			ast.Next,
		},
	}
}

func (ast *Equality) ExpressionType() hindley_milner.ExpressionType {
	if !ast.HasNext() {
		return hindley_milner.E_PROXY
	}
	return hindley_milner.E_APPLICATION
}

//func (ast *Equality) Name() hindley_milner.NameGroup     { return hindley_milner.Name("bool") }
//func (ast *Equality) Body() generic_ast.Expression { return ast }
//func (ast *Equality) Map(mapper generic_ast.ExpressionVisitor) generic_ast.Expression {
//	return mapper(ast)
//}
//func (ast *Equality) Visit(mapper generic_ast.ExpressionVisitor) {
//	mapper(ast)
//}
//func (ast *Equality) Type() hindley_milner.Type {
//	return CreatePrimitive(T_BOOL)
//}
//// TODO: Lit/lambda needed?
//func (ast *Equality) ExpressionType() hindley_milner.ExpressionType { return hindley_milner.E_LITERAL }
//

//

func (ast *Equality) ConstFold() generic_ast.TraversableNode {
	if ast.HasNext() {
		const1, ok1 := ast.Comparison.ExtractConst()
		const2, ok2 := ast.Next.Comparison.ExtractConst()
		if ok1 && ok2 {
			p1 := const1.(*Primary)
			p2 := const2.(*Primary)
			v := p1.Compare(p2, ast.Op)
			// Change pointers
			ast.Comparison.Addition.Multiplication.Unary.UnaryApplication.Index.Primary = v
			ast.Op = ast.Next.Op
			ast.Next = ast.Next.Next
			return ast
		}
	}
	return ast
}
