package log

import (
	"bytes"
	"github.com/go-logfmt/logfmt"
	"sort"
)

// Fields is a map of fields, which may be marshalled into a logfmt-compatible string.
type Fields map[string]any

// String returns a logfmt-compatible string representation of the fields.
func (f Fields) String() string {
	if len(f) == 0 {
		return ""
	}

	var buf bytes.Buffer
	enc := logfmt.NewEncoder(&buf)

	var keys []string
	for k := range f {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		err := enc.EncodeKeyval(k, f[k])
		if err != nil {
			panic(err)
		}
	}

	err := enc.EndRecord()
	if err != nil {
		panic(err)
	}

	return buf.String()
}
