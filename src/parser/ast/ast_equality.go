package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Equality struct {
	BaseASTNode
	Comparison *Comparison `@@`
	Op         string      `[ @( "!" "=" | "=" "=" )`
	Next       *Equality   `  @@ ]`
	ParentNode TraversableNode
}

func (ast *Equality) Parent() TraversableNode {
	return ast.ParentNode
}

func (ast *Equality) OverrideParent(node TraversableNode) {
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

func (ast *Equality) GetChildren() []TraversableNode {
	return []TraversableNode{
		ast.Comparison,
		MakeTraversableNodeToken(ast, ast.Op, ast.Pos, ast.EndPos),
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

func (ast *Equality) Map(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	next := ast.Next
	if ast.HasNext() {
		next = mapper(ast, ast.Next).(*Equality)
	}
	return mapper(parent, &Equality{
		BaseASTNode: ast.BaseASTNode,
		Comparison:    mapper(ast, ast.Comparison).(*Comparison),
		Op:          ast.Op,
		Next:        next,
		ParentNode: parent.(TraversableNode),
	})
}

func (ast *Equality) Visit(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) {
	mapper(ast, ast.Comparison)
	if ast.HasNext() {
		mapper(ast, ast.Next)
	}
	mapper(parent, ast)
}

func (ast *Equality) Fn() hindley_milner.Expression {
	return &BuiltinFunction{
		BaseASTNode: ast.BaseASTNode,
		name: ast.Op,
	}
}

func (ast *Equality) Body() hindley_milner.Expression {
	if !ast.HasNext() {
		return ast.Comparison
	}
	return hindley_milner.Batch{
		Exp: []hindley_milner.Expression{
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
//func (ast *Equality) Body() hindley_milner.Expression { return ast }
//func (ast *Equality) Map(mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
//	return mapper(ast)
//}
//func (ast *Equality) Visit(mapper hindley_milner.ExpressionMapper) {
//	mapper(ast)
//}
//func (ast *Equality) Type() hindley_milner.Type {
//	return CreatePrimitive(T_BOOL)
//}
//// TODO: Lit/lambda needed?
//func (ast *Equality) ExpressionType() hindley_milner.ExpressionType { return hindley_milner.E_LITERAL }
//
