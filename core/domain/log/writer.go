package log

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

// writer is the log writer interface.
type writer interface {
	Write(level level, message string, fields Fields) error
}

// fileWriter is the file log writer.
type fileWriter struct {
	file *os.File
}

// Write writes the log message to the file.
func (f *fileWriter) Write(level level, message string, fields Fields) error {
	fields["level"] = level.String()
	fields["message"] = message
	fields["time"] = time.Now().Format("2006-01-02 15:04:05")

	var stack string

	for i := 0; i < 4; i++ {
		_, file, line, ok := runtime.Caller(i + 1)
		if !ok {
			break
		}

		stack += fmt.Sprintf("%s:%d;", path.Base(file), line)
	}

	fields["stack"] = strings.TrimRight(stack, ";")

	_, err := f.file.Write([]byte(fields.String() + "\n"))
	return err
}

// compositeWriter logs a given message to multiple writers.
type compositeWriter struct {
	writers []writer
}

// Write writes the log message to all writers.
func (c *compositeWriter) Write(level level, message string, fields Fields) error {
	for _, writer := range c.writers {
		err := writer.Write(level, message, fields)
		if err != nil {
			return err
		}
	}

	return nil
}

// TestEvent contains the log values for testing.
type TestEvent struct {
	Level   level
	Message string
	Fields  Fields
}

// MemoryWriter is a writer for testing only.
// It stores logs in memory.
type MemoryWriter struct {
	Events []TestEvent
}

// Write writes the log message to the memory.
func (m *MemoryWriter) Write(level level, message string, fields Fields) error {
	m.Events = append(m.Events, TestEvent{level, message, fields})
	return nil
}

// Last is a convenience method for getting the last event.
func (m *MemoryWriter) Last() TestEvent {
	return m.Events[len(m.Events)-1]
}

// LastN is a convenience method for getting the last N events.
func (m *MemoryWriter) LastN(n int) []TestEvent {
	if len(m.Events) < n {
		return m.Events
	}

	return m.Events[len(m.Events)-n:]
}

func (m *MemoryWriter) Len() int {
	return len(m.Events)
}
