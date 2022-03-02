package datadir

import (
	"gotest.tools/v3/assert"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestDir(t *testing.T) {
	defer resetDataDir()

	home, err := os.MkdirTemp("", "vite-datadir")
	assert.NilError(t, err)

	homeDir = home

	assert.Equal(t, homeDir+"/"+dataDirName, Dir())
}

func TestSetHomeDir(t *testing.T) {
	defer resetDataDir()

	home, err := os.MkdirTemp("", "vite-datadir")
	assert.NilError(t, err)

	homeDir = home

	assert.Equal(t, homeDir+"/"+dataDirName, Dir())

	newHome := "/tmp/vite-datadir-" + strconv.Itoa(int(time.Now().UnixMilli()))

	SetHomeDir(newHome)

	assert.Equal(t, newHome+"/"+dataDirName, Dir())
}
