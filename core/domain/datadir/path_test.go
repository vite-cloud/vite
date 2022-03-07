package datadir

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/docker/docker/pkg/homedir"
	"gotest.tools/v3/assert"
)

func TestDir(t *testing.T) {
	defer resetDataDir()

	home, err := os.MkdirTemp("", "vite-datadir")
	assert.NilError(t, err)

	homeDir = home

	dir, err := Dir()
	assert.NilError(t, err)

	assert.Equal(t, homeDir+"/"+dataDirName, dir)
}

func TestSetHomeDir(t *testing.T) {
	defer resetDataDir()

	home, err := os.MkdirTemp("", "vite-datadir")
	assert.NilError(t, err)

	homeDir = home

	dir, err := Dir()
	assert.NilError(t, err)

	assert.Equal(t, homeDir+"/"+dataDirName, dir)

	newHome := "/tmp/vite-datadir-" + strconv.Itoa(int(time.Now().UnixMilli()))

	SetHomeDir(newHome)

	dir, err = Dir()
	assert.NilError(t, err)

	assert.Equal(t, newHome+"/"+dataDirName, dir)
}

func TestStore_Dir(t *testing.T) {
	defer resetDataDir()

	UseTestHome(t)

	dir, err := Dir()
	assert.NilError(t, err)

	_, err = os.Stat(dir + "/this")
	assert.ErrorIs(t, err, os.ErrNotExist)

	path, err := Store("this").Dir()
	assert.NilError(t, err)

	_, err = os.Stat(path)
	assert.NilError(t, err)

	assert.Equal(t, path, dir+"/this")
}

func TestStore_Dir2(t *testing.T) {
	defer resetDataDir()

	SetHomeDir("/nop")

	_, err := Store("this").Dir()
	assert.ErrorIs(t, err, os.ErrPermission)
}

func TestStore_Open(t *testing.T) {
	defer resetDataDir()

	UseTestHome(t)

	dir, err := Dir()
	assert.NilError(t, err)

	_, err = os.Stat(dir + "/this")
	assert.ErrorIs(t, err, os.ErrNotExist)

	f, err := Store("this").Open("file", os.O_RDWR|os.O_CREATE, 0600)
	assert.NilError(t, err)

	_, err = os.Stat(dir + "/this/file")
	assert.NilError(t, err)

	f.Close()
}

func TestStore_Open2(t *testing.T) {
	defer resetDataDir()

	// can not use setHomeDir as it calls setDataDir
	// which panics if it can not create the directory
	homeDir = "/nop"
	dataDir = "/nop/.vite"

	_, err := Store("this").Open("file", os.O_RDWR|os.O_CREATE, 0600)
	assert.ErrorIs(t, err, os.ErrPermission)
}

func TestStore_String(t *testing.T) {
	assert.Equal(t, Store("something").String(), "something")
}

func TestGetHomeDir(t *testing.T) {
	homeDir = ""

	home := homedir.Get()

	assert.Equal(t, home, getHomeDir())
}

func TestSetDataDir(t *testing.T) {
	defer resetDataDir()

	// ensures that setDataDir does not update dataDir if it is already set.
	dataDir = "/something"

	err := setDataDir()
	assert.NilError(t, err)

	assert.Equal(t, "/something", dataDir)
}

func TestSetHomeDir2(t *testing.T) {
	defer resetDataDir()

	SetHomeDir("/nop")

	_, err := Dir()
	assert.ErrorIs(t, err, os.ErrPermission)
}
