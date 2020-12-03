package ast

import (
	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/parser/context"
)

type DeclarationItem struct {
	BaseASTNode
	Name string `@Ident`
	Initializer *Expression `( "=" @@ )?`
}

func (ast *DeclarationItem) Begin() lexer.Position {
	return ast.Pos
}

func (ast *DeclarationItem) End() lexer.Position {
	return ast.EndPos
}

func (ast *DeclarationItem) GetNode() interface{} {
	return ast
}

func (ast *DeclarationItem) GetChildren() []TraversableNode {
	return []TraversableNode{
		MakeTraversableNodeToken(ast.Name, ast.Pos, ast.EndPos),
		ast.Initializer,
	}
}

func (ast *DeclarationItem) HasInitializer() bool {
	return ast.Initializer != nil
}

func (ast *DeclarationItem) Print(c *context.ParsingContext) string {
	if ast.HasInitializer() {
		return printNode(c, ast, "%s = %s", ast.Name, ast.Initializer.Print(c))
	}
	return printNode(c, ast, "%s", ast.Name)
}
