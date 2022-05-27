package log

import (
	"github.com/vite-cloud/go-zoup"
	"github.com/vite-cloud/vite/core/domain/datadir"
	"os"
)

// logger contains an instance of the global logger.
var logger zoup.Writer

const (
	// Store is the unique name of the logger store
	Store = datadir.Store("logs")
	// LogFile is the name of the log file
	LogFile = "internal.log"
)

// SetLogger sets the global logger to a given writer.
func SetLogger(w zoup.Writer) {
	logger = w
}

// GetLogger returns the global logger.
func GetLogger() zoup.Writer {
	return logger
}

// defaultLogger creates a default logger.
func defaultLogger() (zoup.Writer, error) {
	dir, err := Store.Dir()
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(dir+"/"+LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}

	return &zoup.FileWriter{file}, nil
}

// Log logs an internal event to the global logger
func Log(level zoup.Level, message string, fields zoup.Fields) {
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
