package ir

import (
	"fmt"

	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
)

func convertToSSA(graph *cfg.CFG, ir *IRGeneratorState) {
	visitedIDs := map[int]struct{}{}
	nameMapping := cfg.VariableSubstitutionMap{}
	allBlockOutputMappings := map[int]cfg.VariableSubstitutionMap{}
	allBlockInputMappings := map[int]cfg.VariableSubstitutionMap{}
	varTypes := map[string]IRType{}
	graph.VisitGraph(graph.Entry, func(g *cfg.CFG, block *cfg.Block, next func(blockID int)) {
		if _, wasVisited := visitedIDs[block.ID]; wasVisited {
			return
		}
		visitedIDs[block.ID] = struct{}{}
		code := graph.GetBlockCode(block.ID)
		if codeBlock, ok := code.(*IRBlock); ok {
			blockOutputsMapping := cfg.VariableSubstitutionMap{}
			allBlockInputMappings[block.ID] = nameMapping.Copy()
			for _, stmt := range codeBlock.Statements {
				vars := graph.ReferencedVars(stmt)
				newNameMapping := cfg.VariableSubstitutionMap{}
				for varName, _ := range vars.Assigned() {
					if !ir.isTempVar(varName) {
						newNameMapping[varName] = ir.NextVar(block.ID, varName)
						varTypes[varName] = stmt.ResolveTypeOfVar(varName)
					}
				}

				// for varName, _ := range vars.Used() {
				// 	if !nameMapping.Has(varName) && !ir.isTempVar(varName) {
				// 		nameMapping[varName] = ir.NextVar(varName)
				// 	}
				// }
				cfg.ReplaceVariables(stmt, nameMapping, newNameMapping, map[generic_ast.TraversableNode]struct{}{})
				nameMapping.Join(newNameMapping)
				blockOutputsMapping.Join(nameMapping)
			}
			allBlockOutputMappings[block.ID] = blockOutputsMapping
		}

		for _, stmt := range block.GetSuccs() {
			if _, wasVisited := visitedIDs[stmt]; wasVisited {
				continue
			}
			next(stmt)
		}
	})

	visitedIDs = map[int]struct{}{}
	graph.VisitGraph(graph.Entry, func(g *cfg.CFG, block *cfg.Block, next func(blockID int)) {
		if _, wasVisited := visitedIDs[block.ID]; wasVisited {
			return
		}
		visitedIDs[block.ID] = struct{}{}
		code := graph.GetBlockCode(block.ID)
		if codeBlock, ok := code.(*IRBlock); ok {
			blockPreds := block.GetPreds()
			if len(blockPreds) > 1 {
				//blockOutputMappings := allBlockOutputMappings[block.ID]
				blockInputMappings := allBlockInputMappings[block.ID]
				subst := cfg.VariableSubstitutionMap{}

				headers := []*IRStatement{}

				for origVarName, varName := range blockInputMappings {
					phiTarget := fmt.Sprintf("%s_phi_%d", origVarName, block.ID)
					phiBlocks := map[int]string{}
					phiValues := map[string]struct{}{}

					for _, pred := range blockPreds {
						predOutputs := allBlockOutputMappings[pred]
						if predVarName, ok := predOutputs[origVarName]; ok {
							phiBlocks[pred] = predVarName
							phiValues[predVarName] = struct{}{}
						}
					}

					if len(phiBlocks) > 0 && len(phiValues) > 1 {
						headers = append(headers, WrapIRPhi(CreateIRPhi(
							phiTarget,
							varTypes[origVarName],
							phiBlocks,
						)))
						subst[varName] = phiTarget
					}
				}

				if len(headers) > 0 {
					headers = append(headers, codeBlock.Statements...)
					codeBlock.Statements = headers
				}
				cfg.ReplaceVariables(code, subst, cfg.VariableSubstitutionMap{}, map[generic_ast.TraversableNode]struct{}{})
			}
		}
		for _, stmt := range block.GetSuccs() {
			if _, wasVisited := visitedIDs[stmt]; wasVisited {
				continue
			}
			next(stmt)
		}
	})
}
