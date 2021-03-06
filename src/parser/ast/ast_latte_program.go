package ast

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type LatteProgram struct {
	generic_ast.BaseASTNode
	Definitions []*TopDef `@@*`
	ParentNode  generic_ast.TraversableNode
}

func (ast *LatteProgram) Parent() generic_ast.TraversableNode {
	return nil
}

func (ast *LatteProgram) OverrideParent(node generic_ast.TraversableNode) {
	// No-op
}

func (ast *LatteProgram) GetIdentifierDeps(c hindley_milner.InferContext, pre bool) hindley_milner.NameGroup {
	idents := []string{}
	nameMapping := map[string]*hindley_milner.Scheme{}
	for _, def := range ast.Definitions {
		names, types := def.GetDefinedIdentifier(c, pre)
		for i, name := range names {
			idents = append(idents, name)
			nameMapping[name] = types[i]
		}
	}
	return hindley_milner.NamesWithTypes(idents, nameMapping)
}

func (ast *LatteProgram) Begin() lexer.Position {
	return ast.Pos
}

func (ast *LatteProgram) End() lexer.Position {
	return ast.EndPos
}

func (ast *LatteProgram) GetNode() interface{} {
	return ast
}

func (ast *LatteProgram) GetChildren() []generic_ast.TraversableNode {
	nodes := []generic_ast.TraversableNode{}
	for _, child := range ast.Definitions {
		nodes = append(nodes, child)
	}
	return nodes
}

func (ast *LatteProgram) Print(c *context.ParsingContext) string {
	defs := []string{}
	for _, def := range ast.Definitions {
		defs = append(defs, def.Print(c))
	}
	return printNode(c, ast, "%s\n", strings.Join(defs, "\n\n"))
}

func (ast *LatteProgram) Body() generic_ast.Expression {
	panic(fmt.Errorf("Batch Body() method cannot be called."))
}

func (ast *LatteProgram) Validate(c *context.ParsingContext) generic_ast.NodeError {
	//fmt.Printf("SUKA BLYAT!\n")
	hasMain := false
	for _, def := range ast.Definitions {
		if def.IsFunction() {
			if def.Function.Name == "main" {
				hasMain = true
			}
		}
	}
	if !hasMain {
		message := fmt.Sprintf("main() function is missing. Please create a top-level function with signature int main() { ... }")
		return generic_ast.NewNodeError(
			"Missing main()",
			ast,
			message,
			message)
	}
	return nil
}

/////

func (ast *LatteProgram) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	mappedDef := []*TopDef{}
	for _, def := range ast.Definitions {
		mappedDef = append(mappedDef, mapper(ast, def, context, false).(*TopDef))
	}
	return mapper(parent, &LatteProgram{
		BaseASTNode: ast.BaseASTNode,
		Definitions: mappedDef,
		ParentNode:  parent.(generic_ast.TraversableNode),
	}, context, true).(*LatteProgram)
}

func (ast *LatteProgram) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	for _, def := range ast.Definitions {
		mapper(ast, def, context)
	}
	mapper(parent, ast, context)
}

func (ast *LatteProgram) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_OPAQUE_BLOCK
}

func (ast *LatteProgram) GetContents() hindley_milner.Batch {
	exp := []generic_ast.Expression{}
	for _, def := range ast.Definitions {
		exp = append(exp, def)
	}
	return hindley_milner.Batch{
		Exp: exp,
	}
}

//

func (ast *LatteProgram) BuildFlowGraph(builder cfg.CFGBuilder) {
	builder.BuildBlock(ast)
}
