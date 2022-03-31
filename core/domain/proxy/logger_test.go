package proxy

import (
	"fmt"
	"github.com/vite-cloud/vite/core/domain/log"
	"os"
	"testing"

	panics "github.com/magiconair/properties/assert"
	"github.com/vite-cloud/vite/core/domain/datadir"
	"gotest.tools/v3/assert"
)

type failingWriter struct{}

func (f failingWriter) Write(level log.Level, message string, fields log.Fields) error {
	return fmt.Errorf("failed")
}

func TestLog(t *testing.T) {
	l := &log.MemoryWriter{}
	SetLogger(l)

	Log(log.DebugLevel, "hello world", log.Fields{
		"a": "b",
	})

	assert.Equal(t, l.Len(), 1)
	assert.Equal(t, l.Last().Level, log.DebugLevel)
	assert.Equal(t, l.Last().Message, "hello world")
	assert.Equal(t, l.Last().Fields["a"], "b")

}

func TestLog2(t *testing.T) {
	datadir.UseTestHome(t)

	SetLogger(nil)

	Log(log.DebugLevel, "hello world", log.Fields{
		"foo":    "bar baz",
		"_stack": "@", // simplifies testing
		"_time":  "@", // simplifies testing
	})

	logDir, err := log.Store.Dir()
	assert.NilError(t, err)

	logFile := logDir + "/" + LogFile

	contents, err := os.ReadFile(logFile)
	assert.NilError(t, err)

	assert.Equal(t, string(contents), "_stack=@ _time=@ foo=\"bar baz\" level=debug message=\"hello world\"\n")
}

func TestGetLogger(t *testing.T) {
	l := &log.MemoryWriter{}
	SetLogger(l)

	assert.Equal(t, GetLogger(), l)
}

func TestLog3(t *testing.T) {
	SetLogger(nil)
	datadir.SetHomeDir("/nop")

	panics.Panic(t, func() {
		Log(log.DebugLevel, "hello world", log.Fields{})
	}, "mkdir /nop: permission denied")
}

func TestLog4(t *testing.T) {
	datadir.UseTestHome(t)

	dir, err := log.Store.Dir()
	assert.NilError(t, err)

	err = os.Mkdir(dir+"/"+LogFile, 0600)
	assert.NilError(t, err)

	panics.Panic(t, func() {
		Log(log.DebugLevel, "hello world", log.Fields{})
	}, fmt.Sprintf("open %s: is a directory", dir+"/"+LogFile))
}

func TestLog5(t *testing.T) {
	SetLogger(&failingWriter{})

	panics.Panic(t, func() {
		Log(log.DebugLevel, "hello world", log.Fields{})
	}, "failed")
}
