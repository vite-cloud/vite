package log

import (
	"fmt"
	"github.com/vite-cloud/go-zoup"
	"os"
	"testing"

	panics "github.com/magiconair/properties/assert"
	"github.com/vite-cloud/vite/core/domain/datadir"
	"gotest.tools/v3/assert"
)

func TestLog(t *testing.T) {
	l := &zoup.MemoryWriter{}
	SetLogger(l)

	Log(zoup.DebugLevel, "hello world", zoup.Fields{
		"a": "b",
	})

	assert.Equal(t, l.Len(), 1)
	assert.Equal(t, l.Last().Level, zoup.DebugLevel)
	assert.Equal(t, l.Last().Message, "hello world")
	assert.Equal(t, l.Last().Fields["a"], "b")

}

func TestLog2(t *testing.T) {
	datadir.UseTestHome(t)

	SetLogger(nil)

	Log(zoup.DebugLevel, "hello world", zoup.Fields{
		"foo":    "bar baz",
		"_stack": "@", // simplifies testing
		"_time":  "@", // simplifies testing
	})

	logDir, err := Store.Dir()
	assert.NilError(t, err)

	logFile := logDir + "/" + LogFile

	contents, err := os.ReadFile(logFile)
	assert.NilError(t, err)

	assert.Equal(t, string(contents), "_stack=@ _time=@ foo=\"bar baz\" level=debug message=\"hello world\"\n")
}

func TestGetLogger(t *testing.T) {
	l := &zoup.MemoryWriter{}
	SetLogger(l)

	assert.Equal(t, GetLogger(), l)
}

func TestLog3(t *testing.T) {
	SetLogger(nil)
	datadir.SetHomeDir("/nop")

	panics.Panic(t, func() {
		Log(zoup.DebugLevel, "hello world", zoup.Fields{})
	}, "mkdir /nop: permission denied")
}

func TestLog4(t *testing.T) {
	datadir.UseTestHome(t)

	dir, err := Store.Dir()
	assert.NilError(t, err)

	err = os.Mkdir(dir+"/"+LogFile, 0600)
	assert.NilError(t, err)

	panics.Panic(t, func() {
		Log(zoup.DebugLevel, "hello world", zoup.Fields{})
	}, fmt.Sprintf("open %s: is a directory", dir+"/"+LogFile))
}
