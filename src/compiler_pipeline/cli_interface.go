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
		Name:        "latc",
		Description: "Latte compiler written in Go",
		Version:     "1.0.0",
		Authors: []*cli.Author{
			{
				Name:  "Piotr Styczyński",
				Email: "piotr@styczynski.in",
			},
		},
		Copyright: "MIT License",
		Flags:     flags,
		Commands: []*cli.Command{
			{
				Name:    "shell",
				Aliases: []string{"sh"},
				Usage:   "Run interactive shell",
				Action: func(c *cli.Context) error {
					shellInterface := &RunCompilerPipelineInteractiveCliInterface{}
					shellInterface.Run()
					return nil
				},
			},
		},
		Action: func(c *cli.Context) error {
			p := config.CreateEntity(config.ENTITY_COMPILER_PIPELINE, "compiler-pipeline", c).(CompilerPipeline)
			message, _, ok := p.RunPipeline(c, input_reader.CreateLatteInputReader(c.Args().Slice()))
			if !ok {
				os.Stderr.WriteString("ERROR\n")
				fmt.Print(message)
				fmt.Print("\n")
				os.Exit(1)
			} else {
				os.Stderr.WriteString("OK\n")
				fmt.Print(message)
				fmt.Print("\n")
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
