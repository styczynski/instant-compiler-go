package compiler_pipeline

import (
	"log"
	"os"

	"github.com/styczynski/latte-compiler/src/config"
	"github.com/urfave/cli/v2"
)

type RunCompilerPipelineCliInterface struct {
}

func (RunCompilerPipelineCliInterface) Run() {
	flags := config.GetEntityParams()

	app := &cli.App{
		Flags: flags,
		Action: func(c *cli.Context) error {
			pipeline := config.CreateEntity(config.ENTITY_COMPILER_PIPELINE, "compiler-pipeline", c).(CompilerPipeline)
			return pipeline.RunPipeline(c, c.Args().Slice())
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
