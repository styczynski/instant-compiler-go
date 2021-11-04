package jasmine

import "fmt"

type IntOpType int64

const (
	Add IntOpType = iota
	Mul
	Sub
	Div
)

type JasmineIntOp struct {
	Operation IntOpType
}

func CreateJasmineIntOp(Operation string) *JasmineIntOp {
	if (Operation == "+") {
		return &JasmineIntOp{ Operation: Add, }
	}
	if (Operation == "-") {
		return &JasmineIntOp{ Operation: Sub, }
	}
	if (Operation == "*") {
		return &JasmineIntOp{ Operation: Mul, }
	}
	if (Operation == "/") {
		return &JasmineIntOp{ Operation: Div, }
	}
	panic(fmt.Sprintf("Unknown operation symbol in CreateJasmineIntOp(): %s", Operation))
}

func (p *JasmineIntOp) Type() JasmineInstructionType {
	return IntOp
}

func (p *JasmineIntOp) ToText(emitter EmitterConfig) string {
	if p.Operation == Add {
		return emitter.Emit("iadd")
	}
	if p.Operation == Sub {
		return emitter.Emit("isub")
	}
	if p.Operation == Mul {
		return emitter.Emit("imul")
	}
	if p.Operation == Div {
		return emitter.Emit("idiv")
	}
	panic(fmt.Sprintf("Invalid operation was used in context of JasmineIntOp: %d", p.Operation))
}

func (p *JasmineIntOp) StackSize(previousStackSize int) int {
	return previousStackSize
}
