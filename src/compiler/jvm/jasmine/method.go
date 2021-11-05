package jasmine

import (
	"fmt"
	"strings"

	"github.com/styczynski/latte-compiler/src/compiler"
	pc "github.com/styczynski/latte-compiler/src/parser/context"
)

type JasmineMethod struct {
	Name        string
	Args        []string
	Returns     string
	StackLimit  int64
	LocalsLimit int64
	Body        []JasmineInstruction
}

func (p *JasmineMethod) Type() JasmineInstructionType {
	return Method
}

func (p *JasmineMethod) ToText(emitter EmitterConfig) string {
	subemitter := emitter.ApplyIdent(1)
	methodContents := []string{
		emitter.Emit(".method %s(%s)%s", p.Name, strings.Join(p.Args, ""), p.Returns),
		subemitter.Emit(".limit stack %d", p.StackLimit),
		subemitter.Emit(".limit locals %d", p.LocalsLimit),
	}
	s := 0
	for _, ins := range p.Body {
		s = ins.StackSize(s)
		methodContents = append(methodContents, ins.ToText(subemitter)) //+fmt.Sprintf(" // %d", s))
	}
	methodContents = append(methodContents, emitter.Emit(".end method"))
	return strings.Join(methodContents, "\n")
}

func (p *JasmineMethod) StackSize(previousStackSize int) int {
	s := 0
	for _, ins := range p.Body {
		s = ins.StackSize(s)
	}
	return s
}

func (p *JasmineMethod) Validate() *compiler.CompilationError {
	s := 0
	sMax := 0
	for _, stmt := range p.Body {
		s = stmt.StackSize(s)
		if s > sMax {
			sMax = s
		}
	}

	if sMax != int(p.StackLimit) {
		codeContext := pc.IndentCodeLines(p.ToText(EmitterConfig{Ident: 0}), 1, 0)
		return compiler.CreateCompilationError(
			"Internal assertion has failed",
			fmt.Sprintf("    | Expected stack limit: %d.\n    | Actual emitted stack limit: %d\n\n%s", sMax, p.StackLimit, codeContext))
	}

	vars := map[int64]struct{}{}
	for _, stmt := range p.Body {
		if stmt.Type() == LoadInt {
			vars[stmt.(*JasmineLoadInt).Index] = struct{}{}
		}
		if stmt.Type() == StoreInt {
			vars[stmt.(*JasmineStoreInt).Index] = struct{}{}
		}
		if stmt.Type() == ReferenceLoad {
			vars[stmt.(*JasmineReferenceLoad).Index] = struct{}{}
		}
	}

	localsCount := len(vars)
	specLimit := 0
	if p.Name == "<init>" {
		specLimit = 0
	}
	expectedLocalsLimit := localsCount + len(p.Args) + specLimit

	if int(p.LocalsLimit) != expectedLocalsLimit {
		codeContext := pc.IndentCodeLines(p.ToText(EmitterConfig{Ident: 0}), 2, 0)
		return compiler.CreateCompilationError(
			"Internal assertion has failed",
			fmt.Sprintf("    | Expected locals limit: %d.\n    | Actual emitted locals limit: %d\n\n%s", expectedLocalsLimit, p.LocalsLimit, codeContext))
	}
	return nil
}
