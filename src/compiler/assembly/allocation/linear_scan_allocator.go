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

///
type AllocConsSkip struct {
}

func (cons *AllocConsSkip) String() string {
	return "SKIP_ALLOCATION"
}

func (cons *AllocConsSkip) IsAllocationConstraint() {}

///

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

type AllocConsRequireMemoryStackTop struct {
	Offset int
	Size   int
}

func (cons *AllocConsRequireMemoryStackTop) String() string {
	return fmt.Sprintf("MUST_BE_ON_TOP_OF_THE_STACK(OFFSET=%d, SIZE=%d)", cons.Offset, cons.Size)
}

func (cons *AllocConsRequireMemoryStackTop) IsAllocationConstraint() {}

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
	Current                ir.IRAllocationMap
	All                    ir.IRAllocationMap
	LastBlock              int
	AvailableRegistries    map[x86.Reg]*RegistryState
	StackTopOffsetSequence int
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

func (state *LinearScanAllocatorState) allocateAvailableBlock(blockSize int) *LocationMemory {
	freeMemoryIndex := 0
	for {
		// Check if block is available
		if state.IsBlockAvailable(freeMemoryIndex, blockSize) {
			break
		}
		freeMemoryIndex++
	}
	state.AllocateBlock(freeMemoryIndex, blockSize)
	return &LocationMemory{
		Index: freeMemoryIndex,
		Size:  blockSize,
	}
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

func (state *LinearScanAllocatorState) LastMemoryBlock() (*LocationMemory, int) {
	lastIndex := 0
	lastSize := 0
	var lastMem *LocationMemory = nil
	for _, alloc := range state.Current {
		if mem, ok := alloc.(*LocationMemory); ok {
			if lastIndex < mem.Index {
				lastMem, lastSize, lastIndex = mem, mem.Size, mem.Index
			}
		}
	}
	return lastMem, lastIndex + lastSize
}

func (state *LinearScanAllocatorState) PreserveOnly(allVars cfg.VariableSet) {
	state.Current.PreserveOnly(allVars)

	for _, regState := range state.AvailableRegistries {
		regState.Full = false
		regState.Var = ""
	}

	for varName, alloc := range state.Current {
		if reg, ok := IsAllocReg(alloc); ok {
			state.AvailableRegistries[reg.Reg].Full = true
			state.AvailableRegistries[reg.Reg].Var = varName
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
		} else if _, ok := con.(*AllocConsSkip); ok {
			return true
		} else if conStackTop, ok := con.(*AllocConsRequireMemoryStackTop); ok {
			if mem, ok := IsAllocMem(currentAlloc); ok {
				// Calculate stack top
				//_, lastIndex := alloc.state.LastMemoryBlock()

				// Check the block offsets
				return -conStackTop.Offset == mem.Index
			}
			return false
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
		//return currentAlloc
		existingAlloc := currentAlloc
		if reg, ok := existingAlloc.(*LocationRegister); ok {
			if regState, ok := alloc.state.AvailableRegistries[reg.Reg]; ok {
				regState.Full = true
				regState.Var = name
				fmt.Printf(" -> Allocated: Exisitng alloc (reg)\n")
				return existingAlloc
			} else {
				panic("Not allowed: Same register reallocated?")
			}
		} else if existingMem, ok := existingAlloc.(*LocationMemory); ok {
			ok := alloc.state.IsBlockAvailable(existingMem.Index, existingMem.Size)
			if ok {
				fmt.Printf(" -> Allocated: Exisitng alloc (mem - free)\n")
				return existingAlloc
			} else {
				for varName, alloc := range alloc.state.Current {
					if mem, ok := alloc.(*LocationMemory); ok {
						if mem.Index == existingMem.Index && mem.Size == existingMem.Size && varName == name {
							// Ok
							fmt.Printf(" -> Allocated: Exisitng alloc (mem)\n")
							return existingAlloc
						}
					}
				}
				panic("Not allowed: Same memory reallocated?")
			}
		}
		panic("Not allowed: Reaclloacted entity.")
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
			fmt.Printf(" -> Allocated: Strong reg requirement (%v)\n", regState.Specs.Normalized)
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
					fmt.Printf(" -> Allocated: Exisitng alloc (contraints for reg)\n")
					return allocReg
				}
			} else if mem, ok := existingAlloc.(*LocationMemory); ok {
				ok := alloc.state.IsBlockAvailable(mem.Index, mem.Size)
				if ok {
					fmt.Printf(" -> Allocated: Exisitng alloc (contraints for mem)\n")
					return mem
				}
			}
		}
	}

	// Check if memory should be allocated
	for _, con := range cons {
		if conStackTop, ok := con.(*AllocConsRequireMemoryStackTop); ok {
			//stackBlockSize := 4
			//_, end := alloc.state.LastMemoryBlock()
			//if alloc.state.StackTopOffsetSequence > 0 {
			//	end = alloc.state.StackTopOffsetSequence
			//}
			//alloc.state.StackTopOffsetSequence = end + stackBlockSize
			//fmt.Printf("ALLOCATE STACK TOP OFFSET=%d ON POS=%d SIZE=%d\n", conStackTop.Offset, end, stackBlockSize)
			//ok := alloc.state.IsBlockAvailable(end, stackBlockSize)
			//if !ok {
			//	panic(fmt.Sprintf("Failed to allocate block on the stack top with offset: %d", conStackTop.Offset))
			//}
			fmt.Printf(" -> Allocated: Stack top (mem)\n")
			return &LocationMemory{
				Index: -conStackTop.Offset,
				Size:  conStackTop.Size,
			}
		}
	}
	alloc.state.StackTopOffsetSequence = 0

	allocReg, regOk := alloc.state.allocateAvailableRegistry(name, blockSize)
	if regOk {
		fmt.Printf(" -> Allocated: Casual reg\n")
		return allocReg
	}

	fmt.Printf(" -> Allocated: Casual mem\n")
	return alloc.state.allocateAvailableBlock(blockSize)
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
			alloc.state.Current = alloc.state.All.Copy()
			allVars := stmt.VarIn.Copy()
			allVars.Insert(stmt.VarOut)
			alloc.state.PreserveOnly(allVars)
			cons := ir.IRAllocationContraintsMap{}

			existingAllocName, existingAlloc, hasExistingAlloc := stmt.TryToGetAllocationTarget()
			_, existingCons := stmt.GetAllocationTargetContraints()

			shouldSkip := false
			for _, cons := range existingCons {
				if _, ok := cons.(*AllocConsSkip); ok {
					shouldSkip = true
					break
				}
			}
			if !shouldSkip {
				for varName, _ := range decl {
					// TODO: Fix ir type
					loc := alloc.allocateVar(varName, ir.IR_INT32, hasExistingAlloc, existingAllocName, existingAlloc, existingCons)
					fmt.Printf("Allocated var %s => %v {%s}\n", varName, loc, alloc.state.Current.String())
					alloc.state.Current[varName] = loc
					alloc.state.All[varName] = loc
					stmtAlloc[varName] = loc
					cons[varName] = ir.IRAllocationConstraints{
						&AllocConsAllowAll{},
					}
				}
			} else {
				for varName, _ := range decl {
					// TODO: Fix ir type
					alloc.state.Current[varName] = nil
					alloc.state.All[varName] = nil
					stmtAlloc[varName] = nil
					cons[varName] = ir.IRAllocationConstraints{
						&AllocConsAllowAll{},
					}
				}
			}

			stmt.SetAllocationInfo(stmtAlloc, alloc.state.All.Copy())
			stmt.SetTargetAllocationConstraintsMap(cons)
		}
	}
	return alloc.state
}
