package jasmine

import (
	"github.com/styczynski/latte-compiler/src/compiler"
)

type JasmineStoreInt struct {
	Index compiler.Location
}

func (p *JasmineStoreInt) Type() JasmineInstructionType {
	return StoreInt
}

func (p *JasmineStoreInt) ToText(emitter EmitterConfig) string {
	if p.Index >= 0 && p.Index <= 3 {
		return emitter.Emit("istore_%d", p.Index)
	}
	return emitter.Emit("istore %d", p.Index)
}

func (p *JasmineStoreInt) StackSize(previousStackSize int) int {
	return previousStackSize - 1
}
