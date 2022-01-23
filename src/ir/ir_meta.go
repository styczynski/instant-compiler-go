package ir

import (
	"fmt"
	"strings"

	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
)

type IRMeta struct {
	TargetAllocation     IRAllocationMap
	ContextAllocation    IRAllocationMap
	AllocationContraints IRAllocationContraintsMap
}

func (meta IRMeta) String() string {
	return fmt.Sprintf("Meta{Target:%s, Context:%s}", meta.TargetAllocation.String(), meta.ContextAllocation.String())
}

type IRAllocation interface {
	String() string
	IsAllocation()
}

type IRAllocationConstraint interface {
	String() string
	IsAllocationConstraint()
}

type IRAllocationConstraints []IRAllocationConstraint

type IRAllocationMap map[string]IRAllocation
type IRAllocationContraintsMap map[string]IRAllocationConstraints

func (cons IRAllocationConstraints) String() string {
	descr := []string{}
	for _, con := range cons {
		descr = append(descr, con.String())
	}
	return fmt.Sprintf("Require[%s]", strings.Join(descr, ", "))
}

func (alloc IRAllocationMap) Copy() IRAllocationMap {
	newAlloc := IRAllocationMap{}
	for varName, loc := range alloc {
		newAlloc[varName] = loc
	}
	return newAlloc
}

func (alloc IRAllocationMap) String() string {
	descr := []string{}
	for varName, loc := range alloc {
		descr = append(descr, fmt.Sprintf("%s: %s", varName, loc.String()))
	}
	return fmt.Sprintf("Allocation(%s)", strings.Join(descr, ", "))
}

func (alloc IRAllocationMap) PreserveOnly(vars cfg.VariableSet) {
	allocCopy := alloc
	for k, _ := range allocCopy {
		if vars.HasVariable(k) {
			//newMap[k] = v
		} else {
			delete(alloc, k)
		}
	}
}

func (alloc IRAllocationContraintsMap) String() string {
	descr := []string{}
	for varName, loc := range alloc {
		descr = append(descr, fmt.Sprintf("%s: %s", varName, loc.String()))
	}
	return fmt.Sprintf("AllocConsMap(%s)", strings.Join(descr, ", "))
}

func (alloc IRAllocationContraintsMap) PreserveOnly(vars cfg.VariableSet) {
	allocCopy := alloc
	for k, _ := range allocCopy {
		if vars.HasVariable(k) {
			//newMap[k] = v
		} else {
			delete(alloc, k)
		}
	}
}
