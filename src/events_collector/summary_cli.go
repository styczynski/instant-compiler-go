package events_collector

import (
	"fmt"

	"github.com/styczynski/latte-compiler/src/config"
)

func init() {
	config.RegisterEntityFactory(config.ENTITY_SUMMARIZER, CliSummaryFactory{})
}

type CliSummaryFactory struct{}

func (CliSummaryFactory) CreateEntity(c config.EntityConfig) interface{} {
	return CreateCliSummary(c.Int("error-limit"))
}

func (CliSummaryFactory) Params(argSpec *config.EntityArgSpec) {
	argSpec.AddInt("error-limit", -1, "Maximum error limit")
}

func (CliSummaryFactory) EntityName() string {
	return "summary-cli"
}

type CliSummary struct {
	errorsLimit int
}

func CreateCliSummary(errorsLimit int) CliSummary {
	return CliSummary{
		errorsLimit: errorsLimit,
	}
}

func (s CliSummary) FormatCliSummary(metricsPromise CollectedMetricsPromise) string {
	metrics := metricsPromise.Resolve()
	timings := metrics.GetTimingsAggregation()

	return fmt.Sprintf("%s: Processed everything in %s (%d inputs):\n%s\n\n%s",
		formatOkMessageBg(formatOkMessageFg("Done")),
		timings.Duration,
		len(metrics.Inputs()),
		FormatTimingAggregation(timings),
		FormatOutputFilesList(metrics.GetOutputs()))
}

func (s CliSummary) Summarize(metricsPromise CollectedMetricsPromise) (string, bool) {
	metrics := metricsPromise.Resolve()
	errors := metrics.GetAllErrors()
	if len(errors) > 0 {
		for i, err := range errors {
			if i >= s.errorsLimit && s.errorsLimit != -1 {
				break
			}
			return err.CliMessage(), false
		}
	}
	return s.FormatCliSummary(metrics), true
}
