package datadir

import (
	"errors"
	"github.com/vite-cloud/vite/pkg/homedir"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	homeDir     string
	dataDir     string
	initDataDir = new(sync.Once)
)

const (
	configFileDir = ".vite"
)

func getHomeDir() string {
	if homeDir == "" {
		homeDir = homedir.Get()
	}

	return homeDir
}

// resetDataDir is used in testing to reset the "dataDir" package variable
// and its sync.Once to force re-lookup between tests.
func resetDataDir() {
	dataDir = ""
	initDataDir = new(sync.Once)
}

func Dir() string {
	initDataDir.Do(setDataDir)
	return dataDir
}

func setDataDir() {
	if dataDir != "" {
		return
	}

	dataDir = getHomeDir() + "/" + configFileDir

	if err := os.MkdirAll(dataDir, 0700); err != nil {
		panic(err)
	}
}

func Path(id string, p ...string) (string, error) {
	path := filepath.Join(append([]string{Dir(), id}, p...)...)
	if !strings.HasPrefix(path, Dir()+"/") {
		return "", errors.New("path is not in data dir")
	}

	if err := os.MkdirAll(filepath.Join(Dir(), id), 0700); err != nil {
		return "", err
	}

	return path, nil
}
