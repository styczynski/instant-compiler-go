package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type LogicalOperation struct {
	generic_ast.BaseASTNode
	Equality     *Equality         `@@`
	Op           string            `[ @( "|" "|" | "&" "&" )`
	Next         *LogicalOperation `  @@ ]`
	ParentNode   generic_ast.TraversableNode
	ResolvedType hindley_milner.Type
}

func (ast *LogicalOperation) OnTypeReturned(t hindley_milner.Type) {
	ast.ResolvedType = t
}

func (ast *LogicalOperation) ExtractConst() (generic_ast.TraversableNode, bool) {
	if ast.HasNext() {
		return nil, false
	}
	return ast.Equality.ExtractConst()
}

func (ast *LogicalOperation) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *LogicalOperation) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *LogicalOperation) Begin() lexer.Position {
	return ast.Pos
}

func (ast *LogicalOperation) End() lexer.Position {
	return ast.EndPos
}

func (ast *LogicalOperation) GetNode() interface{} {
	return ast
}

func (ast *LogicalOperation) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{
		ast.Equality,
		generic_ast.MakeTraversableNodeToken(ast, ast.Op, ast.Pos, ast.EndPos),
		ast.Next,
	}
}

func (ast *LogicalOperation) HasNext() bool {
	return ast.Next != nil
}

func (ast *LogicalOperation) Print(c *context.ParsingContext) string {
	if ast.HasNext() {
		return printBinaryOperation(c, ast, ast.Equality.Print(c), ast.Op, ast.Next.Print(c))
	}
	return ast.Equality.Print(c)
}

/////

func (ast *LogicalOperation) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	next := ast.Next
	if ast.HasNext() {
		next = mapper(ast, ast.Next, context, false).(*LogicalOperation)
	}
	return mapper(parent, &LogicalOperation{
		BaseASTNode: ast.BaseASTNode,
		Equality:    mapper(ast, ast.Equality, context, false).(*Equality),
		Op:          ast.Op,
		Next:        next,
		ParentNode:  parent.(generic_ast.TraversableNode),
	}, context, true)
}

func (ast *LogicalOperation) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(ast, ast.Equality, context)
	if ast.HasNext() {
		mapper(ast, ast.Next, context)
	}
	mapper(parent, ast, context)
}

func (ast *LogicalOperation) Fn(c hindley_milner.InferContext) generic_ast.Expression {
	return &BuiltinFunction{
		BaseASTNode: ast.BaseASTNode,
		name:        ast.Op,
	}
}

func (ast *LogicalOperation) Body() generic_ast.Expression {
	if !ast.HasNext() {
		return ast.Equality
	}
	return hindley_milner.Batch{
		Exp: []generic_ast.Expression{
			ast.Equality,
			ast.Next,
		},
	}
}

func (ast *LogicalOperation) ExpressionType() hindley_milner.ExpressionType {
	if !ast.HasNext() {
		return hindley_milner.E_PROXY
	}
	return hindley_milner.E_APPLICATION
}

///

func (ast *LogicalOperation) ConstFold() generic_ast.TraversableNode {
	//fmt.Printf("Optimize cosntFold() LogicalOperation\n")
	if ast.HasNext() {
		const1, ok1 := ast.Equality.ExtractConst()
		if ok1 && ast.Op == "||" && const1.(*Primary).IsBool() {
			//fmt.Printf("Fold left ||\n")
			if *const1.(*Primary).Bool {
				v := true
				ast.Equality.Comparison.Addition.Multiplication.Unary.UnaryApplication.Index.Primary = &Primary{
					BaseASTNode: ast.BaseASTNode,
					Bool:        &v,
				}
				ast.Next = nil
				return ast
			}
		}
		if ok1 && ast.Op == "&&" && const1.(*Primary).IsBool() {
			//fmt.Printf("Fold left &&\n")
			if !(*const1.(*Primary).Bool) {
				v := false
				ast.Equality.Comparison.Addition.Multiplication.Unary.UnaryApplication.Index.Primary = &Primary{
					BaseASTNode: ast.BaseASTNode,
					Bool:        &v,
				}
				ast.Next = nil
				return ast
			}
		}

		const2, ok2 := ast.Next.Equality.ExtractConst()
		if ok1 && ok2 {
			//fmt.Printf("Fold all\n")
			p1 := const1.(*Primary)
			p2 := const2.(*Primary)
			v := p1.And(p2, ast.Op)
			// Change pointers
			ast.Equality.Comparison.Addition.Multiplication.Unary.UnaryApplication.Index.Primary = v
			ast.Op = ast.Next.Op
			ast.Next = ast.Next.Next
			return ast
		}
	}
	//fmt.Printf("Fold nothing\n")
	return ast
}
