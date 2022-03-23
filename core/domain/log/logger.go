package log

import (
	"errors"
	"github.com/hpcloud/tail"
	"github.com/vite-cloud/vite/core/domain/datadir"
	"io"
	"os"
)

// logger contains an instance of the global logger.
var logger writer

const (
	// Store is the unique name of the logger store
	Store = datadir.Store("logs")
	// LogFile is the name of the log file
	LogFile = "internal.log"
)

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

	file, err := os.OpenFile(dir+"/"+LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
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

func Tail(follow bool, n int) (*tail.Tail, error) {
	dir, err := Store.Dir()
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(dir+"/"+LogFile, os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	ret, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, err
	}

	offset := ret

	for {
		if n == 0 || offset == 0 {
			break
		}

		buf := make([]byte, 1)
		_, err = file.ReadAt(buf, offset-1)
		if err != nil && errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, err
		}

		// Skip last trailing newline
		if offset == ret && buf[0] == '\n' {
			offset--
			continue
		}

		if buf[0] == '\n' {
			n--

			if n == 0 {
				break
			}

			offset--
			continue
		}

		offset--
	}

	return tail.TailFile(dir+"/"+LogFile, tail.Config{
		Logger: tail.DiscardingLogger,
		Follow: follow,
		ReOpen: follow, // not a mistake, follow needs to be enabled to re-open the file
		Location: &tail.SeekInfo{
			Offset: offset,
			Whence: io.SeekStart,
		},
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
