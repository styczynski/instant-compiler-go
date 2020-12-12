package ast

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type TopDef struct {
	BaseASTNode
	Class *Class `@@`
	Function *FnDef `| @@`
	ParentNode TraversableNode
}

func (ast *TopDef) Parent() TraversableNode {
	return ast.ParentNode
}

func (ast *TopDef) OverrideParent(node TraversableNode) {
	ast.ParentNode = node
}

func (ast *TopDef) GetDefinedIdentifier() []string {
	if ast.IsFunction() {
		return []string{
			ast.Function.Name,
		}
	} else if ast.IsClass() {
		return []string{
			ast.Class.Name,
		}
	}
	return []string{}
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

func (ast *TopDef) Map(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) hindley_milner.Expression {
	if ast.IsFunction() {
		return mapper(parent, &TopDef{
			BaseASTNode: ast.BaseASTNode,
			Function:    mapper(ast, ast.Function).(*FnDef),
			ParentNode: parent.(TraversableNode),
		})
	} else if ast.IsClass() {
		return mapper(parent, &TopDef{
			BaseASTNode: ast.BaseASTNode,
			Class:    mapper(ast, ast.Class).(*Class),
			ParentNode: parent.(TraversableNode),
		})
	} else {
		panic("Invalid TopDef type.")
	}
}

func (ast *TopDef) Visit(parent hindley_milner.Expression, mapper hindley_milner.ExpressionMapper) {
	if ast.IsFunction() {
		mapper(ast, ast.Function)
	} else if ast.IsClass() {
		mapper(ast, ast.Class)
	} else {
		panic("Invalid TopDef type.")
	}
	mapper(parent, ast)
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

func (ast *TopDef) Expressions() []hindley_milner.Expression {
	if ast.IsFunction() {
		return []hindley_milner.Expression{ ast.Function, }
	} else if ast.IsClass() {
		return []hindley_milner.Expression{ ast.Class, }
	} else {
		panic("Invalid TopDef type.")
	}
}

func (ast *TopDef) Body() hindley_milner.Expression {
	panic(fmt.Errorf("Batch Body() method cannot be called."))
}
