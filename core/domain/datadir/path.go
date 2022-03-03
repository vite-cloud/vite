package datadir

import (
	"github.com/docker/docker/pkg/homedir"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	// homeDir is the home directory for the current user.
	homeDir string
	// dataDir is the data directory for the current user. (usually ~/.vite)
	dataDir string
	// initDataDir is used to ensure that dataDir is initialized only once.
	initDataDir = new(sync.Once)
)

// dataDirName is the name of the data directory.
const dataDirName = ".vite"

// Store contains the name of the subdirectory in the data directory
// For example Store(certs) would return ~/.vite/certs
type Store string

// String returns the string representation of the Store
func (s Store) String() string {
	return string(s)
}

// Open is a convenience method to open a file from the current Store
func (s Store) Open(path string, flags int, perm os.FileMode) (*os.File, error) {
	dir, err := s.Dir()
	if err != nil {
		return nil, err
	}

	return os.OpenFile(filepath.Join(dir, path), flags, perm)
}

// Dir returns the store directory for the current user.
func (s Store) Dir() (string, error) {
	path := filepath.Join(Dir(), s.String())

	if err := os.MkdirAll(path, 0700); err != nil {
		return "", err
	}

	return path, nil
}

// Dir returns the path to the data directory for the current user
// Usually, this is ~/.vite
func Dir() string {
	initDataDir.Do(setDataDir)
	return dataDir
}

// SetHomeDir sets the home directory for the current user
func SetHomeDir(home string) {
	homeDir = home

	resetDataDir()
}

// getHomeDir caches and returns the home directory for the current user.
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

// setDataDir sets the data directory (usually, ~/.vite), if not already set.
// It may create the directory if it does not exist yet.
func setDataDir() {
	if dataDir != "" {
		return
	}

	dataDir = strings.TrimRight(getHomeDir(), "/") + "/" + dataDirName

	if err := os.MkdirAll(dataDir, 0700); err != nil {
		panic(err)
	}
}