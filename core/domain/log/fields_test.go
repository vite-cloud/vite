package log

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestFields_Marshal(t *testing.T) {
	f := Fields{
		"foo":    "bar",
		"bar":    "baz",
		"_stack": "@", // simplifies testing
		"_time":  "@", // simplifies testing
	}

	l, err := f.Marshal(DebugLevel, "message")
	assert.NilError(t, err)
	assert.Equal(t, string(l), "_stack=@ _time=@ bar=baz foo=bar level=debug message=message\n")
}

func TestFields_Marshal3(t *testing.T) {
	f := Fields{
		"foo":    []string{"bar", "baz"},
		"_stack": "@", // simplifies testing
		"_time":  "@", // simplifies testing
	}

	l, err := f.Marshal(DebugLevel, "message")
	assert.NilError(t, err)
	assert.Equal(t, string(l), "_stack=@ _time=@ foo=\"bar baz\" level=debug message=message\n")
}

func TestFields_Marshal2(t *testing.T) {
	f := Fields{
		"foo": [][]string{{"hello", "world"}},
	}

	_, err := f.Marshal(DebugLevel, "hello world")
	assert.Error(t, err, "unsupported value type")
}
