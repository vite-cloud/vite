package datadir

import (
	"gotest.tools/v3/assert"
	"os"
	"strconv"
	"testing"
	"time"
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

func TestPath3(t *testing.T) {
	defer resetDataDir()

	home, err := os.MkdirTemp("", "vite-datadir")
	assert.NilError(t, err)

	homeDir = home

	// ensure that an error is returned if it can not create the directory
	err = os.WriteFile(Dir()+"/test", []byte{}, 0644)
	assert.NilError(t, err)

	_, err = Store("test").Path("hello.world")
	assert.ErrorContains(t, err, "no such file or directory")
}

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
