package log

import (
	"io"
)

// Writer is the log Writer interface.
type Writer interface {
	Write(level Level, message string, fields Fields) error
}

// FileWriter is the file log Writer.
type FileWriter struct {
	File io.Writer
}

// Write writes the log message to the file.
func (f *FileWriter) Write(level Level, message string, fields Fields) error {
	marshalled, err := fields.Marshal(level, message)
	if err != nil {
		return err
	}

	_, err = f.File.Write(marshalled)
	return err
}

// CompositeWriter logs a given message to multiple writers.
type CompositeWriter struct {
	Writers []Writer
}

// Write writes the log message to all writers.
func (c *CompositeWriter) Write(level Level, message string, fields Fields) error {
	for _, writer := range c.Writers {
		err := writer.Write(level, message, fields)
		if err != nil {
			return err
		}
	}

	return nil
}

// TestEvent contains the log values for testing.
type TestEvent struct {
	Level   Level
	Message string
	Fields  Fields
}

// MemoryWriter is a writer for testing only.
// It stores logs in memory.
type MemoryWriter struct {
	Events []TestEvent
}

// Write writes the log message to the memory.
func (m *MemoryWriter) Write(level Level, message string, fields Fields) error {
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
