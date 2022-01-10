package main

import (
	_ "github.com/styczynski/latte-compiler/src/compiler/assembly"
	"github.com/styczynski/latte-compiler/src/compiler_pipeline"
)

func main() {
	(compiler_pipeline.RunCompilerPipelineCliInterface{}).Run()
}
