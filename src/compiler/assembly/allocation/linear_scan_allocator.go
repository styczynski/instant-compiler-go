package allocation

import (
	"fmt"
	"strings"

	"github.com/styczynski/latte-compiler/src/compiler/assembly/x86"
	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/ir"
)

type AllocConsAllowAll struct{}

func (cons *AllocConsAllowAll) String() string {
	return "ALLOW_ALL"
}

func (cons *AllocConsAllowAll) IsAllocationConstraint() {}

type AllocConsRequireRegisters struct {
}

func (cons *AllocConsRequireRegisters) String() string {
	return "MUST_BE_REGISTER"
}

func (cons *AllocConsRequireRegisters) IsAllocationConstraint() {}

type AllocConsRequireSpecificRegisters struct {
	AllowedRegisters []x86.Reg
}

func (cons *AllocConsRequireSpecificRegisters) String() string {
	regNames := []string{}
	for _, reg := range cons.AllowedRegisters {
		regNames = append(regNames, fmt.Sprintf("%v", reg))
	}
	return fmt.Sprintf("MUST_BE_REGISTER_IN[%s]", strings.Join(regNames, ","))
}

func (cons *AllocConsRequireSpecificRegisters) IsAllocationConstraint() {}

type LocationMemory struct {
	Index int
	Size  int
}

func (l *LocationMemory) String() string {
	return fmt.Sprintf("Memory[%d(size=%d)]", l.Index, l.Size)
}

func (l *LocationMemory) IsAllocation() {}

type LocationRegister struct {
	Reg   x86.Reg
	Size  int
	State *RegistryState
}

func (l *LocationRegister) IsAllocation() {}

func (l *LocationRegister) String() string {
	return fmt.Sprintf("Register[%s(size=%d)]", l.Reg.String(), l.Size)
}

func IsAllocMem(alloc ir.IRAllocation) (*LocationMemory, bool) {
	if lm, ok := alloc.(*LocationMemory); ok {
		return lm, ok
	}
	return nil, false
}

func IsAllocReg(alloc ir.IRAllocation) (*LocationRegister, bool) {
	if lreg, ok := alloc.(*LocationRegister); ok {
		return lreg, ok
	}
	return nil, false
}

type RegistryState struct {
	Full  bool
	Var   string
	Specs *x86.RegistrySpecs
}

func (regState *RegistryState) Copy() *RegistryState {
	return &RegistryState{
		Full:  regState.Full,
		Var:   regState.Var,
		Specs: regState.Specs,
	}
}

type AssemblyFunctionMeta struct {
	VarLen int
}

type LinearScanAllocatorState struct {
	Current             ir.IRAllocationMap
	All                 ir.IRAllocationMap
	LastBlock           int
	AvailableRegistries map[x86.Reg]*RegistryState
}

func (state *LinearScanAllocatorState) allocateBlockUsing(start int, size int) bool {
	if state.IsBlockAvailable(start, size) {
		return false
	}
	newLast := start + size
	if newLast > state.LastBlock {
		state.LastBlock = newLast
	}
	return true
}

func (state *LinearScanAllocatorState) AllocateBlock(start int, size int) {
	newLast := start + size
	if newLast > state.LastBlock {
		state.LastBlock = newLast
	}
}

func (state *LinearScanAllocatorState) GetFunctionMeta(fn *ir.IRFunction) ir.FunctionMeta {
	return AssemblyFunctionMeta{
		VarLen: state.LastBlock,
	}
}

func (state *LinearScanAllocatorState) IsBlockAvailable(start int, size int) bool {
	for _, alloc := range state.Current {
		if mem, ok := alloc.(*LocationMemory); ok {
			if (mem.Index >= start && mem.Index < start+size) || (mem.Index+mem.Size-1 >= start && mem.Index+mem.Size-1 < start+size) {
				return false
			}
		}
	}
	return true
}

func (state *LinearScanAllocatorState) allocateAvailableRegistryUsing(name string, size int, reg x86.Reg) (*LocationRegister, bool) {
	if regState, ok := state.AvailableRegistries[reg]; ok {
		if !regState.Full && regState.Specs.DefaultSize >= size {
			regState.Full = true
			regState.Var = name
			return &LocationRegister{
				Reg:   reg,
				Size:  size,
				State: regState.Copy(),
			}, true
		} else {
			return nil, false
		}
	}
	return nil, false
}

func (state *LinearScanAllocatorState) allocateAvailableRegistry(name string, size int) (*LocationRegister, bool) {
	for reg, regState := range state.AvailableRegistries {
		if !regState.Full && regState.Specs.DefaultSize >= size {
			regState.Full = true
			regState.Var = name
			return &LocationRegister{
				Reg:   reg,
				Size:  size,
				State: regState.Copy(),
			}, true
		}
	}
	return nil, false
}

func (state *LinearScanAllocatorState) PreserveOnly(allVars cfg.VariableSet) {
	state.Current.PreserveOnly(allVars)
	for _, regState := range state.AvailableRegistries {
		isOk := false
		for varName, _ := range allVars {
			if regState.Var == varName && regState.Full {
				isOk = true
				break
			}
		}
		if !isOk {
			regState.Var = ""
			regState.Full = false
		}
	}
}

type LinearScanAllocator struct {
	state      *LinearScanAllocatorState
	lockedRegs []x86.Reg
}

func (alloc *LinearScanAllocator) AreContraintsMet(currentAlloc ir.IRAllocation, cons ir.IRAllocationConstraints) bool {
	for _, con := range cons {
		if _, ok := con.(*AllocConsAllowAll); ok {
			return true
		} else if consReg, ok := con.(*AllocConsRequireSpecificRegisters); ok {
			if alloc, ok := IsAllocReg(currentAlloc); ok {
				a := alloc.Reg
				isOk := false
				for _, allowedReg := range consReg.AllowedRegisters {
					if x86.AreRegsColliding(&a, &allowedReg) {
						isOk = true
						break
					}
				}
				if !isOk {
					return false
				}
			}
		} else if _, ok := con.(*AllocConsRequireRegisters); ok {
			if _, ok := IsAllocReg(currentAlloc); !ok {
				return false
			}
		} else {
			panic(fmt.Sprintf("Unknown constraint was met: %s", con.String()))
		}
	}
	return false
}

func (alloc *LinearScanAllocator) getStrongRegRequirements(cons ir.IRAllocationConstraints) []x86.Reg {
	reqs := map[x86.Reg]struct{}{}
	for _, con := range cons {
		if conReg, ok := con.(*AllocConsRequireSpecificRegisters); ok {
			for _, reg := range conReg.AllowedRegisters {
				reqs[reg] = struct{}{}
			}
		}
	}
	reqsRegs := []x86.Reg{}
	for reg, _ := range reqs {
		reqsRegs = append(reqsRegs, reg)
	}
	return reqsRegs
}

func (alloc *LinearScanAllocator) allocateVar(name string, varType ir.IRType, hasExistingAlloc bool, existingAllocName string, existingAlloc ir.IRAllocation, cons ir.IRAllocationConstraints) ir.IRAllocation {
	// Check if the node has allocation data
	blockSize := ir.GetIRTypeSize(varType) / 8
	if currentAlloc, ok := alloc.state.All[name]; ok {
		// Skip already allocated variables
		return currentAlloc
	}

	strongRegReqs := alloc.getStrongRegRequirements(cons)
	if len(strongRegReqs) > 0 {
		// Check registers with strong constraints
		var newReg *x86.Reg = nil
		for _, reg := range strongRegReqs {
			if regState, ok := alloc.state.AvailableRegistries[reg]; ok {
				if !regState.Full {
					newReg = &reg
					break
				}
			}
		}
		if newReg != nil {
			regState := alloc.state.AvailableRegistries[*newReg]
			regState.Full = true
			regState.Var = name
			return &LocationRegister{
				Reg:   *newReg,
				Size:  regState.Specs.DefaultSize,
				State: regState.Copy(),
			}
		} else {
			// No registry is free
			// We choose the first one
			// TODO: Implement?
			panic("Critical problem cannot allocate register for strong register requirements")
		}
	}

	if hasExistingAlloc {
		if alloc.AreContraintsMet(existingAlloc, cons) {
			// Try to alloc

			if reg, ok := existingAlloc.(*LocationRegister); ok {
				allocReg, regOk := alloc.state.allocateAvailableRegistryUsing(name, blockSize, reg.Reg)
				if regOk {
					return allocReg
				}
			} else if mem, ok := existingAlloc.(*LocationMemory); ok {
				ok := alloc.state.allocateBlockUsing(mem.Index, mem.Size)
				if ok {
					return mem
				}
			}
		}
	}

	freeMemoryIndex := 0
	allocReg, regOk := alloc.state.allocateAvailableRegistry(name, blockSize)
	if regOk {
		return allocReg
	}

	for {
		// Check if block is available
		if alloc.state.IsBlockAvailable(freeMemoryIndex, blockSize) {
			break
		}
		freeMemoryIndex++
	}
	alloc.state.AllocateBlock(freeMemoryIndex, blockSize)
	return &LocationMemory{
		Index: freeMemoryIndex,
		Size:  blockSize,
	}
}

func (alloc *LinearScanAllocator) Initialize() {
	// Empty
}

func (alloc *LinearScanAllocator) Lock(lockedRegs []x86.Reg) *LinearScanAllocator {
	alloc.lockedRegs = lockedRegs
	return alloc
}

func (alloc *LinearScanAllocator) ResetSettings() {
	alloc.state = nil
	alloc.lockedRegs = nil
}

func (alloc *LinearScanAllocator) PerformAllocationForBlocks(blocks []*ir.IRBlock) AllocatorState {

	regs := map[x86.Reg]*RegistryState{}
	for reg, regSpecs := range x86.ALL_REGS {
		if regSpecs.ForbidForAllocation {
			continue
		}
		isLocked := false
		for _, lockedReg := range alloc.lockedRegs {
			if lockedReg == reg {
				isLocked = true
				break
			}
		}
		if !isLocked {
			regs[reg] = &RegistryState{
				Full:  false,
				Var:   "",
				Specs: regSpecs,
			}
		}
	}

	alloc.state = &LinearScanAllocatorState{
		Current:             ir.IRAllocationMap{},
		All:                 ir.IRAllocationMap{},
		AvailableRegistries: regs,
	}
	for _, block := range blocks {
		for _, stmt := range block.Statements {
			decl := cfg.GetAllDeclaredVariables(stmt, map[generic_ast.TraversableNode]struct{}{})
			stmtAlloc := ir.IRAllocationMap{}
			// Cleanup allocation using live ariables
			allVars := stmt.VarIn.Copy()
			allVars.Insert(stmt.VarOut)
			alloc.state.PreserveOnly(allVars)
			cons := ir.IRAllocationContraintsMap{}

			existingAllocName, existingAlloc, hasExistingAlloc := stmt.TryToGetAllocationTarget()
			_, existingCons := stmt.GetAllocationTargetContraints()

			for varName, _ := range decl {
				// TODO: Fix ir type
				loc := alloc.allocateVar(varName, ir.IR_INT32, hasExistingAlloc, existingAllocName, existingAlloc, existingCons)
				alloc.state.Current[varName] = loc
				alloc.state.All[varName] = loc
				stmtAlloc[varName] = loc
				cons[varName] = ir.IRAllocationConstraints{
					&AllocConsAllowAll{},
				}
			}
			stmt.SetAllocationInfo(stmtAlloc, alloc.state.All.Copy())
			stmt.SetTargetAllocationConstraintsMap(cons)
		}
	}
	return alloc.state
}
