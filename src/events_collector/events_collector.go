package events_collector

import (
	"strings"
	"sync"
	"time"

	"github.com/styczynski/latte-compiler/src/compiler"
	"github.com/styczynski/latte-compiler/src/flow_analysis"
	"github.com/styczynski/latte-compiler/src/parser"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/runner"
	"github.com/styczynski/latte-compiler/src/type_checker"
)

type EventsCollector struct {
	eventStream        chan EventMessage
	done               chan bool
	statuses           map[string][]InputStatus
	timingsAggregation map[string]time.Duration
	timingsLabels      []string
	updater            StatusUpdater
	outputFiles        map[string]map[string]string
	outputFilsMutex    sync.Mutex
}

type InputStatus struct {
	processName string
	c           *context.ParsingContext
	input       context.EventCollectorMessageInput
	start       time.Time
	end         time.Time
}

type EventMessage struct {
	eventType   string
	processName string
	c           *context.ParsingContext
	input       context.EventCollectorMessageInput
}

func (collector *EventsCollector) EmitOutputFiles(processName string, c *context.ParsingContext, outputFiles map[string]map[string]string) {
	collector.outputFilsMutex.Lock()
	defer collector.outputFilsMutex.Unlock()
	for path, files := range outputFiles {
		if _, ok := collector.outputFiles[path]; !ok {
			collector.outputFiles[path] = map[string]string{}
		}
		for file, description := range files {
			collector.outputFiles[path][file] = description
		}
	}
}

func (collector *EventsCollector) Start(processName string, c *context.ParsingContext, input context.EventCollectorMessageInput) {
	collector.eventStream <- EventMessage{
		eventType:   "start",
		processName: processName,
		c:           c,
		input:       input,
	}
}

func (collector *EventsCollector) End(processName string, c *context.ParsingContext, input context.EventCollectorMessageInput) {
	collector.eventStream <- EventMessage{
		eventType:   "end",
		processName: processName,
		c:           c,
		input:       input,
	}
}

func maxDuration(t1 time.Duration, t2 time.Duration) time.Duration {
	if t1 > t2 {
		return t1
	}
	return t2
}

func (collector *EventsCollector) insertTimeAggregation(ids []string, t time.Duration) {
	id := strings.Join(ids, "|")
	if _, ok := collector.timingsAggregation[id]; ok {
		collector.timingsAggregation[id] = maxDuration(collector.timingsAggregation[id], t)
	} else {
		collector.timingsAggregation[id] = t
		collector.timingsLabels = append(collector.timingsLabels, id)
	}
}

func (collector *EventsCollector) debugOutputEvent(format string, a ...interface{}) {
	// Enable on debug
	//fmt.Printf(format, a...)
}

func (collector *EventsCollector) runEventsCollectorDeamon() {
	eventStream := collector.eventStream
	go func() {
		defer close(eventStream)
		gracefulShutdown := false
		updaterDeinitialized := true
		for {
			message := <-eventStream
			switch message.eventType {
			case "start":
				if updaterDeinitialized {
					collector.updater.Init()
					updaterDeinitialized = false
				}
				filename := message.input.Filename()
				collector.statuses[filename] = append(collector.statuses[filename], InputStatus{
					processName: message.processName,
					c:           message.c,
					input:       message.input,
					start:       time.Now(),
				})
				collector.updater.UpdateStatus(message.processName)
				collector.debugOutputEvent("[Event] %s - Start %s\n", filename, message.processName)
			case "end":
				filename := message.input.Filename()
				statusLen := len(collector.statuses[filename])
				topStatus := collector.statuses[filename][statusLen-1]
				topStatus.end = time.Now()
				collector.statuses[filename] = collector.statuses[filename][:statusLen-1]

				duration := topStatus.end.Sub(topStatus.start)
				ids := []string{}
				collector.insertTimeAggregation(ids, duration)
				for _, parentStatus := range collector.statuses[filename] {
					ids = append(ids, parentStatus.processName)
					collector.insertTimeAggregation(ids, duration)
				}
				ids = append(ids, topStatus.processName)
				collector.insertTimeAggregation(ids, duration)

				collector.debugOutputEvent("[Event] %s - End %s\n", filename, message.processName)

				if gracefulShutdown {
					stillProcessing := false
					for _, statuses := range collector.statuses {
						if len(statuses) > 0 {
							stillProcessing = true
						}
					}
					if !stillProcessing {
						collector.done <- true
						return
					}
				}
			case "done":
				collector.updater.Done()
				updaterDeinitialized = true

				gracefulShutdown = true
				stillProcessing := false
				for _, statuses := range collector.statuses {
					if len(statuses) > 0 {
						stillProcessing = true
					}
				}
				if !stillProcessing {
					collector.done <- true
					return
				}
			}
		}
	}()
}

func StartEventsCollector(updater StatusUpdater) *EventsCollector {
	eventStream := make(chan EventMessage)
	done := make(chan bool)
	collector := &EventsCollector{
		done:               done,
		eventStream:        eventStream,
		statuses:           map[string][]InputStatus{},
		timingsAggregation: map[string]time.Duration{},
		updater:            updater,
		outputFiles:        map[string]map[string]string{},
	}
	defer collector.runEventsCollectorDeamon()
	return collector
}

type InternalError interface {
	CliMessage() string
	Error() string
	ErrorName() string
}

type CollectedError struct {
	filename string
	err      InternalError
}

func (e CollectedError) ErrorName() string {
	return e.err.ErrorName()
}

func (e CollectedError) Filename() string {
	return e.filename
}

func (e CollectedError) Error() string {
	return e.err.Error()
}

func (e CollectedError) CliMessage() string {
	return e.err.CliMessage()
}

type TimingsAggreagation struct {
	Duration time.Duration
	Name     string
	Children []*TimingsAggreagation
}

type CollectedMetrics interface {
	GetAllErrors() []CollectedError
	GetTimingsAggregation() TimingsAggreagation
	Resolve() CollectedMetrics
	Inputs() []context.EventCollectorMessageInput
	GetOutputs() map[string]map[string]string
}

type CollectedMetricsImpl struct {
	errs               []CollectedError
	timingsAggregation map[string]time.Duration
	timingsLabels      []string
	inputs             []context.EventCollectorMessageInput
	outputs            map[string]map[string]string
}

func (c CollectedMetricsImpl) Inputs() []context.EventCollectorMessageInput {
	return c.inputs
}

func (c CollectedMetricsImpl) GetAllErrors() []CollectedError {
	return c.errs
}

func (c CollectedMetricsImpl) GetTimingsAggregation() TimingsAggreagation {
	ret := TimingsAggreagation{
		Duration: c.timingsAggregation[""],
		Name:     "",
		Children: []*TimingsAggreagation{},
	}
	for _, name := range c.timingsLabels {
		val := c.timingsAggregation[name]
		if len(name) == 0 {
			continue
		}
		tokens := strings.Split(name, "|")
		scope := &ret
		for _, t := range tokens {
			token := t
			foundScope := false
			for _, agg := range scope.Children {
				if agg.Name == token {
					foundScope = true
					scope = agg
					break
				}
			}
			if !foundScope {
				v := &TimingsAggreagation{
					Duration: 0,
					Name:     token,
					Children: []*TimingsAggreagation{},
				}
				scope.Children = append(scope.Children, v)
				scope = v
			}
		}
		scope.Duration = val
	}
	return ret
}

func (c CollectedMetricsImpl) Resolve() CollectedMetrics {
	return c
}

func (c CollectedMetricsImpl) GetOutputs() map[string]map[string]string {
	return c.outputs
}

type CollectedMetricsPromise interface {
	Resolve() CollectedMetrics
}

type CollectedErrorsPromiseChan <-chan CollectedMetrics

func (p CollectedErrorsPromiseChan) Resolve() CollectedMetrics {
	return <-p
}

func (ec *EventsCollector) collectSyncParsingError(program parser.LatteParsedProgramPromise, c *context.ParsingContext, out []CollectedError) ([]CollectedError, bool, parser.LatteParsedProgram) {
	result := program.Resolve()
	if result.ParsingError() != nil {
		return append(out, CollectedError{
			filename: result.Filename(),
			err:      result.ParsingError(),
		}), false, result
	} else {
		// Do nothing
		return out, true, result
	}
}

func (ec *EventsCollector) collectSyncParsingErrors(programs []parser.LatteParsedProgramPromise, c *context.ParsingContext, out []CollectedError) ([]CollectedError, []context.EventCollectorMessageInput) {
	inputs := []context.EventCollectorMessageInput{}
	for _, program := range programs {
		var result context.EventCollectorMessageInput
		out, _, result = ec.collectSyncParsingError(program, c, out)
		inputs = append(inputs, result)
	}
	return out, inputs
}

func (ec *EventsCollector) collectSyncTypecheckingError(program type_checker.LatteTypecheckedProgramPromise, c *context.ParsingContext, out []CollectedError) ([]CollectedError, bool, type_checker.LatteTypecheckedProgram) {
	result := program.Resolve()
	if result.TypeCheckingError != nil {
		return append(out, CollectedError{
			filename: result.Filename(),
			err:      result.TypeCheckingError,
		}), false, result
	} else {
		out, ok, _ := ec.collectSyncParsingError(result.Program, c, out)
		return out, ok, result
	}
}

func (ec *EventsCollector) collectSyncTypecheckingErrors(programs []type_checker.LatteTypecheckedProgramPromise, c *context.ParsingContext, out []CollectedError) ([]CollectedError, []context.EventCollectorMessageInput) {
	inputs := []context.EventCollectorMessageInput{}
	for _, program := range programs {
		var result context.EventCollectorMessageInput
		out, _, result = ec.collectSyncTypecheckingError(program, c, out)
		inputs = append(inputs, result)
	}
	return out, inputs
}

func (ec *EventsCollector) collectSyncFlowAnalysisError(program flow_analysis.LatteAnalyzedProgramPromise, c *context.ParsingContext, out []CollectedError) ([]CollectedError, bool, flow_analysis.LatteAnalyzedProgram) {
	result := program.Resolve()
	if result.FlowAnalysisError != nil {
		return append(out, CollectedError{
			filename: result.Filename(),
			err:      result.FlowAnalysisError,
		}), false, result
	} else {
		out, ok, _ := ec.collectSyncTypecheckingError(result.Program, c, out)
		return out, ok, result
	}
}

func (ec *EventsCollector) collectSyncCompilationError(program compiler.LatteCompiledProgramPromise, c *context.ParsingContext, out []CollectedError) ([]CollectedError, bool, compiler.LatteCompiledProgram) {
	result := program.Resolve()
	if result.CompilationError != nil {
		return append(out, CollectedError{
			filename: result.Filename(),
			err:      result.CompilationError,
		}), false, result
	} else {
		out, ok, _ := ec.collectSyncFlowAnalysisError(result.Program, c, out)
		return out, ok, result
	}
}

func (ec *EventsCollector) collectSyncCompiledCodeRunError(program runner.LatteRunnedProgramPromise, c *context.ParsingContext, out []CollectedError) ([]CollectedError, bool, runner.LatteRunnedProgram) {
	result := program.Resolve()
	if result.RunError != nil {
		return append(out, CollectedError{
			filename: result.Filename(),
			err:      result.RunError,
		}), false, result
	} else {
		out, ok, _ := ec.collectSyncCompilationError(result.Program, c, out)
		return out, ok, result
	}
}

func (ec *EventsCollector) collectSyncCompilationErrors(programs []compiler.LatteCompiledProgramPromise, c *context.ParsingContext, out []CollectedError) ([]CollectedError, []context.EventCollectorMessageInput) {
	inputs := []context.EventCollectorMessageInput{}
	for _, program := range programs {
		var result context.EventCollectorMessageInput
		out, _, result = ec.collectSyncCompilationError(program, c, out)
		inputs = append(inputs, result)
	}
	return out, inputs
}

func (ec *EventsCollector) collectSyncCompiledCodeRunErrors(programs []runner.LatteRunnedProgramPromise, c *context.ParsingContext, out []CollectedError) ([]CollectedError, []context.EventCollectorMessageInput) {
	inputs := []context.EventCollectorMessageInput{}
	for _, program := range programs {
		var result context.EventCollectorMessageInput
		out, _, result = ec.collectSyncCompiledCodeRunError(program, c, out)
		inputs = append(inputs, result)
	}
	return out, inputs
}

///

func (ec *EventsCollector) HandleCompiledCodeRunning(programs []runner.LatteRunnedProgramPromise, c *context.ParsingContext) CollectedMetricsPromise {
	ret := make(chan CollectedMetrics)
	go func() {
		defer close(ret)
		out := []CollectedError{}
		var inputs []context.EventCollectorMessageInput
		out, inputs = ec.collectSyncCompiledCodeRunErrors(programs, c, out)
		ec.eventStream <- EventMessage{
			eventType: "done",
		}
		<-ec.done
		ret <- CollectedMetricsImpl{
			errs:               out,
			timingsAggregation: ec.timingsAggregation,
			timingsLabels:      ec.timingsLabels,
			inputs:             inputs,
			outputs:            ec.outputFiles,
		}
	}()
	return CollectedErrorsPromiseChan(ret)
}

func (ec *EventsCollector) SummarizeCompiledCodeRunning(summarizer Summarizer, programs []runner.LatteRunnedProgramPromise, c *context.ParsingContext) (string, []runner.LatteRunnedProgram, bool) {
	progs := []runner.LatteRunnedProgram{}
	progsP := []runner.LatteRunnedProgramPromise{}
	for _, p := range programs {
		pResolved := p.Resolve()
		progs = append(progs, pResolved)
		progsP = append(progsP, pResolved)
	}
	message, ok := summarizer.Summarize(ec.HandleCompiledCodeRunning(progsP, c))
	return message, progs, ok
}

///

func (ec *EventsCollector) HandleCompilation(programs []compiler.LatteCompiledProgramPromise, c *context.ParsingContext) CollectedMetricsPromise {
	ret := make(chan CollectedMetrics)
	go func() {
		defer close(ret)
		out := []CollectedError{}
		var inputs []context.EventCollectorMessageInput
		out, inputs = ec.collectSyncCompilationErrors(programs, c, out)
		ec.eventStream <- EventMessage{
			eventType: "done",
		}
		<-ec.done
		ret <- CollectedMetricsImpl{
			errs:               out,
			timingsAggregation: ec.timingsAggregation,
			timingsLabels:      ec.timingsLabels,
			inputs:             inputs,
			outputs:            ec.outputFiles,
		}
	}()
	return CollectedErrorsPromiseChan(ret)
}

func (ec *EventsCollector) SummarizeCompilation(summarizer Summarizer, programs []compiler.LatteCompiledProgramPromise, c *context.ParsingContext) (string, bool) {
	return summarizer.Summarize(ec.HandleCompilation(programs, c))
}

func (ec *EventsCollector) HandleTypechecking(programs []type_checker.LatteTypecheckedProgramPromise, c *context.ParsingContext) CollectedMetricsPromise {
	ret := make(chan CollectedMetrics)
	go func() {
		defer close(ret)
		out := []CollectedError{}
		var inputs []context.EventCollectorMessageInput
		out, inputs = ec.collectSyncTypecheckingErrors(programs, c, out)
		ec.eventStream <- EventMessage{
			eventType: "done",
		}
		<-ec.done
		ret <- CollectedMetricsImpl{
			errs:               out,
			timingsAggregation: ec.timingsAggregation,
			timingsLabels:      ec.timingsLabels,
			inputs:             inputs,
			outputs:            ec.outputFiles,
		}
	}()
	return CollectedErrorsPromiseChan(ret)
}

func (ec *EventsCollector) SummarizeTypechecking(summarizer Summarizer, programs []type_checker.LatteTypecheckedProgramPromise, c *context.ParsingContext) (string, bool) {
	return summarizer.Summarize(ec.HandleTypechecking(programs, c))
}

func (ec *EventsCollector) HandleParsing(programs []parser.LatteParsedProgramPromise, c *context.ParsingContext) CollectedMetricsPromise {
	ret := make(chan CollectedMetrics)
	go func() {
		defer close(ret)
		out := []CollectedError{}
		var inputs []context.EventCollectorMessageInput
		out, inputs = ec.collectSyncParsingErrors(programs, c, out)
		ec.eventStream <- EventMessage{
			eventType: "done",
		}
		<-ec.done
		ret <- CollectedMetricsImpl{
			errs:               out,
			timingsAggregation: ec.timingsAggregation,
			timingsLabels:      ec.timingsLabels,
			inputs:             inputs,
			outputs:            ec.outputFiles,
		}
	}()
	return CollectedErrorsPromiseChan(ret)
}

func (ec *EventsCollector) SummarizeParsing(summarizer Summarizer, programs []parser.LatteParsedProgramPromise, c *context.ParsingContext) (string, bool) {
	return summarizer.Summarize(ec.HandleParsing(programs, c))
}
