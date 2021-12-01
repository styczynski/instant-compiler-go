package logs

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/styczynski/latte-compiler/src/logs/formatter"
	"github.com/styczynski/latte-compiler/src/parser/context"
)

var enableLogging = true

var tabcount uint32

var replacement = "\n"

type LogContext interface {
	LogContext(c *context.ParsingContext) map[string]interface{}
}

func InitializeLogger(level string, output string) error {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&formatter.Formatter{
		ShowFullLevel:   true,
		HideKeys:        true,
		NoColors:        false,
		NoFieldsColors:  false,
		TimestampFormat: "15:04:05.000",
		FieldsOrder:     []string{"component", "category"},
	})

	if output == "stdout" {
		log.SetOutput(os.Stdout)
	} else if output == "stderr" {
		log.SetOutput(os.Stderr)
	} else if output == "none" {
		enableLogging = false
	} else {
		file, err := os.OpenFile("latc.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		log.SetOutput(file)
	}

	log.SetLevel(log.WarnLevel)
	if level == "debug" {
		log.SetLevel(log.DebugLevel)
	} else if level == "error" {
		log.SetLevel(log.ErrorLevel)
	} else if level == "warning" {
		log.SetLevel(log.WarnLevel)
	} else if level == "info" {
		log.SetLevel(log.InfoLevel)
	} else {
		return fmt.Errorf("Invalid logging level was specified: %s. Expected: debug, error, warning or info.", level)
	}

	return nil
}

var logsGlobalParsingContext *context.ParsingContext

func UseParsingContext(c *context.ParsingContext) {
	logsGlobalParsingContext = c
}

// func Debug(context LogContext, format string, others ...interface{}) {

func logf(context LogContext, format string, others ...interface{}) (*log.Entry, string) {
	s := fmt.Sprintf(format, others...)
	s = strings.Replace(s, "\n", replacement, -1)

	componentTokens := strings.Split(reflect.TypeOf(context).String(), ".")

	l := log.WithField("component", componentTokens[len(componentTokens)-1])
	for k, v := range context.LogContext(logsGlobalParsingContext) {
		l = l.WithField(k, v)
	}
	contextLogger := l.WithFields(log.Fields{})
	return contextLogger, s
}

func Debug(context LogContext, format string, others ...interface{}) {
	if enableLogging {
		logger, text := logf(context, format, others...)
		logger.Info(text)
	}
}

func Info(context LogContext, format string, others ...interface{}) {
	if enableLogging {
		logger, text := logf(context, format, others...)
		logger.Info(text)
	}
}

func Warning(context LogContext, format string, others ...interface{}) {
	if enableLogging {
		logger, text := logf(context, format, others...)
		logger.Warning(text)
	}
}

func Error(context LogContext, format string, others ...interface{}) {
	if enableLogging {
		logger, text := logf(context, format, others...)
		logger.Error(text)
	}
}
