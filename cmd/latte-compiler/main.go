package main

import (
	_ "github.com/styczynski/latte-compiler/src/compiler/jvm"
	_ "github.com/styczynski/latte-compiler/src/compiler/llvm"
	"github.com/styczynski/latte-compiler/src/compiler_pipeline"
)

func main() {
	(compiler_pipeline.RunCompilerPipelineCliInterface{}).Run()
}
