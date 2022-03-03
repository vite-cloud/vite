package config

import (
	"github.com/vite-cloud/vite/core/domain/datadir"
	"path/filepath"
	"strings"
)

// configStore is the storage used by locator to store cloned configs.
const configStore = datadir.Store("locator")

// Locator contains the configuration for the locator.
type Locator struct {
	Provider   Provider
	UseHTTPS   bool
	Repository string
	Branch     string
	Commit     string
	Path       string
}

// protocolName returns the protocol name, either ssh or https.
func (l *Locator) protocolName() string {
	if l.UseHTTPS {
		return "https"
	}
	return "ssh"
}

// Read a file from the repository.
func (l *Locator) Read(file string) ([]byte, error) {
	git, err := l.git()
	if err != nil {
		return nil, err
	}

	if !git.RepoExists() {
		err = git.Clone(l.Provider.URL(!l.UseHTTPS, l.Repository), l.Branch)
		if err != nil {
			return nil, err
		}
	}

	contents, err := git.Read(l.Commit, filepath.Join(l.Path, file))
	if err != nil {
		return nil, err
	}

	return contents, nil
}

// git returns a Git object for the locator.
func (l *Locator) git() (Git, error) {
	dir, err := configStore.Dir()
	if err != nil {
		return "", err
	}

	path := dir + "/" + l.Branch + "-" + strings.Replace(l.Repository, "/", "-", -1)

	return Git(path), nil
}
