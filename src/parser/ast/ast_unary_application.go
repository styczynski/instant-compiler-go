package ast

import (
	"strings"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type UnaryApplication struct {
	generic_ast.BaseASTNode
	Index        *Index        `@@`
	AppToken     string        `( @"("`
	Arguments    []*Expression `(@@ ("," @@)*)? ")" )?`
	ParentNode   generic_ast.TraversableNode
	ResolvedType hindley_milner.Type
}

func (ast *UnaryApplication) OnTypeReturned(t hindley_milner.Type) {
	ast.ResolvedType = t
}

func (ast *UnaryApplication) ExtractConst() (generic_ast.TraversableNode, bool) {
	if ast.IsIndex() {
		return ast.Index.ExtractConst()
	}
	return nil, false
}

func (ast *UnaryApplication) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *UnaryApplication) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *UnaryApplication) Begin() lexer.Position {
	return ast.Pos
}

func (ast *UnaryApplication) End() lexer.Position {
	return ast.EndPos
}

func (ast *UnaryApplication) GetNode() interface{} {
	return ast
}

func (ast *UnaryApplication) GetChildren() []generic_ast.TraversableNode {
	if ast.IsApplication() {
		nodes := make([]generic_ast.TraversableNode, len(ast.Arguments)+1)
		nodes = append(nodes, ast.Index)
		for _, child := range ast.Arguments {
			nodes = append(nodes, child)
		}
		return nodes
	} else if ast.IsIndex() {
		return []generic_ast.TraversableNode{
			ast.Index,
		}
	}
	return []generic_ast.TraversableNode{}
}

func (ast *UnaryApplication) IsApplication() bool {
	return len(ast.AppToken) > 0
}

func (ast *UnaryApplication) IsIndex() bool {
	return len(ast.AppToken) == 0
}

func (ast *UnaryApplication) Print(c *context.ParsingContext) string {
	if ast.IsApplication() {
		args := []string{}
		for _, argument := range ast.Arguments {
			args = append(args, argument.Print(c))
		}
		return printNode(c, ast, "%s(%s)", ast.Index.Print(c), strings.Join(args, ", "))
	} else if ast.IsIndex() {
		return ast.Index.Print(c)
	}
	panic("Unvalid UnaryApplication value")
}

////

func (ast *UnaryApplication) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	if ast.IsApplication() {
		args := []*Expression{}
		for _, arg := range ast.Arguments {
			args = append(args, mapper(ast, arg, context, false).(*Expression))
		}
		return mapper(parent, &UnaryApplication{
			BaseASTNode: ast.BaseASTNode,
			Index:       mapper(ast, ast.Index, context, false).(*Index),
			Arguments:   args,
			ParentNode:  parent.(generic_ast.TraversableNode),
		}, context, true)
	} else if ast.IsIndex() {
		return mapper(parent, &UnaryApplication{
			BaseASTNode: ast.BaseASTNode,
			Index:       mapper(ast, ast.Index, context, false).(*Index),
			ParentNode:  parent.(generic_ast.TraversableNode),
		}, context, true)
	}
	panic("Invalid UnaryApplication operation type")
}

func (ast *UnaryApplication) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(ast, ast.Index, context)
	if ast.IsIndex() {
		// Do nothing
	} else if ast.IsApplication() {
		for _, arg := range ast.Arguments {
			mapper(ast, arg, context)
		}
	}
	mapper(parent, ast, context)
}

func (ast *UnaryApplication) Fn(c hindley_milner.InferContext) generic_ast.Expression {
	return ast.Index
}

func (ast *UnaryApplication) Body() generic_ast.Expression {
	if ast.IsIndex() {
		return ast.Index
	}
	if len(ast.Arguments) == 0 {
		return hindley_milner.Batch{
			Exp: []generic_ast.Expression{
				hindley_milner.EmbeddedTypeExpr{
					GetType: func() *hindley_milner.Scheme {
						return hindley_milner.NewScheme(nil, CreatePrimitive(T_VOID_ARG))
					},
					Source: ast,
				},
			},
		}
	}
	args := []generic_ast.Expression{}
	for _, arg := range ast.Arguments {
		args = append(args, arg)
	}
	return hindley_milner.Batch{
		Exp: args,
	}
}

func (ast *UnaryApplication) ExpressionType() hindley_milner.ExpressionType {
	if ast.IsIndex() {
		return hindley_milner.E_PROXY
	}
	return hindley_milner.E_APPLICATION
}

///

func (ast *UnaryApplication) RenameVariables(subst cfg.VariableSubstitution) {
	// if ast.IsApplication() {
	// 	v := subst.Replace(*ast.Target)
	// 	ast.Target = &v
	// }
}

func (ast *UnaryApplication) GetUsedVariables(vars cfg.VariableSet, visitedMap map[generic_ast.TraversableNode]struct{}) cfg.VariableSet {
	if ast.IsApplication() {
		vars.Insert(cfg.GetAllUsagesVariables(ast.Index, map[generic_ast.TraversableNode]struct{}{}))
		for _, arg := range ast.Arguments {
			vars.Insert(cfg.GetAllUsagesVariables(arg, map[generic_ast.TraversableNode]struct{}{}))
		}
	} else if ast.IsIndex() {
		vars.Insert(cfg.GetAllUsagesVariables(ast.Index, map[generic_ast.TraversableNode]struct{}{}))
	}
	return vars
}
