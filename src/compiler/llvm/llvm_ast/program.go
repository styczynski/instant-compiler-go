package llvm_ast

import (
	"fmt"
	"strings"

	"github.com/styczynski/latte-compiler/src/compiler"
)

type LLVMProgram struct {
	Instructions []LLVMInstruction
}

type LLVMInstructionType int64

const (
	PrintInt LLVMInstructionType = iota
	IntOp
	Val
	Program
	Var
)

type LLVMInstruction interface {
	Type() LLVMInstructionType
	ToText(emitter EmitterConfig) string
	GetDeclaredVariables() []string
	GetUsedVariables() []string
	ReplaceVariable(oldName string, newName string)
}

type LLVMTargetableInstruction interface {
	GetTarget(withType bool) string
	MoveTarget(newTarget string)
	IsMovable() bool
}

type LLVMValidableInstruction interface {
	Validate() *compiler.CompilationError
}

func (p *LLVMProgram) Type() LLVMInstructionType {
	return Program
}

func (p *LLVMProgram) ProgramToText() string {
	return p.ToText(EmitterConfig{
		Ident: 1,
	})
}

func (p *LLVMProgram) ToText(emitter EmitterConfig) string {
	return strings.Join(p.ToLines(emitter, true), "\n")
}

func (p *LLVMProgram) ToLines(emitter EmitterConfig, isTop bool) []string {
	output := []string{}
	for _, v := range p.Instructions {
		if v.Type() == Program {
			// Embedded program
			output = append(output, v.(*LLVMProgram).ToLines(emitter, false)...)
		} else {
			output = append(output, v.ToText(emitter))
		}
	}

	return strings.Split(fmt.Sprintf(`@dnl = internal constant [4 x i8] c"%%d\0A\00"
declare i32 @printf(i8*, ...)
declare i32 @puts(i8*)
define void @printInt(i32 %%x) {
   %%t0 = getelementptr [4 x i8], [4 x i8]* @dnl, i32 0, i32 0
   call i32 (i8*, ...) @printf(i8* %%t0, i32 %%x)
   ret void
}
define i32 @main() {
%s
	ret i32 0
}`, strings.Join(output, "\n")), "\n")
}

func (p *LLVMProgram) Validate() *compiler.CompilationError {
	for _, ins := range p.Instructions {
		if val, ok := ins.(LLVMValidableInstruction); ok {
			err := val.Validate()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *LLVMProgram) ReplaceVariable(oldName string, newName string) {
	panic("Operation not supported")
}

func (p *LLVMProgram) GetDeclaredVariables() []string {
	panic("Operation not supported")
}

func (p *LLVMProgram) GetUsedVariables() []string {
	panic("Operation not supported")
}

func (p *LLVMProgram) NormalizeVariables() {
	freeVar := 1
	varMap := map[string]string{}
	for _, stmt := range p.Instructions {
		allVars := stmt.GetUsedVariables()
		allVars = append(allVars, stmt.GetDeclaredVariables()...)
		for _, v := range allVars {
			if _, ok := varMap[v]; !ok {
				varMap[v] = fmt.Sprintf("%%%d", freeVar)
				freeVar++
			}
		}
		for _, v := range allVars {
			stmt.ReplaceVariable(v, varMap[v])
		}
	}
}
