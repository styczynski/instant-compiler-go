package cfg

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/styczynski/latte-compiler/src/generic_ast"
)

type VariableSet map[string]Variable

func (vars VariableSet) Copy() VariableSet {
	out := VariableSet{}
	for k,v := range vars {
		out[k] = v
	}
	return out
}

func (vars VariableSet) ReplaceBlock(old generic_ast.NormalNode, new generic_ast.NormalNode) {
	for name, varDef := range vars {
		if varDef.Value() == old {
			vars[name] = NewVariable(varDef.Name(), new)
		}
	}
}

func (vars VariableSet) String() string {
	varStrs := []string{}
	for _, varInfo := range vars {
		varStrs = append(varStrs, varInfo.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(varStrs, ", "))
}

func NewVariableSet(vars ...Variable) VariableSet {
	ret := VariableSet{}
	for _, v := range vars {
		ret[v.Name()] = v
	}
	return ret
}

type NodeWithVariables interface {
	GetAllVariables() VariableSet
}

type NodeWithAssignedVariables interface {
	GetAssignedVariables(wantMembers bool) VariableSet
}

type NodeWithDeclaredVariables interface {
	GetDeclaredVariables() VariableSet
}

type NodeWithUsedVariables interface {
	GetUsedVariables(vars VariableSet) VariableSet
}

type NodeWithVariableReplacement interface {
	RenameVariables(subst VariableSubstitution)
}

type VariableSubstitution interface {
	Get(name string) string
	Has(name string) bool
	Replace(name string) string
}

type VariableSubstitutionMap map[string]string

func (s VariableSubstitutionMap) Has(name string) bool {
	_, ok := s[name]
	return ok
}

func (s VariableSubstitutionMap) Replace(name string) string {
	v, ok := s[name]
	if !ok {
		return name
	}
	return v
}

func (s VariableSubstitutionMap) Get(name string) string {
	v, ok := s[name]
	if !ok {
		return ""
	}
	return v
}

type Variable interface {
	 Name() string
	 String() string
	 Value() generic_ast.NormalNode
}

type VariableImpl struct {
	varName string
	value generic_ast.NormalNode
}

func NewVariable(name string, value generic_ast.NormalNode) Variable {
	return &VariableImpl{
		varName: name,
		value: value,
	}
}

func (v *VariableImpl) Value() generic_ast.NormalNode {
	return v.value
}

func (v *VariableImpl) String() string {
	return v.Name()
}

func (v *VariableImpl) Name() string {
	return v.varName
}

func (v VariableSet) Add(varInfo Variable) {
	v[varInfo.Name()] = varInfo
}

func (v VariableSet) Insert(src VariableSet) {
	for name, varValue := range src {
		v[name] = varValue
	}
}

func GetAllUsagesVariables(node generic_ast.TraversableNode) VariableSet {
	if isNilNode(node) {
		return NewVariableSet()
	}
	vars := VariableSet{}
	for _, child := range node.GetChildren() {
		if child != nil {
			vars.Insert(GetAllUsagesVariables(child))
		}
	}
	if nodeWithUsedVariables, ok := node.(NodeWithUsedVariables); ok {
		return nodeWithUsedVariables.GetUsedVariables(vars)
	} else {
		//vars.Insert(GetAllVariables(node))
	}
	return vars
}

func GetAllVariables(node generic_ast.TraversableNode) VariableSet {
	if isNilNode(node) {
		return NewVariableSet()
	}
	vars := VariableSet{}
	if nodeWithVariables, ok := node.(NodeWithVariables); ok {
		vars = nodeWithVariables.GetAllVariables()
	} else {
		for _, child := range node.GetChildren() {
			if child != nil {
				vars.Insert(GetAllVariables(child))
			}
		}
	}
	vars.Insert(GetAllAssignedVariables(node, true))
	vars.Insert(GetAllAssignedVariables(node, false))
	vars.Insert(GetAllDeclaredVariables(node))
	return vars
}

func isNilNode(node generic_ast.TraversableNode) bool {
	return node == nil || (reflect.ValueOf(node).Kind() == reflect.Ptr && reflect.ValueOf(node).IsNil())
}

func GetAllAssignedVariables(node generic_ast.TraversableNode, wantMembers bool) VariableSet {
	if isNilNode(node) {
		return NewVariableSet()
	}
	if nodeWithAssignedVariables, ok := node.(NodeWithAssignedVariables); ok {
		return nodeWithAssignedVariables.GetAssignedVariables(wantMembers)
	}
	vars := VariableSet{}
	for _, child := range node.GetChildren() {
		if child != nil {
			vars.Insert(GetAllAssignedVariables(child, wantMembers))
		}
	}
	return vars
}

func GetAllDeclaredVariables(node generic_ast.TraversableNode) VariableSet {
	if isNilNode(node) {
		return NewVariableSet()
	}
	if nodeWithDeclaredVariables, ok := node.(NodeWithDeclaredVariables); ok {
		return nodeWithDeclaredVariables.GetDeclaredVariables()
	}
	vars := VariableSet{}
	for _, child := range node.GetChildren() {
		if child != nil {
			vars.Insert(GetAllDeclaredVariables(child))
		}
	}
	return vars
}
