package cfg

import (
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
)

type ConstFoldingError interface {
	GetSource() generic_ast.TraversableNode
	GetMessage() string
	Error() string
}

type BlockLiveVariables interface {
	BlockIn(blockID int) VariableSet
	BlockOut(blockID int) VariableSet
}

type FlowAnalysis interface {
	Graph() *CFG
	Liveness() BlockLiveVariables
	Reaching() ReachingVariablesInfo
	Print(c *context.ParsingContext) string
	Optimize(c *context.ParsingContext)
	ConstFold(c *context.ParsingContext) ConstFoldingError
	Output() generic_ast.NormalNode
	Rebuild()
	RebuildCalculatedProperties()
}

type FlowAnalysisImpl struct {
	input    []generic_ast.NormalNode
	graph    *CFG
	liveness *LiveVariablesInfo
	reaching *ReachingVariablesInfo
}

func (flow *FlowAnalysisImpl) Output() generic_ast.NormalNode {
	return flow.input[0]
}

func (flow *FlowAnalysisImpl) RebuildCalculatedProperties() {
	flow.reaching = nil
	flow.liveness = nil
	flow.Liveness()
	flow.Reaching()
}

func (flow *FlowAnalysisImpl) Rebuild() {
	flow.reaching = nil
	flow.liveness = nil
	flow.graph = nil
}

func (flow *FlowAnalysisImpl) ReplaceBlock(blockID int, new generic_ast.NormalNode) {
	g := flow.Graph()
	for i, block := range flow.input {
		//if block == old {
		//	flow.input[i] = new
		//}
		flow.input[i] = generic_ast.ReplaceExpressionRecursively(block, g.codeMapping[blockID], new).(generic_ast.NormalNode)
	}
	if flow.graph != nil {
		flow.graph.ReplaceBlockCode(blockID, new)
	}
	// if flow.liveness != nil {
	// 	flow.liveness.ReplaceBlockCode(blockID, new)
	// }
	// if flow.reaching != nil {
	// 	flow.reaching.ReplaceBlockCode(blockID, new)
	// }
}

func CreateFlowAnalysis(input generic_ast.NormalNode) FlowAnalysis {
	return &FlowAnalysisImpl{
		input: []generic_ast.NormalNode{input},
	}
}

func (flow *FlowAnalysisImpl) Reaching() ReachingVariablesInfo {
	if flow.reaching == nil {
		v := DefUse(flow.Graph())
		flow.reaching = &v
	}
	return *flow.reaching
}

func (flow *FlowAnalysisImpl) Graph() *CFG {
	if flow.graph == nil {
		flow.graph = CreateCFGFromNodes(flow.input)
	}
	return flow.graph
}

func (flow *FlowAnalysisImpl) Liveness() BlockLiveVariables {
	if flow.liveness == nil {
		v := LiveVars(flow.Graph())
		flow.liveness = &v
	}
	return *flow.liveness
}
