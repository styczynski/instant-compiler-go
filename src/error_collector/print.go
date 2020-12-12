package error_collector

import (
	"fmt"
	"os"
)

func DumpCliMessageAndExit(errorsPromise CollectedErrorsPromise, limit int) {
	errors := errorsPromise.Resolve().GetAll()
	if len(errors) > 0 {
		for i, err := range errors {
			if i >= limit && limit != -1 {
				break
			}
			fmt.Print(err.CliMessage())
		}
		os.Exit(1)
	}
}
