package x86

import (
	"fmt"
	"strings"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"golang.org/x/arch/x86/x86asm"
)

type GenerationContext struct {
	pc     uint64
	indent int
}

func (c *GenerationContext) ShiftPC(distance uint64) {
	c.pc = c.pc + distance
}

func (c *GenerationContext) Indent() {
	c.indent = c.indent + 1
}

func (c *GenerationContext) IndentBack() {
	c.indent = c.indent - 1
}

type Entry interface {
	Generate(c *GenerationContext) []string
}

type Function struct {
	Name   string
	Source generic_ast.BaseASTNode
	Body   []*Instruction
}

func (f *Function) Generate(c *GenerationContext) []string {
	retInstrs := []string{}
	headers := []string{
		fmt.Sprintf("; Function %s", f.Name),
		fmt.Sprintf("; Source: %s", f.Source.Pos.String()),
		fmt.Sprintf("_%s:", f.Name),
	}
	footers := []string{
		fmt.Sprintf("; End of function %s", f.Name),
	}
	c.Indent()
	retInstrs = append(retInstrs, headers...)
	for _, instr := range f.Body {
		gen := instr.Generate(c)
		retInstrs = append(retInstrs, gen...)
	}
	c.IndentBack()
	retInstrs = append(retInstrs, footers...)
	return retInstrs
}

type Instruction struct {
	x86asm.Inst
}

func (f *Instruction) Generate(c *GenerationContext) []string {
	prefix := strings.Repeat("  ", c.indent)
	ret := []string{
		prefix + x86asm.IntelSyntax(f.Inst, c.pc, nil),
	}
	c.ShiftPC(1)
	return ret
}

type Program struct {
	Entries []Entry
}

func (p *Program) ProgramToText() string {
	c := &GenerationContext{
		pc:     uint64(0),
		indent: 0,
	}
	code := []string{}
	for _, entry := range p.Entries {
		codeChunk := entry.Generate(c)
		code = append(code, codeChunk...)
	}
	return strings.Join(code, "\n")
}
