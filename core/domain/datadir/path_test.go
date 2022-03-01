package datadir

import (
	"gotest.tools/v3/assert"
	"os"
	"testing"
)

func TestPath(t *testing.T) {
	defer resetDataDir()

	home, err := os.MkdirTemp("", "vite-datadir")
	assert.NilError(t, err)

	homeDir = home

	path, err := Store("test").Path("hello.world")
	assert.NilError(t, err)

	if path != homeDir+"/"+dataDirName+"/test/hello.world" {
		t.Fatal("path is not correct")
	}

	_, err = os.Stat(homeDir + "/" + dataDirName + "/test")
	assert.NilError(t, err)

	// Ensure that it creates the parent directory but not the file
	_, err = os.Stat(path)
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestPath2(t *testing.T) {
	defer resetDataDir()

	home, err := os.MkdirTemp("", "vite-datadir")
	assert.NilError(t, err)

	homeDir = home

	_, err = Store("test").Path("../../../../..", "etc", "passwd")
	assert.Error(t, err, "path is not in data dir")
}

func TestDir(t *testing.T) {
	defer resetDataDir()

	home, err := os.MkdirTemp("", "vite-datadir")
	assert.NilError(t, err)

	homeDir = home

	assert.Equal(t, homeDir+"/"+dataDirName, Dir())
}
