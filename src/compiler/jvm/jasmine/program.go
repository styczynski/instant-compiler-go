package jasmine

import (
	"strings"

	"github.com/styczynski/latte-compiler/src/compiler"
)

type JasmineProgram struct {
	Instructions []JasmineInstruction
}

type JasmineInstructionType int64

const (
	Push JasmineInstructionType = iota
	StoreInt
	ConstInt
	LoadInt
	InvokeStatic
	Return
	Pop
	Swap
	IntOp
	Program
	ReferenceLoad
	Method
	Class
)

type JasmineInstruction interface {
	Type() JasmineInstructionType
	ToText(emitter EmitterConfig) string
	StackSize(previousStackSize int) int
}

type JasmineValidableInstruction interface {
	Validate() *compiler.CompilationError 
}

func (p *JasmineProgram) Type() JasmineInstructionType {
	return Program
}

func (p *JasmineProgram) ProgramToText() string {
	return p.ToText(EmitterConfig{
		Ident: 0,
	})
}

func (p *JasmineProgram) ToText(emitter EmitterConfig) string {
	return strings.Join(p.ToLines(emitter), "\n")
}

func (p *JasmineProgram) ToLines(emitter EmitterConfig) []string {
	output := []string{}
	for _, v := range p.Instructions {
		if v.Type() == Method {

		} else if v.Type() == Program {
			// Embedded program
			output = append(output, v.(*JasmineProgram).ToLines(emitter)...)
		} else {
			output = append(output, v.ToText(emitter))
		}
	}
	return output
}

func (p *JasmineProgram) StackSize(previousStackSize int) int {
	for _, v := range p.Instructions {
		previousStackSize = v.StackSize(previousStackSize)
	}
	return previousStackSize
}

func (p *JasmineProgram) Validate() *compiler.CompilationError {
	for _, ins := range p.Instructions {
		if val, ok := ins.(JasmineValidableInstruction); ok {
			err := val.Validate()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
