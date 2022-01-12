package ir

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/parser/utils"
)

type IRPhi struct {
	generic_ast.BaseASTNode
	TargetName string
	Type       IRType
	Values     []string
	Blocks     []int
	ParentNode generic_ast.TraversableNode
}

func CreateIRPhi(target string, varType IRType, phiBlocks map[int]string) *IRPhi {

	phiBlocksIDs := []int{}
	phiValues := []string{}
	for phiBlock, phiValue := range phiBlocks {
		phiBlocksIDs = append(phiBlocksIDs, phiBlock)
		phiValues = append(phiValues, phiValue)
	}

	return &IRPhi{
		TargetName: target,
		Values:     phiValues,
		Type:       varType,
		Blocks:     phiBlocksIDs,
	}
}

func (ast *IRPhi) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *IRPhi) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *IRPhi) Begin() lexer.Position {
	return ast.Pos
}

func (ast *IRPhi) End() lexer.Position {
	return ast.EndPos
}

func (ast *IRPhi) GetNode() interface{} {
	return ast
}

func (ast *IRPhi) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{}
}

func (ast *IRPhi) Print(c *context.ParsingContext) string {
	description := []string{}
	for i, varName := range ast.Values {
		description = append(description, fmt.Sprintf("%d: value %s", ast.Blocks[i], varName))
	}

	return utils.PrintASTNode(c, ast, "%s %s = Phi [%s]", ast.Type, ast.TargetName, strings.Join(description, ", "))
}

func (ast *IRPhi) GetUsedVariables(vars cfg.VariableSet, visitedMap map[generic_ast.TraversableNode]struct{}) cfg.VariableSet {

	return vars
}

func (ast *IRPhi) RenameVariables(substUsed, substDecl cfg.VariableSubstitution) {

}
