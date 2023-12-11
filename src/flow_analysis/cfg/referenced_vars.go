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

	//fmt.Printf("REFA\n")
	asgt := GetAllAssignedVariables(node, false, map[generic_ast.TraversableNode]struct{}{})
	//fmt.Printf("REFB\n")
	updt := GetAllAssignedVariables(node, true, map[generic_ast.TraversableNode]struct{}{})
	//fmt.Printf("REFC\n")
	decl := GetAllDeclaredVariables(node, map[generic_ast.TraversableNode]struct{}{})
	//fmt.Printf("REFD %v\n", reflect.TypeOf(node))
	use := GetAllUsagesVariables(node, map[generic_ast.TraversableNode]struct{}{})
	//fmt.Printf("REFE\n")
	return ReferencedVars{
		asgt: asgt,
		updt: updt,
		decl: decl,
		use:  use,
	}
	
	//return r
}
