package log

import (
	"os"
	"testing"

	"github.com/vite-cloud/vite/core/domain/datadir"
	"gotest.tools/v3/assert"
)

func TestLog(t *testing.T) {
	l := &MemoryWriter{}
	SetLogger(l)

	Log(DebugLevel, "hello world", Fields{
		"a": "b",
	})

	assert.Equal(t, l.Len(), 1)
	assert.Equal(t, l.Last().Level, DebugLevel)
	assert.Equal(t, l.Last().Message, "hello world")
	assert.Equal(t, l.Last().Fields["a"], "b")

}

func TestLog2(t *testing.T) {
	dir, err := os.MkdirTemp("", "logger-test")
	assert.NilError(t, err)

	datadir.SetHomeDir(dir)

	SetLogger(nil)

	Log(DebugLevel, "hello world", Fields{
		"foo":    "bar baz",
		"_stack": "@",
		"_time":  "@",
	})

	logDir, err := Store.Dir()
	assert.NilError(t, err)

	logFile := logDir + "/internal.log"

	contents, err := os.ReadFile(logFile)
	assert.NilError(t, err)

	assert.Equal(t, string(contents), "_level=debug _message=\"hello world\" _stack=@ _time=@ foo=\"bar baz\"\n")
}

func TestGetLogger(t *testing.T) {
	l := &MemoryWriter{}
	SetLogger(l)

	assert.Equal(t, GetLogger(), l)
}
