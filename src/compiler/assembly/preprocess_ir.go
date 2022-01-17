package assembly

import (
	"github.com/styczynski/latte-compiler/src/ir"
	"github.com/styczynski/latte-compiler/src/parser/context"
)

func (backend CompilerX86Backend) preprocessIR(c *context.ParsingContext, code *ir.IRProgram) error {
	for _, fn := range code.Statements {
		for _, fnBlock := range fn.FunctionBody {
			err := backend.preprocessIRBlock(c, fn, fnBlock)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
