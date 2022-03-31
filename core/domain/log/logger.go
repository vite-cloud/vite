package log

import (
	"github.com/vite-cloud/vite/core/domain/datadir"
	"os"
)

// logger contains an instance of the global logger.
var logger Writer

const (
	// Store is the unique name of the logger store
	Store = datadir.Store("logs")
	// LogFile is the name of the log file
	LogFile = "internal.log"
)

// SetLogger sets the global logger to a given writer.
func SetLogger(w Writer) {
	logger = w
}

// GetLogger returns the global logger.
func GetLogger() Writer {
	return logger
}

// defaultLogger creates a default logger.
func defaultLogger() (Writer, error) {
	dir, err := Store.Dir()
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(dir+"/"+LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}

	return &FileWriter{file}, nil
}

// Log logs an internal event to the global logger
func Log(level Level, message string, fields Fields) {
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
