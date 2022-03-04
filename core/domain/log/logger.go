package log

import (
	"os"

	"github.com/vite-cloud/vite/core/domain/datadir"
)

// logger contains an instance of the global logger.
var logger writer

// Store is the unique name of the logger store
const Store = datadir.Store("logs")

// SetLogger sets the global logger to a given writer.
func SetLogger(w writer) {
	logger = w
}

// GetLogger returns the global logger.
func GetLogger() writer {
	return logger
}

// defaultLogger creates a default logger.
func defaultLogger() (writer, error) {
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

// Log logs an internal event to the global logger
func Log(level level, message string, fields Fields) {
	if logger == nil {
		w, err := defaultLogger()
		if err != nil {
			panic(err)
		}

		SetLogger(w)
	}

	err := logger.Write(level, message, fields)
	if err != nil {
		panic(err)
	}
}
