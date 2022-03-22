package locator

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/vite-cloud/vite/core/domain/datadir"
)

const (
	// ConfigStore is the storage used by locator to store cloned configs.
	ConfigStore = datadir.Store("locator")
	// ConfigFile is the name of the config file for the locator.
	ConfigFile = "config.json"
)

// ErrInvalidCommit is returned when the commit is invalid.
var ErrInvalidCommit = errors.New("invalid commit")

// Locator contains the configuration for the locator.
type Locator struct {
	Provider   Provider `json:"provider,omitempty"`
	Protocol   string   `json:"protocol,omitempty"`
	Repository string   `json:"repository,omitempty"`
	Branch     string   `json:"branch,omitempty"`
	Commit     string   `json:"commit,omitempty"`
	Path       string   `json:"path,omitempty"`
}

// Read a file from the repository.
func (l *Locator) Read(file string) ([]byte, error) {
	if l.Commit == "" {
		return nil, ErrInvalidCommit
	}

	git, err := l.git()
	if err != nil {
		return nil, err
	}

	if !git.RepoExists() {
		err = git.Clone(l.Provider.URL(l.Protocol, l.Repository), l.Branch)
		if errors.Is(err, ErrEmptyBranch) {
			return nil, fmt.Errorf("could not clone repository %s: no branch specified (run `vite setup` again)", l.Provider.URL(l.Protocol, l.Repository))
		}
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
	dir, err := ConfigStore.Dir()
	if err != nil {
		return "", err
	}

	path := dir + "/" + l.Branch + "-" + strings.Replace(l.Repository, "/", "-", -1)

	return Git(path), nil
}

// Save the locator to the config store.
func (l *Locator) Save() error {
	dir, err := ConfigStore.Dir()
	if err != nil {
		return err
	}

	contents, _ := json.Marshal(l)

	return os.WriteFile(filepath.Join(dir, ConfigFile), contents, 0600)
}

// LoadFromStore loads a Locator from a config.json in store or fails if it does not exist.
func LoadFromStore() (*Locator, error) {
	f, err := ConfigStore.Open(ConfigFile, os.O_CREATE|os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var locator Locator
	err = json.NewDecoder(f).Decode(&locator)
	if errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("config locator hasn't been configured yet, run `vite setup` first")
	} else if err != nil {
		return nil, err
	}

	return &locator, nil
}

// Checksum returns the unique id of the locator, it may be used as a file name.
func (l *Locator) Checksum() string {
	return base64.StdEncoding.EncodeToString([]byte(l.Branch + l.Repository + l.Provider.Name() + l.Commit + l.Path))
}

func (l *Locator) Commits() (CommitList, error) {
	git, err := l.git()
	if err != nil {
		return nil, err
	}

	return git.Commits(l.Branch)
}
