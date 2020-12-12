package events_collector

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type CliSummaryShortStatus struct {
}

var formatShortStatusOkFg = color.New(color.FgHiWhite).SprintFunc()
var formatShortStatusOkBg = color.New(color.BgGreen).SprintFunc()

var formatShortStatusErrFg = color.New(color.FgHiWhite).SprintFunc()
var formatShortStatusErrBg = color.New(color.BgRed).SprintFunc()

var formatShortStatusInputFilenameFg = color.New(color.FgBlue).SprintFunc()

func CreateCliSummaryShortStatus() CliSummaryShortStatus {
	return CliSummaryShortStatus{
	}
}

func (s CliSummaryShortStatus) FormatCliSummaryShortStatus(metricsPromise CollectedMetricsPromise) string {
	metrics := metricsPromise.Resolve()
	timings := metrics.GetTimingsAggregation()

	return fmt.Sprintf("%s: Processed everything in %s (%d inputs):\n%s\n",
			formatOkMessageBg(formatOkMessageFg("Done")),
			timings.Duration,
			len(metrics.Inputs()),
			FormatTimingAggregation(timings))
}

func (s CliSummaryShortStatus) Summarize(metricsPromise CollectedMetricsPromise) (string, bool) {
	metrics := metricsPromise.Resolve()
	errors := metrics.GetAllErrors()
	lines := []string{}
	for i, input := range metrics.Inputs() {
		var inputErr *CollectedError = nil
		for _, err := range errors {
			if err.Filename() == input.Filename() {
				inputErr = &err
				break
			}
		}
		errMessage := formatShortStatusOkFg(formatShortStatusOkBg(fmt.Sprintf("%-20s", fmt.Sprintf("%10s", "OK"))))
		if inputErr != nil {
			errMessage = formatShortStatusErrFg(formatShortStatusErrBg(fmt.Sprintf("%-20s", fmt.Sprintf("%10s", inputErr.ErrorName()))))
		}
		lines = append(lines, fmt.Sprintf(" %3d: %10s - %s",
			i+1,
			formatShortStatusInputFilenameFg(input.Filename()),
			errMessage))
	}

	return fmt.Sprintf("  %s\n", strings.Join(lines, "\n  ")), true
}
