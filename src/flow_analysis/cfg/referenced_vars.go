package cfg

import (
	"fmt"

	"github.com/styczynski/latte-compiler/src/generic_ast"
)

type ReferencedVars struct {
	asgt VariableSet
	updt VariableSet
	decl VariableSet
	use VariableSet
}

func (vars ReferencedVars) Print() string {
	return fmt.Sprintf("{%s; %s; %s; %s}", vars.asgt, vars.updt, vars.decl, vars.use)
}

func (c *CFG) ReferencedVars(node generic_ast.TraversableNode) ReferencedVars {
	return ReferencedVars{
		asgt: GetAllAssignedVariables(node, false, map[generic_ast.TraversableNode]struct{}{}),
		updt: GetAllAssignedVariables(node,true, map[generic_ast.TraversableNode]struct{}{}),
		decl: GetAllDeclaredVariables(node, map[generic_ast.TraversableNode]struct{}{}),
		use:  GetAllUsagesVariables(node, map[generic_ast.TraversableNode]struct{}{}),
	}
}
