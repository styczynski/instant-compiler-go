package ir

import (
	"fmt"
	"strings"

	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
)

type IRMeta struct {
	Allocation IRAllocationMap
}

func (meta IRMeta) String() string {
	return fmt.Sprintf("Meta{%s}", meta.Allocation.String())
}

type IRAllocation interface {
	String() string
}

type IRAllocationMap map[string]IRAllocation

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
