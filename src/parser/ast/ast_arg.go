package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type Arg struct {
	generic_ast.BaseASTNode
	ArgumentType Type   `@@`
	Name         string `@Ident`
	ParentNode   generic_ast.TraversableNode
	ResolvedType hindley_milner.Type
}

func (ast *Arg) OnTypeReturned(t hindley_milner.Type) {
	ast.ResolvedType = t
}

func (ast *Arg) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *Arg) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *Arg) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Arg) End() lexer.Position {
	return ast.EndPos
}

func (ast *Arg) GetNode() interface{} {
	return ast
}

func (ast *Arg) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{
		&ast.ArgumentType,
		generic_ast.MakeTraversableNodeToken(ast, ast.Name, ast.Pos, ast.EndPos),
	}
}

func (ast *Arg) Print(c *context.ParsingContext) string {
	return printNode(c, ast, "%s %s", ast.ArgumentType.Print(c), ast.Name)
}
