package compiler_pipeline

import (
	"github.com/styczynski/latte-compiler/src/config"
	"github.com/styczynski/latte-compiler/src/input_reader"
	"github.com/styczynski/latte-compiler/src/printer"
	"github.com/styczynski/latte-compiler/src/runner"
)

type CompilationRequest struct {
	input string
}

type CompilationResponse struct {
	Summary string
	Ok      bool
	Program runner.LatteRunnedProgram
}

func (in CompilationRequest) Read() ([]byte, error) {
	return []byte(in.input), nil
}

func (in CompilationRequest) Filename() string {
	return "input"
}

type HeadlessCompilerPipeline struct {
	inputRequests chan CompilationRequest
	outputs       chan CompilationResponse

	responseCache map[string]*CompilationResponse

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
		responseCache: map[string]*CompilationResponse{},
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
			"backend":        "x86",
			"runner":         "runner",
			"summary":        "summary-cli",
		},
		Ints: map[string]int{
			"summary-cli-error-limit": 900,
		},
		Bools: map[string]bool{
			"runner-always-run": true,
		},
	}
	for {
		request := <-pipeline.inputRequests

		if cachedResponse, ok := pipeline.responseCache[request.input]; false && ok {
			pipeline.responseCache = map[string]*CompilationResponse{}
			pipeline.outputs <- *cachedResponse
			return
		}
		// Flush cache
		pipeline.responseCache = map[string]*CompilationResponse{}

		p := config.CreateEntity(config.ENTITY_COMPILER_PIPELINE, "compiler-pipeline", c).(CompilerPipeline)
		message, progs, ok := p.RunPipeline(c, input_reader.CreateLatteConstInputReader([]input_reader.LatteInput{request}, input_reader.DEFAULT_RUNTIME_INCLUDES))

		response := &CompilationResponse{
			Summary: message,
			Ok:      ok,
			Program: progs[0],
		}

		pipeline.responseCache[request.input] = response
		pipeline.outputs <- *response
	}
}
