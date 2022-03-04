package log

import (
	"testing"

	panics "github.com/magiconair/properties/assert"
	"gotest.tools/v3/assert"
)

func TestFields_String(t *testing.T) {
	f := Fields{
		"foo":    "bar",
		"bar":    "baz",
		"_stack": "@", // simplifies testing
		"_time":  "@", // simplifies testing
	}

	assert.Equal(t, f.Marshal(DebugLevel, "message"), "_level=debug _message=message _stack=@ _time=@ bar=baz foo=bar\n")
}

func TestFields_String3(t *testing.T) {
	panics.Panic(t, func() {
		f := Fields{
			"foo": []string{"hello", "world"},
		}

		_ = f.Marshal(DebugLevel, "hello world")
	}, "unsupported value type")
}
