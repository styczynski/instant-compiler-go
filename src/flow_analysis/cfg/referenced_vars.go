package cfg

import (
	"fmt"

	"github.com/styczynski/latte-compiler/src/generic_ast"
)

type ReferencedVars struct {
	asgt VariableSet
	updt VariableSet
	decl VariableSet
	use  VariableSet
}

func (vars ReferencedVars) Assigned() VariableSet {
	return vars.asgt
}

func (vars ReferencedVars) All() VariableSet {
	ret := NewVariableSet()
	ret.Insert(vars.asgt)
	ret.Insert(vars.decl)
	ret.Insert(vars.updt)
	ret.Insert(vars.use)
	return ret
}

func (vars ReferencedVars) Updated() VariableSet {
	return vars.updt
}

func (vars ReferencedVars) Declared() VariableSet {
	return vars.decl
}

func (vars ReferencedVars) Used() VariableSet {
	return vars.use
}

func (vars ReferencedVars) Print() string {
	return fmt.Sprintf("as=%s up=%s de=%s us=%s", vars.asgt, vars.updt, vars.decl, vars.use)
}

func (c *CFG) ReferencedVars(node generic_ast.TraversableNode) ReferencedVars {
	return ReferencedVars{
		asgt: GetAllAssignedVariables(node, false, map[generic_ast.TraversableNode]struct{}{}),
		updt: GetAllAssignedVariables(node, true, map[generic_ast.TraversableNode]struct{}{}),
		decl: GetAllDeclaredVariables(node, map[generic_ast.TraversableNode]struct{}{}),
		use:  GetAllUsagesVariables(node, map[generic_ast.TraversableNode]struct{}{}),
	}
	//fmt.Printf("REF %s\n", r.Print())
	//return r
}
