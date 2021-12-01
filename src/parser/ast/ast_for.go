package ast

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type For struct {
	generic_ast.BaseASTNode
	ElementType *Type          `"for" "(" @@`
	Destructor  *ForDestructor `@@ ")"`
	Do          *Statement     `@@`
	ParentNode  generic_ast.TraversableNode
}

func (ast *For) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *For) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *For) Begin() lexer.Position {
	return ast.Pos
}

func (ast *For) End() lexer.Position {
	return ast.EndPos
}

func (ast *For) GetNode() interface{} {
	return ast
}

func (ast *For) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{
		ast.Destructor,
		ast.Do,
	}
}

func (ast *For) Print(c *context.ParsingContext) string {
	return printNode(c, ast, "for (%s %s: %s) %s",
		ast.ElementType.Print(c),
		ast.Destructor.Print(c),
		ast.Do.Print(c),
	)
}

func (ast *For) Validate(c *context.ParsingContext) generic_ast.NodeError {
	if ast.Do != nil {
		if ast.Do.IsDeclaration() {
			message := fmt.Sprintf("Declaration as a non-block expression inside for statement is forbidden. Please use { } brackets and create the definition there.")
			return generic_ast.NewNodeError(
				"Declaration not allowed",
				ast,
				message,
				message)
		}
	}
	return nil
}

///

func (ast *For) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	// TODO
	return mapper(parent, &For{
		BaseASTNode: ast.BaseASTNode,
		ElementType: ast.ElementType,
		Destructor:  mapper(ast, ast.Destructor, context, false).(*ForDestructor),
		Do:          mapper(ast, ast.Do, context, false).(*Statement),
		ParentNode:  parent.(generic_ast.TraversableNode),
	}, context, true).(generic_ast.Expression)
}
func (ast *For) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(ast, ast.Destructor, context)
	mapper(ast, ast.Do, context)
	mapper(parent, ast, context)
}

func (ast *For) Var(c hindley_milner.InferContext) *hindley_milner.NameGroup {
	types := map[string]*hindley_milner.Scheme{}
	types[ast.Destructor.ElementVar] = ast.ElementType.GetType(c)
	return hindley_milner.NamesWithTypes([]string{ast.Destructor.ElementVar}, types)
}

func (ast *For) Def(c hindley_milner.InferContext) generic_ast.Expression {
	return ast.Destructor
}

func (ast *For) Body() generic_ast.Expression {
	return hindley_milner.Batch{Exp: []generic_ast.Expression{
		ast.Do,
	}}
}

func (ast *For) ExpressionType() hindley_milner.ExpressionType {
	return hindley_milner.E_LET_RECURSIVE
}

///

func (ast *For) BuildFlowGraph(builder cfg.CFGBuilder) {
	// flows as such (range same w/o init & post):
	// previous -> [ init -> ] for -> body -> [ post -> ] for -> next

	var post generic_ast.NormalNode = ast

	builder.AddBlockSuccesor(ast.Destructor)
	builder.UpdatePrev([]generic_ast.NormalNode{ast.Destructor})
	builder.AddBlockSuccesor(ast)

	builder.UpdatePrev([]generic_ast.NormalNode{ast})
	builder.BuildNode(ast.Do)

	builder.AddBlockSuccesor(post)

	ctrlExits := []generic_ast.NormalNode{ast}

	// handle any branches; if no label or for me: handle and remove from branches.
	for i := 0; i < len(builder.Branches()); i++ {
		//br := builder.Branches()[i]
		// Deal with continue/break here if such thing will be implemented
	}

	builder.UpdatePrev(ctrlExits) // for stmt and any appropriate break statements
}

func (ast *For) GetUsedVariables(vars cfg.VariableSet, visitedMap map[generic_ast.TraversableNode]struct{}) cfg.VariableSet {
	return cfg.GetAllVariables(ast.Destructor.Target, visitedMap)
}

func (ast *For) GetDeclaredVariables(visitedMap map[generic_ast.TraversableNode]struct{}) cfg.VariableSet {
	return cfg.NewVariableSet(cfg.NewVariable(ast.Destructor.ElementVar, ast.Destructor))
}
