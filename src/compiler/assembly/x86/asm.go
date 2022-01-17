package x86

import (
	"fmt"
	"strings"

	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/ir"
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
	Generate(c *GenerationContext, slFn SymLookup) []string
	GenerateSymbolLookup(c *GenerationContext, sl *SymbolLookup)
}

type Function struct {
	Name   string
	Source generic_ast.BaseASTNode
	Body   []*Instruction
}

func (f *Function) GenerateSymbolLookup(c *GenerationContext, sl *SymbolLookup) {
	for _, instr := range f.Body {
		instr.GenerateSymbolLookup(c, sl)
	}
}

func (f *Function) Generate(c *GenerationContext, slFn SymLookup) []string {
	retInstrs := []string{}
	additionalDescription := ""
	fnLabel := fmt.Sprintf("_%s:", f.Name)
	if f.Name == "main" {
		fnLabel = "main:"
		additionalDescription = " (Entrypoint)"
	}
	headers := []string{
		fmt.Sprintf("# Function %s%s", f.Name, additionalDescription),
		fmt.Sprintf("# Source: %s", f.Source.Pos.String()),
		fnLabel,
	}
	footers := []string{
		fmt.Sprintf("# End of function %s", f.Name),
	}
	c.Indent()
	c.Indent()
	retInstrs = append(retInstrs, headers...)
	for _, instr := range f.Body {
		gen := instr.Generate(c, slFn)
		retInstrs = append(retInstrs, gen...)
	}
	c.IndentBack()
	retInstrs = append(retInstrs, footers...)
	return retInstrs
}

type Instruction struct {
	Inst
	Label   string
	Comment string
}

type SymbolLookup struct {
	table   map[string]uint64
	reverse map[uint64]string
}

func (sl *SymbolLookup) GetLookupFunction() SymLookup {
	return func(addr uint64) (string, uint64) {
		if name, ok := sl.reverse[addr]; ok {
			return name, sl.table[name]
		}
		return "", 0
	}
}

func (f *Instruction) GenerateSymbolLookup(c *GenerationContext, sl *SymbolLookup) {
	if len(f.Label) > 0 {
		sl.table[f.Label] = c.pc
		sl.reverse[c.pc] = f.Label
	} else {
		c.ShiftPC(1)
	}
}

func (f *Instruction) FromIR(originalIR *ir.IRStatement) *Instruction {
	if originalIR != nil {
		f.Comment = originalIR.Comment
	}

	if len(originalIR.BaseASTNode.Begin().Filename) > 0 {
		f.Comment = fmt.Sprintf("%s (%s)", f.Comment, originalIR.BaseASTNode.Begin())
	}

	return f
}

func (f *Instruction) Generate(c *GenerationContext, slFn SymLookup) []string {
	commentStr := ""
	if len(f.Comment) > 0 {
		commentStr = fmt.Sprintf(" # %s", f.Comment)
	}
	if len(f.Label) > 0 {
		c.IndentBack()
		ret := []string{
			fmt.Sprintf("%s%s:%s", strings.Repeat("  ", c.indent), f.Label, commentStr),
		}
		c.Indent()
		return ret
	}
	prefix := strings.Repeat("  ", c.indent)
	//f.Inst.MemBytes = 0
	ret := []string{
		prefix + GNUSyntax(f.Inst, c.pc, slFn) + commentStr,
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
	sl := &SymbolLookup{
		table:   map[string]uint64{},
		reverse: map[uint64]string{},
	}
	// for _, entry := range p.Entries {
	// 	entry.GenerateSymbolLookup(c, sl)
	// }

	c = &GenerationContext{
		pc:     uint64(0),
		indent: 0,
	}

	code := []string{
		".text",
		".global main",
	}
	for _, entry := range p.Entries {
		codeChunk := entry.Generate(c, sl.GetLookupFunction())
		code = append(code, codeChunk...)
	}
	return strings.Join(code, "\n")
}
