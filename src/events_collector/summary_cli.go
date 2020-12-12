package events_collector

import (
	"fmt"
)

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

	return fmt.Sprintf("%s: Processed everything in %s (%d inputs):\n%s\n",
			formatOkMessageBg(formatOkMessageFg("Done")),
			timings.Duration,
			len(metrics.Inputs()),
			FormatTimingAggregation(timings))
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
