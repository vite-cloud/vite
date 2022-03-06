package locator

import (
	"encoding/json"
	"github.com/vite-cloud/vite/core/domain/datadir"
	"os"
	"path/filepath"
	"strings"
)

// configStore is the storage used by locator to store cloned configs.
const configStore = datadir.Store("locator")

const configFile = "config.json"

// Locator contains the configuration for the locator.
type Locator struct {
	Provider   Provider `json:"provider"`
	Protocol   string   `json:"protocol"`
	Repository string   `json:"repository"`
	Branch     string   `json:"branch"`
	Commit     string   `json:"commit"`
	Path       string   `json:"path"`
}

// Read a file from the repository.
func (l *Locator) Read(file string) ([]byte, error) {
	git, err := l.git()
	if err != nil {
		return nil, err
	}

	if !git.RepoExists() {
		err = git.Clone(l.Provider.URL(l.Protocol, l.Repository), l.Branch)
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

func (l *Locator) Save() error {
	dir, err := configStore.Dir()
	if err != nil {
		return err
	}

	contents, err := json.Marshal(l)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(dir, configFile), contents, 0600)
}

// LoadFromStore loads a Locator from a config.json in store or fails if it does not exist.
func LoadFromStore() (*Locator, error) {
	f, err := configStore.Open(configFile, os.O_CREATE|os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	var locator Locator
	err = json.NewDecoder(f).Decode(&locator)
	if err != nil {
		return nil, err
	}

	return &locator, nil
}
