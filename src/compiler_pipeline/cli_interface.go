package compiler_pipeline

import (
	"fmt"
	"log"
	"os"

	"github.com/styczynski/latte-compiler/src/config"
	"github.com/styczynski/latte-compiler/src/input_reader"
	"github.com/urfave/cli/v2"
)

type RunCompilerPipelineCliInterface struct {
}

func (RunCompilerPipelineCliInterface) Run() {
	flags := config.GetEntityParams()

	app := &cli.App{
		Flags: flags,
		Action: func(c *cli.Context) error {
			p := config.CreateEntity(config.ENTITY_COMPILER_PIPELINE, "compiler-pipeline", c).(CompilerPipeline)
			message, ok := p.RunPipeline(c, input_reader.CreateLatteInputReader(c.Args().Slice()))
			if !ok {
				os.Stderr.WriteString("ERROR\n")
				fmt.Print(message)
				os.Exit(1)
			} else {
				os.Stderr.WriteString("OK\n")
				fmt.Print(message)
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
