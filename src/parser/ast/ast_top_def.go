package ast

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type TopDef struct {
	BaseASTNode
	Function *FnDef `@@`
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

func (ast *TopDef) GetChildren() []TraversableNode {
	return []TraversableNode{
		ast.Function,
	}
}

func (ast *TopDef) Print(c *context.ParsingContext) string {
	return printNode(c, ast, "%s", ast.Function.Print(c))
}

func (ast *TopDef) IsFunction() bool {
	return ast.Function != nil
}


///


func (ast *TopDef) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_OPAQUE_BLOCK
}

func (ast *TopDef) Map(mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	if ast.IsFunction() {
		return mapper(&TopDef{
			BaseASTNode: ast.BaseASTNode,
			Function:    mapper(ast.Function).(*FnDef),
		})
	} else {
		panic("Invalid TopDef type.")
	}
}

func (ast *TopDef) Visit(mapper hindley_milner.ExpressionMapper) {
	if ast.IsFunction() {
		mapper(ast.Function)
	} else {
		panic("Invalid TopDef type.")
	}
	mapper(ast)
}

func (ast *TopDef) GetContents() hindley_milner.Batch {
	if ast.IsFunction() {
		return hindley_milner.Batch{Exp: ast.Expressions()}
	} else {
		panic("Invalid TopDef type.")
	}
}

func (ast *TopDef) Expressions() []hindley_milner.Expression {
	if ast.IsFunction() {
		return []hindley_milner.Expression{ ast.Function, }
	} else {
		panic("Invalid TopDef type.")
	}
}

func (ast *TopDef) Body() hindley_milner.Expression {
	panic(fmt.Errorf("Batch Body() method cannot be called."))
}
