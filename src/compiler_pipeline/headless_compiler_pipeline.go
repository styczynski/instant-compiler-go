package compiler_pipeline

import (
	"github.com/styczynski/latte-compiler/src/config"
	"github.com/styczynski/latte-compiler/src/input_reader"
	"github.com/styczynski/latte-compiler/src/printer"
)

type CompilationRequest struct {
	input string
}

type CompilationResponse struct {
	Summary string
	Ok      bool
}

func (in CompilationRequest) Read() ([]byte, error) {
	return []byte(in.input), nil
}

func (in CompilationRequest) Filename() string {
	return "<input>"
}

type HeadlessCompilerPipeline struct {
	inputRequests chan CompilationRequest
	outputs       chan CompilationResponse

	Printer *printer.LattePrinter

	running bool
}

func CreateHeadlessCompilerPipeline() *HeadlessCompilerPipeline {
	pr := printer.CreateLattePrinter()

	return &HeadlessCompilerPipeline{
		inputRequests: make(chan CompilationRequest),
		outputs:       make(chan CompilationResponse),
		Printer:       pr,
		running:       false,
	}
}

func (pipeline *HeadlessCompilerPipeline) ProcessAsync(input string, handler func(response CompilationResponse)) {
	if !pipeline.running {
		go func() {
			pipeline.running = true
			pipeline.deamonHandler()
		}()
	}

	pipeline.inputRequests <- CompilationRequest{
		input: input,
	}

	go func() {
		handler(<-pipeline.outputs)
	}()
}

func (pipeline *HeadlessCompilerPipeline) deamonHandler() {
	c := config.EntityConstConfig{
		Strings: map[string]string{
			"status-updater": "updater-silent",
			"backend":        "jvm",
			"runner":         "runner",
			"summary":        "summary-cli",
		},
	}
	for {
		request := <-pipeline.inputRequests

		p := config.CreateEntity(config.ENTITY_COMPILER_PIPELINE, "compiler-pipeline", c).(CompilerPipeline)
		message, ok := p.RunPipeline(c, input_reader.CreateLatteConstInputReader([]input_reader.LatteInput{request}))

		pipeline.outputs <- CompilationResponse{
			Summary: message,
			Ok:      ok,
		}

	}
}
