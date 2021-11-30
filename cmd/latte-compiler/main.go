package main

import (
	"os"
	"runtime/pprof"

	_ "github.com/styczynski/latte-compiler/src/compiler/jvm"
	"github.com/styczynski/latte-compiler/src/compiler_pipeline"
)

func main() {
	f, _ := os.Create("latc.profile")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	(compiler_pipeline.RunCompilerPipelineCliInterface{}).Run()
}
