package x86

import "fmt"

type CodePattern struct {
	Before  []Op
	Current Op
	After   []Op
}

func (pattern CodePattern) Matches(cursor CodeCursor) bool {
	for i, pat := range pattern.Before {
		instr := cursor.GetInstruction(-len(pattern.Before) + i)
		if instr.IsLabel() {
			return false
		}
		if !instr.IsLabel() && instr.Op != pat {
			return false
		}
	}
	instr := cursor.GetInstruction(0)
	if instr.IsLabel() {
		return false
	}
	if !instr.IsLabel() && instr.Op != pattern.Current {
		return false
	}
	for i, pat := range pattern.After {
		instr := cursor.GetInstruction(i + 1)
		if instr.IsLabel() {
			return false
		}
		if !instr.IsLabel() && instr.Op != pat {
			return false
		}
	}
	return true
}

func (pattern CodePattern) GetOffsetStart() int {
	return -len(pattern.Before)
}

func (pattern CodePattern) GetOffsetEnd() int {
	return len(pattern.After)
}

func (pattern CodePattern) Map(cursor CodeCursor, mapper func(before []*Instruction, current *Instruction, after []*Instruction) (bool, []*Instruction)) CodeCursor {
	for !cursor.IsOutside() {
		if pattern.Matches(cursor) {
			proceed, newSection := mapper(
				cursor.GetSection(pattern.GetOffsetStart(), -1),
				cursor.GetInstruction(0),
				cursor.GetSection(1, pattern.GetOffsetEnd()),
			)
			if proceed {
				fmt.Printf("REPLACE [%v] with [%v]\n", cursor.GetSection(pattern.GetOffsetStart(), pattern.GetOffsetEnd()), newSection)
				cursor.ReplaceSection(pattern.GetOffsetStart(), pattern.GetOffsetEnd(), newSection)
			}
		}
		cursor.Move(1)
	}
	return cursor
}

type CodeCursor interface {
	GetInstruction(offset int) *Instruction
	GetSection(fromOffeset int, toOffset int) []*Instruction
	ReplaceSection(fromOffeset int, toOffset int, newSection []*Instruction)
	Move(offset int)
	IsOutside() bool
	Reset()
	Dump() []*Instruction
}

type CodeCursorSlice struct {
	input    []*Instruction
	position int
}

func CreateCursor(input []*Instruction) CodeCursor {
	return &CodeCursorSlice{
		input:    input,
		position: 0,
	}
}

func (c *CodeCursorSlice) normalizedAbsoluteIndex(offset int) int {
	index := c.position + offset
	if index < 0 || index >= len(c.input) {
		return -1
	}
	return index
}

func (c *CodeCursorSlice) Reset() {
	c.position = 0
}

func (c *CodeCursorSlice) Dump() []*Instruction {
	return c.input
}

func (c *CodeCursorSlice) GetSection(fromOffset int, toOffset int) []*Instruction {
	if len(c.input) == 0 || fromOffset > toOffset {
		return []*Instruction{}
	}
	from := c.normalizedAbsoluteIndex(fromOffset)
	to := c.normalizedAbsoluteIndex(toOffset)
	if from == -1 {
		from = 0
	}
	if to == -1 {
		to = len(c.input) - 1
	}
	return c.input[from : to+1]
}

func (c *CodeCursorSlice) ReplaceSection(fromOffset int, toOffset int, newSection []*Instruction) {
	if len(c.input) == 0 || fromOffset > toOffset {
		return
	}
	from := c.normalizedAbsoluteIndex(fromOffset)
	to := c.normalizedAbsoluteIndex(toOffset)
	if from == -1 {
		from = 0
	}
	if to == -1 {
		to = len(c.input) - 1
	}
	prefix := []*Instruction{}
	if from > 0 {
		prefix = c.input[:from]
	}
	postfix := []*Instruction{}
	if to < len(c.input)-1 {
		postfix = c.input[to+1:]
	}
	prefix = append(prefix, newSection...)
	prefix = append(prefix, postfix...)
	c.input = prefix
}

func (c *CodeCursorSlice) GetInstruction(offset int) *Instruction {
	index := c.normalizedAbsoluteIndex(offset)
	if index == -1 {
		return DoNop()
	}
	return c.input[index]
}

func (c *CodeCursorSlice) Move(offset int) {
	c.position += offset
}

func (c *CodeCursorSlice) IsOutside() bool {
	index := c.normalizedAbsoluteIndex(0)
	if index == -1 {
		return true
	}
	return false
}

func inlineFunction(fn *Function, returnLabel string) []*Instruction {
	output := []*Instruction{}
	wasLabelUsed := false
	for i, stmt := range fn.Body {
		if !stmt.IsLabel() && stmt.Inst.Op == RET {
			if i == len(fn.Body)-1 {
				// We do not add return jump because it's the last statement in the inlined function
				break
			}
			wasLabelUsed = true
			output = append(output, DoJump(returnLabel))
			continue
		}
		output = append(output, stmt)
	}
	if wasLabelUsed {
		output = append(output, Label(returnLabel))
	}
	return output
}

func Optimize(code []Entry) []Entry {
	others := []Entry{}
	funcs := []*Function{}
	funcMap := map[string]*Function{}
	for _, entry := range code {
		if fn, ok := entry.(*Function); ok {
			funcs = append(funcs, fn)
			funcMap[fn.Label()] = fn
		} else {
			others = append(others, entry)
			//panic(fmt.Sprintf("Invalid optimization target: %s", reflect.TypeOf(entry).String()))
		}
	}

	// Filter unused functions
	usedLabels := map[string]int{}
	for _, fn := range funcs {
		cursor := CreateCursor(fn.Body)
		_, usedLabelsPartial := GetUnusedLabels(cursor)
		for label, count := range usedLabelsPartial {
			if count > 0 {
				if val, ok := usedLabels[label]; ok {
					usedLabels[label] = val + count
				} else {
					usedLabels[label] = count
				}
			}
		}
	}
	filteredFns := []*Function{}
	for _, fn := range funcs {
		if _, ok := usedLabels[fn.Label()]; ok || fn.Label() == "main" {
			// Function is used
			filteresInstrs := []*Instruction{}
			for i, instr := range fn.Body {
				if instr.IsLabel() {
					if _, ok := usedLabels[instr.Label]; !ok {
						continue
					}
				} else if instr.Op == CALL {
					label := instr.Args[0].(*RelLabel).label
					if calledFn, ok := funcMap[label]; ok {
						// Inline function
						backLabel := fmt.Sprintf("inline_call_%d_return_%s", i, fn.Label())
						inlinedContent := inlineFunction(calledFn, backLabel)
						filteresInstrs = append(filteresInstrs, inlinedContent...)
						continue
					}
				}
				filteresInstrs = append(filteresInstrs, instr)
			}
			fn.Body = filteresInstrs
			filteredFns = append(filteredFns, fn)
		}
	}
	funcs = filteredFns
	// End

	// Optimizations
	for _, fn := range funcs {
		cursor := CreateCursor(fn.Body)
		cursor = OptimizeRemoveDoubleMovs(cursor)
		cursor.Reset()
		fn.Body = cursor.Dump()
	}

	output := []Entry{}
	for _, other := range others {
		output = append(output, other)
	}
	for _, fn := range funcs {
		output = append(output, fn)
	}
	return output
}

func GetUnusedLabels(cursor CodeCursor) (CodeCursor, map[string]int) {
	usedLabels := map[string]int{}
	for !cursor.IsOutside() {
		instr := cursor.GetInstruction(0)
		if !instr.IsLabel() {
			for _, arg := range instr.Args {
				if l, ok := arg.(*RelLabel); ok {
					if val, ok := usedLabels[l.label]; ok {
						usedLabels[l.label] = val + 1
					} else {
						usedLabels[l.label] = 1
					}
				}
			}
		}
		cursor.Move(1)
	}
	return cursor, usedLabels
}

func OptimizeRemoveDoubleMovs(cursor CodeCursor) CodeCursor {
	p := CodePattern{
		Current: MOV,
		After: []Op{
			MOV,
		},
	}
	return p.Map(cursor, func(before []*Instruction, current *Instruction, after []*Instruction) (bool, []*Instruction) {
		fmt.Printf("Got MATCH: [%v] [%v] [%v]\n", before, current, after)
		if current.Args[0] == after[0].Args[1] {
			return true, []*Instruction{
				doRawMov(after[0].Args[0], current.Args[1], current.MemBytes),
				doRawMov(current.Args[0], current.Args[1], current.MemBytes),
			}
		}
		return false, nil
	})
}
