package cfg

import (
	"reflect"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
)

type ConstFoldingErrorImpl struct {
	Message string
	Source  generic_ast.TraversableNode
}

func (e *ConstFoldingErrorImpl) GetSource() generic_ast.TraversableNode {
	return e.Source
}

func (e *ConstFoldingErrorImpl) GetMessage() string {
	return e.Message
}

func (e *ConstFoldingErrorImpl) Error() string {
	return e.Message
}

func (flow *FlowAnalysisImpl) optimizeConst(node generic_ast.TraversableNode, c *context.ParsingContext, validate bool) (generic_ast.TraversableNode, error) {
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
	if expr, ok := node.(generic_ast.Expression); ok {
		var mapper generic_ast.ExpressionVisitor
		var invalidChecker generic_ast.ExpressionVisitor
		visitedNodes := map[generic_ast.Expression]struct{}{}
		mapper = func(parent generic_ast.Expression, e generic_ast.Expression, context generic_ast.VisitorContext) {
			if _, wasVisited := visitedNodes[e]; wasVisited {
				return
			}
			//
			visitedNodes[e] = struct{}{}
			e.Visit(parent, mapper, context)
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
		}
		expr.Visit(node.Parent().(generic_ast.Expression), mapper, generic_ast.NewEmptyVisitorContext())

		if validate {
			var invalidNodeError error = nil
			var invalidNodeSrc generic_ast.TraversableNode = nil
			visitedNodes = map[generic_ast.Expression]struct{}{}
			invalidChecker = func(parent generic_ast.Expression, e generic_ast.Expression, context generic_ast.VisitorContext) {
				if _, wasVisited := visitedNodes[e]; wasVisited {
					return
				}
				visitedNodes[e] = struct{}{}
				e.Visit(parent, invalidChecker, context)
				if nodeWithFoldingValidation, ok := e.(generic_ast.NodeWithFoldingValidation); ok {
					err, src := nodeWithFoldingValidation.ValidateConstFold()
					if err != nil {
						if invalidNodeError == nil {
							invalidNodeError = err
							invalidNodeSrc = src
						}
						return
					}
				}
			}
			expr.Visit(node.Parent().(generic_ast.Expression), invalidChecker, generic_ast.NewEmptyVisitorContext())
			if invalidNodeError != nil {
				return invalidNodeSrc, invalidNodeError
			}
		}

		return expr.(generic_ast.TraversableNode), nil
	} else {
		return node, nil
	}
}

func (flow *FlowAnalysisImpl) ConstFold(c *context.ParsingContext) ConstFoldingError {

	g := flow.Graph()
	flow.Reaching()

	for i := 1; i <= 10; i++ {
		//fmt.Printf("refold...")
		
		change := false

		// Firstly, run const optimization on each node
		for _, block := range flow.graph.blocksOrder {
			if stmt, ok := g.codeMapping[block].(generic_ast.TraversableNode); ok {
				newNode, err := flow.optimizeConst(stmt, c, false)
				if err != nil {
					return &ConstFoldingErrorImpl{
						Message: err.Error(),
						Source:  newNode,
					}
				}
				generic_ast.ReplaceExpressionRecursively(stmt, stmt, newNode)
			}
		}
		//fmt.Printf("refold A")

		for _, block := range flow.graph.blocksOrder {
			//fmt.Printf("refold beforeG")
			vars := flow.graph.ReferencedVars(g.codeMapping[block])
			//fmt.Printf("refold G")
			for _, variable := range vars.use {
				//fmt.Printf("refold B")

				var variableDecl *Variable = nil
				
				for _, reachingBlock := range flow.graph.blocksOrder {
					//fmt.Printf("refold C")

					if _, hasBlock := flow.reaching.ReachedBlocks(reachingBlock)[block]; hasBlock {
						for _, defVar := range flow.graph.ReferencedVars(g.codeMapping[reachingBlock]).decl {
							if defVar.Name() == variable.Name() {
								variableDecl = &defVar
								break
							}
						}
					}
					
				}
				
				//fmt.Printf("refold F")
				if variableDecl == nil {
					
					continue
				}

				//fmt.Printf("refold D")
				node := (*variableDecl).Value()
				if !isNilNode(node) {
					if constExtractable, ok := node.(generic_ast.ConstExtractableNode); ok {
						varConst, isConst := constExtractable.ExtractConst()
						if isConst {
							val := variable.Value()
							if val == nil || (reflect.ValueOf(val).Kind() == reflect.Ptr && reflect.ValueOf(val).IsNil()) {
								// Do nothing
							} else if val != varConst {
								generic_ast.ReplaceExpressionRecursively(g.codeMapping[block], val, varConst)
							}
						}
					}
				}
				//fmt.Printf("refold E")
				
			}
		}

		if !change {
			break
		}
	}

	
	// Validate
	for _, block := range flow.graph.blocksOrder {
		if stmt, ok := g.codeMapping[block].(generic_ast.TraversableNode); ok {
			newNode, err := flow.optimizeConst(stmt, c, true)
			if err != nil {
				return &ConstFoldingErrorImpl{
					Message: err.Error(),
					Source:  newNode,
				}
			}
		}
	}

	
	return nil
}
