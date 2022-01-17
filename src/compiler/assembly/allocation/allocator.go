package allocation

import (
	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/ir"
)

type AllocatorState interface {
	PreserveOnly(allVars cfg.VariableSet)
	GetFunctionMeta(fn *ir.IRFunction) ir.FunctionMeta
}

type Allocator interface {
	Initialize()
	ResetSettings()
	PerformAllocationForBlocks(blocks []*ir.IRBlock) AllocatorState
}

func RunAllocator(program *ir.IRProgram, alloc Allocator) {
	alloc.Initialize()
	for _, fn := range program.Statements {
		allocState := alloc.PerformAllocationForBlocks(fn.FunctionBody)
		fn.SetMeta(allocState.GetFunctionMeta(fn))
	}
}
