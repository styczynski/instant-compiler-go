package assembly

import (
	"fmt"

	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/ir"
)

type LocationMemory struct {
	Index int
}

type AllocatorState struct {
	Current ir.IRAllocationMap
}

func (l *LocationMemory) String() string {
	return fmt.Sprintf("Memory[%d]", l.Index)
}

func allocateVar(name string, state *AllocatorState) ir.IRAllocation {
	freeMemoryIndex := 0
	for {
		isOk := true
		for _, alloc := range state.Current {
			if mem, ok := alloc.(*LocationMemory); ok {
				if mem.Index == freeMemoryIndex {
					isOk = false
					break
				}
			}
		}
		if isOk {
			break
		}
		freeMemoryIndex++
	}
	return &LocationMemory{
		Index: freeMemoryIndex,
	}
}

func performAllocationForBlocks(blocks []*ir.IRBlock) {
	state := &AllocatorState{
		Current: ir.IRAllocationMap{},
	}
	for _, block := range blocks {
		for _, stmt := range block.Statements {
			decl := cfg.GetAllDeclaredVariables(stmt, map[generic_ast.TraversableNode]struct{}{})
			stmtAlloc := ir.IRAllocationMap{}
			for varName, _ := range decl {
				loc := allocateVar(varName, state)
				state.Current[varName] = loc
				stmtAlloc[varName] = loc
			}
			stmt.SetAllocationInfo(stmtAlloc)
			// Cleanup allocation using live ariables
			state.Current.PreserveOnly(stmt.VarOut)
		}
	}
}
