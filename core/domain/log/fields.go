package log

import (
	"bytes"
	"fmt"
	"github.com/go-logfmt/logfmt"
	"path"
	"runtime"
	"sort"
	"strings"
	"time"
)

// Fields is a map of fields, which may be marshaled into a logfmt-compatible string.
type Fields map[string]any

func (f Fields) stack() string {
	var stack string

	for i := 0; i < 3; i++ {
		_, file, line, ok := runtime.Caller(i + 1)
		if !ok {
			break
		}

		stack += fmt.Sprintf("%s:%d;", path.Base(file), line)
	}

	return strings.TrimRight(stack, ";")
}

// Marshal returns a logfmt-compatible string representation of the fields.
func (f Fields) Marshal(level level, message string) string {
	if _, ok := f["_level"]; !ok {
		f["_level"] = level.String()
	}

	if _, ok := f["_message"]; !ok {
		f["_message"] = message
	}

	if _, ok := f["_time"]; !ok {
		f["_time"] = time.Now().Format("2006-01-02 15:04:05")
	}

	if _, ok := f["_stack"]; !ok {
		f["_stack"] = f.stack()
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
