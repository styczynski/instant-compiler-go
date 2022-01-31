package allocation

import (
	"fmt"

	"github.com/styczynski/latte-compiler/src/compiler/assembly/x86"
	"github.com/styczynski/latte-compiler/src/ir"
)

func getFreeStackTopAlloc(
	allocationOutputContext ir.IRAllocationMap,
) (*LocationMemory, int) {
	allocator := &LinearScanAllocator{
		state: &LinearScanAllocatorState{
			Current: allocationOutputContext.Copy(),
			All:     allocationOutputContext.Copy(),
		},
	}
	return allocator.state.LastMemoryBlock()
}

func allocateRegistryStorage(
	parentFn *ir.IRFunction,
	allocationOutputContext ir.IRAllocationMap,
	sizes []int,
) ([]*LocationMemory, AssemblyFunctionMeta) {
	allocator := &LinearScanAllocator{
		state: &LinearScanAllocatorState{
			Current: allocationOutputContext.Copy(),
			All:     allocationOutputContext.Copy(),
		},
	}
	results := []*LocationMemory{}
	for i, blockSize := range sizes {
		loc := allocator.state.allocateAvailableBlock(blockSize)
		results = append(results, loc)
		allocator.state.Current[fmt.Sprintf("reg_%d", i)] = loc
	}
	meta := allocator.state.GetFunctionMeta(parentFn).(AssemblyFunctionMeta)
	return results, meta
}

func DoCall(
	parentFn *ir.IRFunction,
	label string,
	targetType ir.IRType,
	targetAlloc ir.IRAllocation,
	argsOrder []string,
	argsAllocs ir.IRAllocationMap,
	allocationContext ir.IRAllocationMap,
	allocationOutputContext ir.IRAllocationMap,
	registriesOverrides map[x86.Reg]int64,
	transferResult bool,
	preserveAnything bool,
) (AssemblyFunctionMeta, []*x86.Instruction) {

	fmt.Printf("OUTPUT CONTEXT FOR %s IS %v\n", label, allocationOutputContext)
	usedRegs := []*x86.RegistrySpecs{}
	for _, alloc := range allocationOutputContext {
		if allocReg, ok := IsAllocReg(alloc); ok {
			usedRegs = append(usedRegs, allocReg.State.Specs)
		}
	}

	ret := []*x86.Instruction{}

	targetRegs := []*x86.RegistrySpecs{}
	targetRegsSet := map[x86.Reg]*x86.RegistrySpecs{}
	targetMemsOffsets := []int{}
	firstMemArgNo := -1
	for argNo, _ := range argsOrder {
		specs := x86.GetRegisterForFunctionArg(argNo)
		if specs == nil && firstMemArgNo == -1 {
			firstMemArgNo = argNo
			break
		}
	}
	for argNo, _ := range argsOrder {
		specs := x86.GetRegisterForFunctionArg(argNo)
		if specs != nil {
			targetRegs = append(targetRegs, specs)
			targetRegsSet[specs.Normalized] = specs
			targetMemsOffsets = append(targetMemsOffsets, -1)
		} else {
			// Transfer to memory
			offset, _ := x86.GetMemoryForFunctionArg(firstMemArgNo + len(argsOrder) - 1 - argNo)
			targetMemsOffsets = append(targetMemsOffsets, offset)
		}
	}

	regsToPreserve := map[x86.Reg]*x86.RegistrySpecs{}

	regsArgsLocations := []*x86.RegistrySpecs{}
	regsArgsLocationsIndexes := []int{}

	swapRegistriesMap := map[*x86.RegistrySpecs]*x86.RegistrySpecs{}
	for _, reg := range x86.ALL_REGS {
		swapRegistriesMap[reg] = reg
	}

	regsArgsLocationsSet := map[x86.Reg]struct{}{}
	// regsArgsLocationsIndexesMap := map[*x86.RegistrySpecs]int{}

	argsMem := []*LocationMemory{}
	argsMemIndexes := []int{}

	// Determine registry overrides causing new preserved registries
	if preserveAnything {
		for target, _ := range registriesOverrides {
			_, mappedTarget := x86.ResizeReg(target, 4)
			regsToPreserve[mappedTarget] = x86.ALL_REGS[mappedTarget]
		}

		// Determine registry to preserve
		for _, usedReg := range usedRegs {
			if true || !usedReg.IsPreserved {
				fmt.Printf("REQUEST PRESERVE REGISTRY!: %s %v\n", label, usedReg.Normalized)
				regsToPreserve[usedReg.Normalized] = usedReg
			}
		}
	}

	fmt.Printf("DOCALL FOR %s WITH OUTPUT ALLOC %s\n", label, allocationOutputContext.String())

	for argNo, argName := range argsOrder {
		arg := argsAllocs[argName]
		fmt.Printf("ARGUMENT %d FOR CALL IS ALLOCATED ON: %v\n", argNo, arg)
		if allocReg, ok := IsAllocReg(arg); ok {
			regsArgsLocations = append(regsArgsLocations, allocReg.State.Specs)
			regsArgsLocationsIndexes = append(regsArgsLocationsIndexes, argNo)
			regsArgsLocationsSet[allocReg.State.Specs.Normalized] = struct{}{}
			// regsArgsLocationsIndexesMap[allocReg.State.Specs] = argNo

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
			// if , ok := allocationOutputContext[argName]; ok && preserveAnything {
			// 	// We need to preserve input registry argument
			// 	regsToPreserve[allocReg.Reg] = allocReg.State.Specs
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
	if targetType != ir.IR_VOID && transferResult {
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
	regsSizes := []int{}
	for _, specs := range regsToPreserve {
		if _, ok := doNotPreserveReturn[specs.Normalized]; !ok {
			// Allocate space for the preserved registry
			regsSizes = append(regsSizes, 8)
		}
	}
	currentMeta := parentFn.GetMeta().(AssemblyFunctionMeta)
	registryStorage, meta := allocateRegistryStorage(parentFn, allocationContext, regsSizes)

	if meta.VarLen > currentMeta.VarLen {
		currentMeta.VarLen = meta.VarLen
	}

	registryStorageIndex := 0

	if preserveAnything {
		for _, specs := range regsToPreserve {
			fmt.Printf("CHECK RETURN REG SHOULD PRESERVED BE?: %s %v\n", label, specs.Normalized)
			if _, ok := doNotPreserveReturn[specs.Normalized]; !ok {
				fmt.Printf("   PROCEED I CHUJ!\n")
				regsPreserveOrder = append(regsPreserveOrder, specs.Reg8B)
				//ret = append(ret, x86.DoPush(specs.Reg8B, 8))
				regStorage := registryStorage[registryStorageIndex]
				ret = append(ret, x86.DoMemoryStore(regStorage.Index, regStorage.Size, specs.Reg8B))
				registryStorageIndex++
			}
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

		targetRegArgNo := regsArgsLocationsIndexes[argNo]
		if targetRegArgNo >= len(targetRegs) {
			// We require tranfser to memory
			targetOffset := targetMemsOffsets[argNo]
			_, stackTop := getFreeStackTopAlloc(allocationContext)
			blockAllocSize := 4
			blockAllocIndex := targetOffset + stackTop
			fmt.Printf("ARGUMENT %d FOR CALL WILL BE TRANSFERED TO MEMORY WITH OFFSET: %d\n", argNo, targetOffset)

			srcReg := argReg

			ret = append(ret, x86.DoMemoryStore(blockAllocIndex, blockAllocSize, srcReg.Reg4B))

			newEnd := blockAllocIndex + blockAllocSize
			if newEnd > currentMeta.VarLen {
				currentMeta.VarLen = newEnd
			}

			continue
		}

		targetSpecs := targetRegs[regsArgsLocationsIndexes[argNo]]
		// Registry not used if teraget has override
		if _, hasOverride := registriesOverrides[targetSpecs.Normalized]; hasOverride {
			continue
		}
		fmt.Printf("ARGUMENT %d FOR CALL WILL BE TRANSFERED TO: %v\n", argNo, targetSpecs.Normalized)

		if _, ok := regsArgsLocationsSet[argReg.Normalized]; ok {
			// We have this register as an input so we do swap
			if targetSpecs.Reg8B != argReg.Reg8B {
				swapA := targetSpecs
				swapB := argReg

				fmt.Printf("ARGUMENT SWAP %v %v\n", swapA.Normalized, swapB.Normalized)
				ret = append(ret, x86.DoSwap(swapA.Normalized, swapB.Normalized, argReg.DefaultSize))
				// Swap args in array
				//regsArgsLocations[argNo], regsArgsLocations[regsArgsLocationsIndexesMap[targetSpecs]] = regsArgsLocations[regsArgsLocationsIndexesMap[targetSpecs]], regsArgsLocations[argNo]

				swapRegistriesMap[swapA], swapRegistriesMap[swapB] = swapRegistriesMap[swapB], swapRegistriesMap[swapA]
				regsArgsLocationsSet = map[x86.Reg]struct{}{}
				for i, reg := range regsArgsLocations {
					if reg == swapA {
						regsArgsLocations[i] = swapB
					} else if reg == swapB {
						regsArgsLocations[i] = swapA
					}
					regsArgsLocationsSet[regsArgsLocations[i].Normalized] = struct{}{}
				}

				// fmt.Printf("CALL ARGS ==> [")
				// for argNoI := 0; argNoI < len(argsOrder); argNoI++ {
				// 	fmt.Printf(" {%d: %v} ", argNoI, regsArgsLocations[argNoI].Normalized)
				// }
				// fmt.Printf("]\n")
			}
		} else {
			// We only do mov as it's safe
			ret = append(ret, x86.DoRegistryCopy(targetSpecs.Normalized, argReg.Normalized, argReg.DefaultSize))
		}
	}

	// Transfer args (memory)
	for memNo, argMem := range argsMem {
		tgtReg := targetRegs[argsMemIndexes[memNo]]
		ret = append(ret, x86.DoMemoryLoad(argMem.Index, argMem.Size, tgtReg.Normalized))
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
	if transferResult {
		ret = append(ret, returnTransferInstrs...)
	}

	// We need retrieve values of the registers
	if preserveAnything {
		registryStorageIndex = 0
		for _, reg := range regsPreserveOrder {
			regStorage := registryStorage[registryStorageIndex]
			ret = append(ret, x86.DoMemoryLoad(regStorage.Index, regStorage.Size, reg))
			registryStorageIndex++
		}
	}

	return currentMeta, ret
}
