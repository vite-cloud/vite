package log

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestCompositeWriter_Write(t *testing.T) {
	writers := []writer{
		&MemoryWriter{},
		&MemoryWriter{},
	}

	w := &compositeWriter{writers}
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
