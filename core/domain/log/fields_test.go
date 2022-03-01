package log

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestFields_String(t *testing.T) {
	f := Fields{
		"foo": "bar",
		"bar": "baz",
	}
	if f.String() != "bar=baz foo=bar\n" {
		t.Errorf("Fields.String() returned %s", f.String())
	}
}

func TestFields_String2(t *testing.T) {
	f := Fields{}
	if f.String() != "" {
		t.Errorf("Fields.String() returned %s", f.String())
	}
}

func TestFields_String3(t *testing.T) {
	assert.Panic(t, func() {
		f := Fields{
			"foo": []string{"hello", "world"},
		}

		_ = f.String()
	}, "unsupported value type")
}
