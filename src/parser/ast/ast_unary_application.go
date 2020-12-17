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
	Target *string   `( @Ident`
	Arguments []*Expression   `"(" (@@ ("," @@)*)? ")" )`
	Index *Index `| @@`
	ParentNode generic_ast.TraversableNode
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
		nodes := make([]generic_ast.TraversableNode, len(ast.Arguments) + 1)
		nodes = append(nodes, generic_ast.MakeTraversableNodeToken(ast, *ast.Target, ast.Pos, ast.EndPos))
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
	return ast.Target != nil
}

func (ast *UnaryApplication) IsIndex() bool {
	return ast.Index != nil
}

func (ast *UnaryApplication) Print(c *context.ParsingContext) string {
	if ast.IsApplication() {
		args := []string{}
		for _, argument := range ast.Arguments {
			args = append(args, argument.Print(c))
		}
		return printNode(c, ast, "%s(%s)", *ast.Target, strings.Join(args, ", "))
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
			Target:      ast.Target,
			Arguments:   args,
			ParentNode: parent.(generic_ast.TraversableNode),
		}, context, true)
	} else if ast.IsIndex() {
		return mapper(parent, &UnaryApplication{
			BaseASTNode: ast.BaseASTNode,
			Index: mapper(ast, ast.Index, context, false).(*Index),
			ParentNode: parent.(generic_ast.TraversableNode),
		}, context, true)
	}
	panic("Invalid UnaryApplication operation type")
}

func (ast *UnaryApplication) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	if ast.IsIndex() {
		mapper(ast, ast.Index, context)
	} else if ast.IsApplication() {
		for _, arg := range ast.Arguments {
			mapper(ast, arg, context)
		}
	}
	mapper(parent, ast, context)
}

func (ast *UnaryApplication) Fn() generic_ast.Expression {
	return &VarName{
		BaseASTNode: ast.BaseASTNode,
		name: *ast.Target,
	}
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
						return hindley_milner.NewScheme(nil, CreatePrimitive(T_VOID))
					},
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
	if ast.IsApplication() {
		v := subst.Replace(*ast.Target)
		ast.Target = &v
	}
}

func (ast *UnaryApplication) GetUsedVariables(vars cfg.VariableSet) cfg.VariableSet {
	if ast.IsApplication() {
		vars.Add(cfg.NewVariable(*ast.Target, nil))
	}
	if ast.IsApplication() {
		for _, arg := range ast.Arguments {
			vars.Insert(cfg.GetAllUsagesVariables(arg))
		}
	} else if ast.IsIndex() {
		vars.Insert(cfg.GetAllUsagesVariables(ast.Index))
	}
	return vars
}