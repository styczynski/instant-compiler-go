package context

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var formatErrorBg = color.New(color.BgRed).SprintFunc()
var formatErrorFg = color.New(color.FgHiWhite).SprintFunc()

var formatErrorMessageFg = color.New(color.FgRed).SprintFunc()
var formatErrorMetaInfoFg = color.New(color.FgHiBlue).SprintFunc()

func IndentCodeLines(message string, errorLine int, lineStart int) string {
	lines := strings.Split(message, "\n")
	newLines := []string{}
	curLineNo := lineStart
	startTrimMode := true
	for _, line := range lines {
		add := true
		if len(line) > 0 {
			startTrimMode = false
		} else if startTrimMode {
			add = false
		}
		if add {
			lineMarker := ""
			if curLineNo == errorLine {
				lineMarker = "> "
			}
			newLines = append(newLines, fmt.Sprintf("%6s | %s", fmt.Sprintf("%s%d", lineMarker, curLineNo), line))
		}
		curLineNo++
	}
	return strings.Join(newLines, "\n")
}

func (c *ParsingContext) FormatParsingError(errorType string, message string, line int, col int, filename string, recommendedBracket string, errorMessage string) (string, string) {

	locationMessage := fmt.Sprintf("%s %s %s %d %s %d",
		formatErrorMetaInfoFg("Error in file"),
		filename,
		formatErrorMetaInfoFg("in line"),
		line,
		formatErrorMetaInfoFg("column"),
		col)

	codeContext, lineStart, _ := c.GetFileContext(nil, line, col)
	formattedCode, err := c.Printer.FormatRaw(codeContext, false)
	if err != nil {
		// Ignore error
		formattedCode = codeContext
	}

	textMessage := fmt.Sprintf("%s\n%s\n %s: %s\n",
		locationMessage,
		IndentCodeLines(formattedCode, line, lineStart),
		formatErrorFg(formatErrorBg(fmt.Sprintf(" %s ", errorType))),
		formatErrorMessageFg(errorMessage))

	return message, textMessage
}
