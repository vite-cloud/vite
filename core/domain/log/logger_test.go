package log

import (
	"fmt"
	"os"
	"testing"

	panics "github.com/magiconair/properties/assert"
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
		"_stack": "@", // simplifies testing
		"_time":  "@", // simplifies testing
	})

	logDir, err := Store.Dir()
	assert.NilError(t, err)

	logFile := logDir + "/internal.log"

	contents, err := os.ReadFile(logFile)
	assert.NilError(t, err)

	assert.Equal(t, string(contents), "_stack=@ _time=@ foo=\"bar baz\" level=debug message=\"hello world\"\n")
}

func TestGetLogger(t *testing.T) {
	l := &MemoryWriter{}
	SetLogger(l)

	assert.Equal(t, GetLogger(), l)
}

func TestLog3(t *testing.T) {
	SetLogger(nil)
	datadir.SetHomeDir("/nop")

	panics.Panic(t, func() {
		Log(DebugLevel, "hello world", Fields{})
	}, "mkdir /nop: permission denied")
}

func TestLog4(t *testing.T) {
	home, err := os.MkdirTemp("", "vite-datadir")
	assert.NilError(t, err)

	datadir.SetHomeDir(home)

	dir, err := Store.Dir()
	assert.NilError(t, err)

	err = os.Mkdir(dir+"/internal.log", 0600)
	assert.NilError(t, err)

	panics.Panic(t, func() {
		Log(DebugLevel, "hello world", Fields{})
	}, fmt.Sprintf("open %s: is a directory", dir+"/internal.log"))
}

func TestLog5(t *testing.T) {
	SetLogger(&failingWriter{})

	panics.Panic(t, func() {
		Log(DebugLevel, "hello world", Fields{})
	}, "failed")
}
