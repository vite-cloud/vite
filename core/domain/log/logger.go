package log

import (
	"fmt"
	"github.com/vite-cloud/vite/core/domain/datadir"
	"os"
	"sync"
)

var (
	logger       writer
	createLogger = new(sync.Once)
)

// newLogger creates a new writer to log internal events.
func newLogger() (writer, error) {
	path, err := datadir.Path("logs", "internal.log")
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return &fileWriter{file}, nil
}

// Log logs internal events to the configured log target.
func Log(level level, format string, fields Fields) {
	createLogger.Do(func() {
		w, err := newLogger()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		logger = w
	})

	err := logger.Write(level, format, fields)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
