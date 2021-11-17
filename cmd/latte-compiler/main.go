package main

import (
	_ "github.com/styczynski/latte-compiler/src/compiler/jvm"
	"github.com/styczynski/latte-compiler/src/compiler_pipeline"
	sat_solver "github.com/styczynski/latte-compiler/src/sat_solver/core"
)

func mainx() {
	(compiler_pipeline.RunCompilerPipelineCliInterface{}).Run()
}

func main() {
	sat_solver.GraphColouring()
}
