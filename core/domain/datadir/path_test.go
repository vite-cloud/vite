package datadir

import (
	"os"
	"strconv"
	"testing"
	"time"

	panics "github.com/magiconair/properties/assert"

	"github.com/docker/docker/pkg/homedir"
	"gotest.tools/v3/assert"
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

func TestStore_Dir(t *testing.T) {
	defer resetDataDir()

	home, err := os.MkdirTemp("", "vite-datadir")
	assert.NilError(t, err)

	SetHomeDir(home)

	_, err = os.Stat(Dir() + "/this")
	assert.ErrorIs(t, err, os.ErrNotExist)

	dir, err := Store("this").Dir()
	assert.NilError(t, err)

	_, err = os.Stat(dir)
	assert.NilError(t, err)

	assert.Equal(t, dir, Dir()+"/this")
}

func TestStore_Open(t *testing.T) {
	defer resetDataDir()

	home, err := os.MkdirTemp("", "vite-datadir")
	assert.NilError(t, err)

	SetHomeDir(home)

	_, err = os.Stat(Dir() + "/this")
	assert.ErrorIs(t, err, os.ErrNotExist)

	f, err := Store("this").Open("file", os.O_RDWR|os.O_CREATE, 0600)
	assert.NilError(t, err)

	_, err = os.Stat(Dir() + "/this/file")
	assert.NilError(t, err)

	f.Close()
}

func TestStore_Open2(t *testing.T) {
	defer resetDataDir()

	// can not use SetHomeDir as it calls setDataDir
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

	setDataDir()

	assert.Equal(t, "/something", dataDir)
}

func TestSetHomeDir2(t *testing.T) {
	defer resetDataDir()

	panics.Panic(t, func() {
		SetHomeDir("/nop")

		Dir()
	}, "mkdir /nop: permission denied")
}
