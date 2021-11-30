package events_collector

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/styczynski/latte-compiler/src/config"
)

func init() {
	config.RegisterEntityFactory(config.ENTITY_SUMMARIZER, CliSummaryShortStatusFactory{})
}

type CliSummaryShortStatusFactory struct{}

func (CliSummaryShortStatusFactory) CreateEntity(c config.EntityConfig) interface{} {
	return CreateCliSummaryShortStatus()
}

func (CliSummaryShortStatusFactory) Params(argSpec *config.EntityArgSpec) {
}

func (CliSummaryShortStatusFactory) EntityName() string {
	return "summary-short-cli"
}

type CliSummaryShortStatus struct {
}

var formatShortStatusOkFg = color.New(color.FgHiWhite).SprintFunc()
var formatShortStatusOkBg = color.New(color.BgGreen).SprintFunc()

var formatShortStatusErrFg = color.New(color.FgHiWhite).SprintFunc()
var formatShortStatusErrBg = color.New(color.BgRed).SprintFunc()

var formatShortStatusInputFilenameFg = color.New(color.FgBlue).SprintFunc()

func CreateCliSummaryShortStatus() CliSummaryShortStatus {
	return CliSummaryShortStatus{}
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

func centeredText(text string) string {
	lineWidth := 28
	spaceLeft := lineWidth - len(text)
	leftAlign := spaceLeft / 2
	rightAlign := spaceLeft - leftAlign
	if leftAlign < 0 || rightAlign < 0 {
		return text[:lineWidth]
	}
	return fmt.Sprintf(fmt.Sprintf("%%-%ds", lineWidth), fmt.Sprintf(fmt.Sprintf("%%%ds", rightAlign+len(text)), text))
}

func substr(input string, start int, length int) string {
	if start < 0 {
		start = 0
	}
	end := start + length
	if end < 0 {
		end = 0
	}
	if end > len(input) {
		end = len(input)
	}
	if start > end {
		start = end
	}
	return input[start:end]
}

func fillStringPostfix(input string, maxLength int) string {
	size := len(input)
	if size < maxLength {
		return input + strings.Repeat(" ", maxLength-size)
	}
	return input[:maxLength]
}

func cutFormatFilename(path string, maxLength int) string {
	i := strings.LastIndex(path, "/")

	if len(path) <= maxLength {
		return fillStringPostfix(path, maxLength)
	}

	tokenRight := substr(path, i, len(path)-i)
	if len(tokenRight) > maxLength-3 {
		ret := "..." + substr(path, len(path)-maxLength+3, maxLength-3)
		return fillStringPostfix(ret, maxLength)
	}

	tokenCenter := "..."
	tokenLeft := substr(path, 0, maxLength-3-(len(tokenRight)+len(tokenCenter)))

	ret := tokenLeft + tokenCenter + tokenRight
	return fillStringPostfix(ret, maxLength)
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

		errMessage := formatShortStatusOkFg(formatShortStatusOkBg(centeredText("OK")))
		if inputErr != nil {
			errMessage = formatShortStatusErrFg(formatShortStatusErrBg(centeredText(inputErr.ErrorName())))
		}
		lines = append(lines, fmt.Sprintf(" %3d: %50s - %s",
			i+1,
			formatShortStatusInputFilenameFg(cutFormatFilename(input.Filename(), 50)),
			errMessage))
	}

	return fmt.Sprintf("  %s\n", strings.Join(lines, "\n  ")), true
}
