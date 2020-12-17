package cfg

import (
	"reflect"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
)

func (flow *FlowAnalysisImpl) optimizeConst(node generic_ast.TraversableNode) generic_ast.TraversableNode {
	//if expr, ok := node.(generic_ast.Expression); ok{
	//	var mapper generic_ast.ExpressionMapper
	//	mapper = func (parent generic_ast.Expression, e generic_ast.Expression, context generic_ast.VisitorContext, backwards bool) generic_ast.Expression{
	//		if backwards {
	//			if optimizableNode, ok := e.(generic_ast.ConstFoldableNode); ok {
	//				return optimizableNode.ConstFold().(generic_ast.Expression)
	//			}
	//			return e
	//		}
	//		newNode := e.Map(parent, mapper, context)
	//		if optimizableNode, ok := newNode.(generic_ast.ConstFoldableNode); ok {
	//			return optimizableNode.ConstFold().(generic_ast.Expression)
	//		}
	//		return newNode.Map(parent, mapper, context)
	//	}
	//	return expr.Map(node.Parent().(generic_ast.Expression), mapper, generic_ast.NewEmptyVisitorContext()).(generic_ast.TraversableNode)
	//} else {
	//	return node
	//}
	if expr, ok := node.(generic_ast.Expression); ok{
		var mapper generic_ast.ExpressionVisitor
		visitedNodes := map[generic_ast.Expression]struct{}{}
		mapper = func (parent generic_ast.Expression, e generic_ast.Expression, context generic_ast.VisitorContext) {
			if _, wasVisited := visitedNodes[e]; wasVisited {
				return
			}
			visitedNodes[e] = struct{}{}
			if optimizableNode, ok := e.(generic_ast.ConstFoldableNode); ok {
				//generic_ast.ReplaceExpressionRecursively(e.(generic_ast.TraversableNode), e.(generic_ast.TraversableNode), optimizableNode.ConstFold())
				optimizedNode := optimizableNode.ConstFold()
				if reflect.TypeOf(e).Kind() == reflect.Ptr {
					// Pointer:
					reflect.ValueOf(e).Elem().Set(reflect.ValueOf(optimizedNode).Elem())
				} else {
					panic("Wrong type")
				}
			}
				e.Visit(parent, mapper, context)

		}
		expr.Visit(node.Parent().(generic_ast.Expression), mapper, generic_ast.NewEmptyVisitorContext())
		return expr.(generic_ast.TraversableNode)
	} else {
		return node
	}
}

func (flow *FlowAnalysisImpl) ConstFold(c *context.ParsingContext) {

	flow.Graph()
	flow.Reaching()

 	for true {
 		change := false

		// Firstly, run const optimization on each node
		for _, block := range flow.graph.blocksOrder {
			if block != nil {
				newNode := flow.optimizeConst(block)
				generic_ast.ReplaceExpressionRecursively(block, block, newNode)
			}
		}

		for _, block := range flow.graph.blocksOrder {
			vars := flow.graph.ReferencedVars(block)
			for _, variable := range vars.use {

				var variableDecl *Variable = nil
				for _, reachingBlock := range flow.graph.blocksOrder {
					if _, hasBlock := flow.reaching.ReachedBlocks(reachingBlock)[block]; hasBlock {
						for _, defVar := range flow.graph.ReferencedVars(reachingBlock).decl {
							if defVar.Name() == variable.Name() {
								variableDecl = &defVar
								break
							}
						}
					}
				}

				if variableDecl == nil {
					continue
				}

				node := (*variableDecl).Value()
				if !isNilNode(node) {
					if constExtractable, ok := node.(generic_ast.ConstExtractableNode); ok {
						varConst, isConst := constExtractable.ExtractConst()
						if isConst {
							generic_ast.ReplaceExpressionRecursively(block, variable.Value(), varConst)
							change = true
						}
					}
				}
			}
		}

		if !change {
			break
		}
	}
}