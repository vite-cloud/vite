package log

import (
	"fmt"
	"gotest.tools/v3/assert"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestCompositeWriter_Write(t *testing.T) {
	writers := []Writer{
		&MemoryWriter{},
		&MemoryWriter{},
	}

	w := &CompositeWriter{writers}
	err := w.Write(DebugLevel, "test", Fields{
		"foo": "bar",
	})
	assert.NilError(t, err)

	for _, wr := range writers {
		wr := wr.(*MemoryWriter)
		assert.Equal(t, 1, wr.Len())
		assert.Equal(t, DebugLevel, wr.Last().Level)
		assert.Equal(t, "test", wr.Last().Message)
		assert.Equal(t, "bar", wr.Last().Fields["foo"])
	}
}

type failingWriter struct{}

func (f failingWriter) Write(level Level, message string, fields Fields) error {
	return fmt.Errorf("failed")
}

func TestCompositeWriter_Write2(t *testing.T) {
	w := &CompositeWriter{[]Writer{
		&failingWriter{},
	}}

	err := w.Write(DebugLevel, "test", Fields{
		"foo": "bar",
	})

	assert.ErrorContains(t, err, "failed")
}

func TestMemoryWriter_Len(t *testing.T) {
	w := &MemoryWriter{}
	assert.Equal(t, 0, w.Len())

	for i := 0; i < 3; i++ {
		err := w.Write(DebugLevel, "test", Fields{
			"foo": "bar",
		})
		assert.NilError(t, err)
		assert.Assert(t, w.Len() == i+1)
	}
}

func TestMemoryWriter_Last(t *testing.T) {
	w := &MemoryWriter{}
	err := w.Write(DebugLevel, "test", Fields{
		"foo": "bar",
	})
	assert.NilError(t, err)
	assert.Equal(t, DebugLevel, w.Last().Level)
	assert.Equal(t, "test", w.Last().Message)
	assert.Equal(t, "bar", w.Last().Fields["foo"])
}

func TestMemoryWriter_LastN(t *testing.T) {
	w := &MemoryWriter{}
	for i := 0; i < 3; i++ {
		err := w.Write(DebugLevel, "test", Fields{
			"foo": "bar",
		})
		assert.NilError(t, err)
	}

	assert.Equal(t, 3, w.Len())
}

func TestMemoryWriter_Write(t *testing.T) {
	w := &MemoryWriter{}
	err := w.Write(DebugLevel, "test", Fields{
		"foo": "bar",
	})
	assert.NilError(t, err)
	assert.Equal(t, 1, w.Len())
	assert.Equal(t, DebugLevel, w.Last().Level)
	assert.Equal(t, "test", w.Last().Message)
	assert.Equal(t, "bar", w.Last().Fields["foo"])
}

func TestFileWriter_Write(t *testing.T) {
	path := "/tmp/tmp-filew-" + strconv.Itoa(int(time.Now().UnixMilli()))
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0600)
	assert.NilError(t, err)

	defer file.Close()

	w := FileWriter{
		File: file,
	}

	err = w.Write(DebugLevel, "test", Fields{
		"foo":    "bar",
		"_time":  "@", // simplifies testing
		"_stack": "@", // simplifies testing
	})
	assert.NilError(t, err)

	contents, err := os.ReadFile(path)
	assert.NilError(t, err)

	assert.Equal(t, "_stack=@ _time=@ foo=bar level=debug message=test\n", string(contents))
}

func TestFileWriter_Write2(t *testing.T) {
	w := FileWriter{}

	err := w.Write(DebugLevel, "test", Fields{
		"invalid_value": []int{1, 2, 3},
	})
	assert.Error(t, err, "unsupported value type")
}
