package jasmine

import (
	"strings"
)

type JasmineProgram struct {
	StackLimit   int64
	LocalsLimit  int64
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

func (p *JasmineProgram) Type() JasmineInstructionType {
	return Program
}

func (p *JasmineProgram) ProgramToText() string {
	return p.ToText(EmitterConfig{
		Ident: 0,
	})
}

func (p *JasmineProgram) ToText(emitter EmitterConfig) string {
	return strings.Join(p.ToLines(emitter, true), "\n")
}

func (p *JasmineProgram) ToLines(emitter EmitterConfig, isTop bool) []string {
	output := []string{}
	if isTop {
		output = append(output, []string{
			emitter.Emit(".limit stack %d", p.StackLimit),
			emitter.Emit(".limit locals %d", p.LocalsLimit),
		}...)
	}
	for _, v := range p.Instructions {
		if v.Type() == Method {

		} else if v.Type() == Program {
			// Embedded program
			output = append(output, v.(*JasmineProgram).ToLines(emitter, false)...)
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
