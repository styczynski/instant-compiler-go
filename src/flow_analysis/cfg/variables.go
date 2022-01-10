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
	for k, v := range vars {
		out[k] = v
	}
	return out
}

func (vars VariableSet) HasVariable(name string) bool {
	_, ok := vars[name]
	return ok
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
	GetAllVariables(visitedMap map[generic_ast.TraversableNode]struct{}) VariableSet
}

type NodeWithAssignedVariables interface {
	GetAssignedVariables(wantMembers bool, visitedMap map[generic_ast.TraversableNode]struct{}) VariableSet
}

type NodeWithDeclaredVariables interface {
	GetDeclaredVariables(visitedMap map[generic_ast.TraversableNode]struct{}) VariableSet
}

type NodeWithUsedVariables interface {
	GetUsedVariables(vars VariableSet, visitedMap map[generic_ast.TraversableNode]struct{}) VariableSet
}

type NodeWithVariableReplacement interface {
	RenameVariables(substUsed, substDecl VariableSubstitution)
}

type NodeWithRemovableVariableAsignment interface {
	RemoveVariableAssignment(variables map[string]struct{}) generic_ast.NormalNode
}

type VariableSubstitution interface {
	Get(name string) string
	Has(name string) bool
	Replace(name string) string
}

type VariableSubstitutionMap map[string]string

func (s VariableSubstitutionMap) Join(newMapping VariableSubstitutionMap) {
	for k, v := range newMapping {
		s[k] = v
	}
}

func (s VariableSubstitutionMap) Has(name string) bool {
	_, ok := s[name]
	return ok
}

func (s VariableSubstitutionMap) Copy() VariableSubstitutionMap {
	ret := VariableSubstitutionMap{}
	for k, v := range s {
		ret[k] = v
	}
	return ret
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
	value   generic_ast.NormalNode
}

func NewVariable(name string, value generic_ast.NormalNode) Variable {
	return &VariableImpl{
		varName: name,
		value:   value,
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

func GetAllUsagesVariables(node generic_ast.TraversableNode, visitedMap map[generic_ast.TraversableNode]struct{}) VariableSet {
	if _, wasVisited := visitedMap[node]; wasVisited {
		return NewVariableSet()
	}
	visitedMap[node] = struct{}{}
	if isNilNode(node) {
		//fmt.Printf("* NIL RET\n")
		return NewVariableSet()
	}
	vars := VariableSet{}
	for _, child := range node.GetChildren() {
		if child != nil {
			vars.Insert(GetAllUsagesVariables(child, visitedMap))
		}
	}
	if nodeWithUsedVariables, ok := node.(NodeWithUsedVariables); ok {
		r := nodeWithUsedVariables.GetUsedVariables(vars, visitedMap)
		//fmt.Printf("* RET %v\n", r)
		return r
	} else {
		//vars.Insert(GetAllVariables(node))
	}
	//fmt.Printf("* NORMAL RET %v\n", vars)
	return vars
}

func GetAllVariables(node generic_ast.TraversableNode, visitedMap map[generic_ast.TraversableNode]struct{}) VariableSet {
	if _, wasVisited := visitedMap[node]; wasVisited {
		return NewVariableSet()
	}
	visitedMap[node] = struct{}{}
	if isNilNode(node) {
		return NewVariableSet()
	}
	vars := VariableSet{}
	if nodeWithVariables, ok := node.(NodeWithVariables); ok {
		vars = nodeWithVariables.GetAllVariables(visitedMap)
	} else {
		for _, child := range node.GetChildren() {
			if child != nil {
				vars.Insert(GetAllVariables(child, visitedMap))
			}
		}
	}
	vars.Insert(GetAllAssignedVariables(node, true, map[generic_ast.TraversableNode]struct{}{}))
	vars.Insert(GetAllAssignedVariables(node, false, map[generic_ast.TraversableNode]struct{}{}))
	vars.Insert(GetAllDeclaredVariables(node, map[generic_ast.TraversableNode]struct{}{}))
	return vars
}

func isNilNode(node generic_ast.TraversableNode) bool {
	return node == nil || (reflect.ValueOf(node).Kind() == reflect.Ptr && reflect.ValueOf(node).IsNil())
}

func ReplaceVariables(node generic_ast.TraversableNode, substUsed VariableSubstitution, substDecl VariableSubstitution, visitedMap map[generic_ast.TraversableNode]struct{}) {
	if _, wasVisited := visitedMap[node]; wasVisited {
		return
	}
	visitedMap[node] = struct{}{}
	if isNilNode(node) {
		return
	}
	if nodeWithVariableReplacement, ok := node.(NodeWithVariableReplacement); ok {
		nodeWithVariableReplacement.RenameVariables(substUsed, substDecl)
		return
	}
	for _, child := range node.GetChildren() {
		if child != nil {
			ReplaceVariables(child, substUsed, substDecl, visitedMap)
		}
	}
}

func GetAllAssignedVariables(node generic_ast.TraversableNode, wantMembers bool, visitedMap map[generic_ast.TraversableNode]struct{}) VariableSet {
	if _, wasVisited := visitedMap[node]; wasVisited {
		return NewVariableSet()
	}
	visitedMap[node] = struct{}{}
	if isNilNode(node) {
		return NewVariableSet()
	}
	if nodeWithAssignedVariables, ok := node.(NodeWithAssignedVariables); ok {
		return nodeWithAssignedVariables.GetAssignedVariables(wantMembers, visitedMap)
	}
	vars := VariableSet{}
	for _, child := range node.GetChildren() {
		if child != nil {
			vars.Insert(GetAllAssignedVariables(child, wantMembers, visitedMap))
		}
	}
	return vars
}

func GetAllDeclaredVariables(node generic_ast.TraversableNode, visitedMap map[generic_ast.TraversableNode]struct{}) VariableSet {
	if _, wasVisited := visitedMap[node]; wasVisited {
		return NewVariableSet()
	}
	visitedMap[node] = struct{}{}
	if isNilNode(node) {
		return NewVariableSet()
	}
	if nodeWithDeclaredVariables, ok := node.(NodeWithDeclaredVariables); ok {
		return nodeWithDeclaredVariables.GetDeclaredVariables(visitedMap)
	}
	vars := VariableSet{}
	for _, child := range node.GetChildren() {
		if child != nil {
			vars.Insert(GetAllDeclaredVariables(child, visitedMap))
		}
	}
	return vars
}
