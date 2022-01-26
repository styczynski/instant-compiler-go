package allocation

import (
	"fmt"

	"github.com/styczynski/latte-compiler/src/compiler/assembly/x86"
	"github.com/styczynski/latte-compiler/src/ir"
)

func DoCall(
	label string,
	targetType ir.IRType,
	targetAlloc ir.IRAllocation,
	argsOrder []string,
	argsAllocs ir.IRAllocationMap,
	allocationContext ir.IRAllocationMap,
	allocationOutputContext ir.IRAllocationMap,
	registriesOverrides map[x86.Reg]int64,
) []*x86.Instruction {

	usedRegs := []*x86.RegistrySpecs{}
	for _, alloc := range allocationOutputContext {
		if allocReg, ok := IsAllocReg(alloc); ok {
			usedRegs = append(usedRegs, allocReg.State.Specs)
		}
	}

	ret := []*x86.Instruction{}

	targetRegs := []*x86.RegistrySpecs{}
	targetRegsSet := map[x86.Reg]*x86.RegistrySpecs{}
	for argNo, _ := range argsOrder {
		specs := x86.GetRegisterForFunctionArg(argNo)
		if specs != nil {
			targetRegs = append(targetRegs, specs)
			targetRegsSet[specs.Normalized] = specs
		}
	}

	regsToPreserve := map[x86.Reg]*x86.RegistrySpecs{}

	regsArgsLocations := []*x86.RegistrySpecs{}
	regsArgsLocationsIndexes := []int{}
	regsArgsLocationsSet := map[x86.Reg]struct{}{}
	regsArgsLocationsIndexesMap := map[*x86.RegistrySpecs]int{}

	argsMem := []*LocationMemory{}
	argsMemIndexes := []int{}

	// Determine registry overrides causing new preserved registries
	for target, _ := range registriesOverrides {
		_, mappedTarget := x86.ResizeReg(target, 4)
		regsToPreserve[mappedTarget] = x86.ALL_REGS[mappedTarget]
	}

	// Determine registry to preserve
	for _, usedReg := range usedRegs {
		if true || !usedReg.IsPreserved {
			regsToPreserve[usedReg.Normalized] = usedReg
		}
	}

	for argNo, argName := range argsOrder {
		arg := argsAllocs[argName]
		fmt.Printf("ARGUMENT %d FOR CALL IS ALLOCATED ON: %v\n", argNo, arg)
		if allocReg, ok := IsAllocReg(arg); ok {
			regsArgsLocations = append(regsArgsLocations, allocReg.State.Specs)
			regsArgsLocationsIndexes = append(regsArgsLocationsIndexes, argNo)
			regsArgsLocationsSet[allocReg.State.Specs.Normalized] = struct{}{}
			regsArgsLocationsIndexesMap[allocReg.State.Specs] = argNo
			//var regUsed *x86.RegistrySpecs = nil
			// Check used registries
			// for _, targetReg := range usedRegs {
			// 	if x86.AreRegsCollidingConst(&allocReg.Reg, targetReg.Reg) {
			// 		regUsed = targetReg
			// 		break
			// 	}
			// }
			// if regUsed != nil {
			// 	regsToPreserve[allocReg.Reg] = regUsed
			// }
		} else if allocMem, ok := IsAllocMem(arg); ok {
			argsMem = append(argsMem, allocMem)
			argsMemIndexes = append(argsMemIndexes, argNo)
		}
	}

	// Calculate return location
	doNotPreserveReturn := map[x86.Reg]struct{}{}
	returnTransferInstrs := []*x86.Instruction{}
	returnLocation := x86.GetRegisterForFunctionReturn(0).Reg8B
	if targetType != ir.IR_VOID {
		if targetAllocReg, ok := IsAllocReg(targetAlloc); ok {
			returnTransferInstrs = append(returnTransferInstrs, x86.DoRegistryCopy(targetAllocReg.Reg, returnLocation, targetAllocReg.Size))
			doNotPreserveReturn[targetAllocReg.State.Specs.Normalized] = struct{}{}
			fmt.Printf("RETURN LOC: %v\n", targetAllocReg.State.Specs.Normalized)
		} else if targetAllocMem, ok := IsAllocMem(targetAlloc); ok {
			returnTransferInstrs = append(returnTransferInstrs, x86.DoMemoryStore(targetAllocMem.Index, targetAllocMem.Size, returnLocation))
		}
	}

	// We need to preserver regsToPreserve before the call
	regsPreserveOrder := []x86.Reg{}
	for _, specs := range regsToPreserve {
		fmt.Printf("CHECK RETURN REG SHOULD PRESERVED BE?: %v\n", specs.Normalized)
		if _, ok := doNotPreserveReturn[specs.Normalized]; !ok {
			fmt.Printf("   PROCEED I CHUJ!\n")
			regsPreserveOrder = append(regsPreserveOrder, specs.Reg8B)
			ret = append(ret, x86.DoPush(specs.Reg8B, 8))
		}
	}

	// if targetType != ir.IR_VOID {
	// 	if targetRegAlloc, ok := IsAllocReg(targetAlloc); ok {
	// 		regsPreserveOrder = append(regsPreserveOrder, targetRegAlloc.State.Specs.Reg8B)
	// 		ret = append(ret, x86.DoPush(targetRegAlloc.State.Specs.Reg8B, 8))
	// 	}
	// }

	// Transfer args (handle only registers)
	for argNo, argReg := range regsArgsLocations {
		targetSpecs := targetRegs[regsArgsLocationsIndexes[argNo]]
		// Registry not used if teraget has override
		if _, hasOverride := registriesOverrides[targetSpecs.Normalized]; hasOverride {
			continue
		}
		fmt.Printf("ARGUMENT %d FOR CALL WILL BE TRANSFERED TO: %v\n", argNo, targetSpecs.Normalized)
		if _, ok := regsArgsLocationsSet[argReg.Normalized]; ok {
			// We have this register as an input so we do swap
			if targetSpecs.Reg8B != argReg.Reg8B {
				fmt.Printf("ARGUMENT SWAP %v %v\n", targetSpecs.Normalized, argReg.Normalized)
				ret = append(ret, x86.DoSwap(targetSpecs.Normalized, argReg.Normalized, argReg.DefaultSize))
				// Swap args in array
				regsArgsLocations[argNo], regsArgsLocations[regsArgsLocationsIndexesMap[targetSpecs]] = regsArgsLocations[regsArgsLocationsIndexesMap[targetSpecs]], regsArgsLocations[argNo]
			}
		} else {
			// We only do mov as it's safe
			ret = append(ret, x86.DoRegistryCopy(targetSpecs.Normalized, argReg.Normalized, argReg.DefaultSize))
		}
	}

	// Transfer args (memory)
	for memNo, argMem := range argsMem {
		ret = append(ret, x86.DoMemoryLoad(argMem.Index, argMem.Size, targetRegs[argsMemIndexes[memNo]].Reg8B))
	}

	// Set overrides
	for target, overrideValue := range registriesOverrides {
		ret = append(ret, x86.DoRegStoreConst(target, 4, overrideValue))
	}

	instCall := x86.Inst{}
	instCall.Op = x86.CALL
	instCall.Args = x86.Args{
		x86.CreateRelLabel(label),
	}
	ret = append(ret, &x86.Instruction{
		Inst: instCall,
	})

	// Transfer result
	ret = append(ret, returnTransferInstrs...)

	// We need retrieve values of the registers
	for i := len(regsPreserveOrder) - 1; i >= 0; i-- {
		ret = append(ret, x86.DoPop(regsPreserveOrder[i], 8))
	}

	return ret
}
