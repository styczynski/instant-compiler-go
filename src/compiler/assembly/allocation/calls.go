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
) []*x86.Instruction {

	usedRegs := []*x86.RegistrySpecs{}
	for _, alloc := range allocationContext {
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
			targetRegsSet[specs.Reg] = specs
		}
	}

	regsToPreserve := map[x86.Reg]*x86.RegistrySpecs{}

	regsArgsLocations := []*x86.RegistrySpecs{}
	regsArgsLocationsIndexes := []int{}
	regsArgsLocationsSet := map[x86.Reg]struct{}{}
	regsArgsLocationsIndexesMap := map[*x86.RegistrySpecs]int{}

	argsMem := []*LocationMemory{}
	argsMemIndexes := []int{}

	for argNo, argName := range argsOrder {
		arg := argsAllocs[argName]
		fmt.Printf("ARGUMENT %d FOR CALL IS ALLOCATED ON: %v\n", argNo, arg)
		if allocReg, ok := IsAllocReg(arg); ok {
			regsArgsLocations = append(regsArgsLocations, allocReg.State.Specs)
			regsArgsLocationsIndexes = append(regsArgsLocationsIndexes, argNo)
			regsArgsLocationsSet[allocReg.State.Specs.Reg] = struct{}{}
			regsArgsLocationsIndexesMap[allocReg.State.Specs] = argNo
			var regUsed *x86.RegistrySpecs = nil

			// Check used registries
			for _, targetReg := range usedRegs {
				if x86.AreRegsCollidingConst(&allocReg.Reg, targetReg.Reg) {
					regUsed = targetReg
					break
				}
			}
			if regUsed != nil {
				regsToPreserve[allocReg.Reg] = regUsed
			}
		} else if allocMem, ok := IsAllocMem(arg); ok {
			argsMem = append(argsMem, allocMem)
			argsMemIndexes = append(argsMemIndexes, argNo)
		}
	}

	// Calculate return location
	doNotPreserveReturn := map[x86.Reg]struct{}{}
	returnTransferInstrs := []*x86.Instruction{}
	returnLocation := x86.GetRegisterForFunctionReturn(0).TopReg
	if targetType != ir.IR_VOID {
		if targetAllocReg, ok := IsAllocReg(targetAlloc); ok {
			returnTransferInstrs = append(returnTransferInstrs, x86.DoRegistryCopy(targetAllocReg.Reg, returnLocation, targetAllocReg.Size))
			doNotPreserveReturn[targetAllocReg.State.Specs.Reg] = struct{}{}
		} else if targetAllocMem, ok := IsAllocMem(targetAlloc); ok {
			returnTransferInstrs = append(returnTransferInstrs, x86.DoMemoryStore(targetAllocMem.Index, targetAllocMem.Size, returnLocation))
		}
	}

	// We need to preserver regsToPreserve before the call
	for _, specs := range regsToPreserve {
		if _, ok := doNotPreserveReturn[specs.Reg]; !ok {
			ret = append(ret, x86.DoPush(specs.TopReg, 8))
		}
	}

	// Transfer args (handle only registers)
	for argNo, argReg := range regsArgsLocations {
		targetSpecs := targetRegs[regsArgsLocationsIndexes[argNo]]
		fmt.Printf("ARGUMENT %d FOR CALL WILL BE TRANSFERED TO: %v\n", argNo, targetSpecs.Reg)
		if _, ok := regsArgsLocationsSet[argReg.Reg]; ok {
			// We have this register as an input so we do swap
			if targetSpecs.TopReg != argReg.TopReg {
				fmt.Printf("ARGUMENT SWAP %v %v\n", targetSpecs.Reg, argReg.Reg)
				ret = append(ret, x86.DoSwap(targetSpecs.Reg, argReg.Reg, argReg.Size))
				// Swap args in array
				regsArgsLocations[argNo], regsArgsLocations[regsArgsLocationsIndexesMap[targetSpecs]] = regsArgsLocations[regsArgsLocationsIndexesMap[targetSpecs]], regsArgsLocations[argNo]
			}
		} else {
			// We only do mov as it's safe
			ret = append(ret, x86.DoRegistryCopy(targetSpecs.Reg, argReg.Reg, argReg.Size))
		}
	}

	// Transfer args (memory)
	for memNo, argMem := range argsMem {
		ret = append(ret, x86.DoMemoryLoad(argMem.Index, argMem.Size, targetRegs[argsMemIndexes[memNo]].TopReg))
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
	for _, specs := range regsToPreserve {
		if _, ok := doNotPreserveReturn[specs.Reg]; !ok {
			ret = append(ret, x86.DoPop(specs.TopReg, 8))
		}
	}

	return ret
}
