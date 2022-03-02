package log

import (
	"fmt"
	"github.com/vite-cloud/vite/core/domain/datadir"
	"os"
)

var (
	logger writer
)

// Store is the unique name of the logger store
const Store = datadir.Store("logs")

// newLogger creates a new writer to log internal events.
func newLogger() (writer, error) {
	dir, err := Store.Dir()
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(dir+"/internal.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}

	return &fileWriter{file}, nil
}

func UseTestLogger() *MemoryWriter {
	testWriter := &MemoryWriter{}

	logger = testWriter

	return testWriter
}

// Log logs internal events to the configured log target.
func Log(level level, format string, fields Fields) {
	if logger == nil {
		w, err := newLogger()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		logger = w
	}

	err := logger.Write(level, format, fields)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
