package log

import (
	"os"
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
	marshalled, err := fields.Marshal(level, message)
	if err != nil {
		return err
	}

	_, err = f.file.Write(marshalled)
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

// Len returns the number of events in the memory.
func (m *MemoryWriter) Len() int {
	return len(m.Events)
}
