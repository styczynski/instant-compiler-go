package errors

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/styczynski/latte-compiler/src/parser/context"
)

func GeneralRecovery(ctx *context.ParsingContext, stageName string, filename string, handler func (message string, textMessage string), finalHandler func()) {
	if r := recover(); r != nil {
		traceLines := strings.Split(string(debug.Stack()), "\n")
		newTraceLines := []string{}
		foundPanic := false
		for i, _ := range traceLines {
			if strings.Contains(traceLines[i], "panic.go") {
				foundPanic = true
				i++
				continue
			}
			if foundPanic {
				if strings.Contains(traceLines[i], ".go") {
					line := strings.Split(traceLines[i], " ")
					newTraceLines = append(newTraceLines, line[0])
				}
			}
		}
		traceLines = newTraceLines
		if len(traceLines) > 10 {
			traceLines = traceLines[:10]
			traceLines = append(traceLines, "... and more ...")
		}
		traceMessage := strings.Join(traceLines, "\n            ")
		baseMessage := r.(error).Error()
		fullDescription := fmt.Sprintf("%s\n            %s", baseMessage, traceMessage)
		message, textMessage := ctx.FormatParsingError(
			fmt.Sprintf("PANIC (%s)", stageName),
			fullDescription,
			0,
			0,
			filename,
			"",
			fullDescription)
		handler(message, textMessage)
	}
	finalHandler()
}
