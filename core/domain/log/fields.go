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
func (f Fields) Marshal(level Level, message string) ([]byte, error) {
	if _, ok := f["_stack"]; !ok {
		f["_stack"] = f.stack()
	}

	if _, ok := f["_time"]; !ok {
		f["_time"] = time.Now().Format("2006-01-02 15:04:05")
	}

	f["level"] = level.String()
	f["message"] = message

	var buf bytes.Buffer
	enc := logfmt.NewEncoder(&buf)

	var keys []string
	for k := range f {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		var err error

		switch f[k].(type) {
		case []string:
			// join strings with spaces
			err = enc.EncodeKeyval(k, strings.Join(f[k].([]string), " "))
		default:
			err = enc.EncodeKeyval(k, f[k])
		}

		if err != nil {
			return nil, err
		}
	}

	err := enc.EndRecord()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
