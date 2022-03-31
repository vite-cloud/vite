package proxy

import (
	"github.com/vite-cloud/vite/core/domain/log"
	"os"
)

const (
	LogFile = "proxy.log"
)

// logger contains an instance of the proxy logger.
var logger log.Writer

// SetLogger sets the proxy logger to a given writer.
func SetLogger(w log.Writer) {
	logger = w
}

// GetLogger returns the proxy logger.
func GetLogger() log.Writer {
	return logger
}

// defaultLogger creates a default logger.
func defaultLogger() (log.Writer, error) {
	dir, err := log.Store.Dir()
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(dir+"/"+LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}

	return &log.FileWriter{File: file}, nil
}

// Log logs an internal event to the proxy logger
func Log(level log.Level, message string, fields log.Fields) {
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
