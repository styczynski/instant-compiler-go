package events_collector

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

//func (c *ParsingContext) PrintProcessingInfo() string {
//	timingsDetails := []string{}
//
//	for name, stage := range c.Stages {
//		timingsDetails = append(timingsDetails, fmt.Sprintf("%s: Took %s",
//			formatStageTitleFg(name),
//			stage.End.Sub(*stage.Start),
//			))
//	}
//
//	return fmt.Sprintf("%s: Processed everything in %s:\n   | %s\n",
//		formatOkMessageBg(formatOkMessageFg("Done")),
//		c.End.Sub(*c.Start),
//		strings.Join(timingsDetails, "\n   | "))
//}

var formatStageTitleFg = color.New(color.FgMagenta).SprintFunc()

var formatOkMessageFg = color.New(color.FgHiWhite).SprintFunc()
var formatOkMessageBg = color.New(color.BgGreen).SprintFunc()

var formatOutputFilePathFg = color.New(color.FgMagenta).SprintFunc()

func generateIndent(nestLevel int) string {
	indent := strings.Repeat("  ", nestLevel)
	if nestLevel > 1 {
		indent = "  " + strings.Repeat("│ ", nestLevel-1)
	}
	return indent
}

func formatTimingAggregationRec(aggregation *TimingsAggreagation, nestLevel int) string {
	indent := generateIndent(nestLevel)
	lines := []string{}
	count := len(aggregation.Children)
	for i, agg := range aggregation.Children {
		prefix := "├─"
		if i == count-1 {
			prefix = "└─"
		}
		subtree := formatTimingAggregationRec(agg, nestLevel+1)
		if subtree != generateIndent(nestLevel+1) {
			subtree = "\n" + subtree
		} else {
			subtree = ""
		}
		lines = append(lines, fmt.Sprintf("%s %s - Took %s%s",
			prefix,
			formatStageTitleFg(agg.Name),
			agg.Duration,
			subtree))
	}
	return fmt.Sprintf("%s%s", indent, strings.Join(lines, fmt.Sprintf("\n%s", indent)))
}

func FormatTimingAggregation(aggregation TimingsAggreagation) string {
	t := &aggregation
	return formatTimingAggregationRec(t, 1)
}

func FormatOutputFilesList(outputFiles map[string]map[string]string) string {
	descriptionLines := []string{}
	for path, files := range outputFiles {
		filesTextExt := ""
		if len(files) > 1 {
			filesTextExt = "s"
		}
		descriptionLines = append(descriptionLines, fmt.Sprintf("   → Created %s file%s in %s:", formatOutputFilePathFg(len(files)), filesTextExt, formatOutputFilePathFg(path)))
		for name, description := range files {
			descriptionLines = append(descriptionLines, fmt.Sprintf("     ⚐ %s (%s)", name, description))
		}
	}
	return strings.Join(descriptionLines, "\n")
}
