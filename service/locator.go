package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redwebcreation/nest/loggy"
	"io/fs"
	"log"
	"os"
	"strings"
)

// Locator contains nest's configuration
type Locator struct {
	Provider   string `json:"provider"`
	Repository string `json:"repository"`
	Branch     string `json:"branch"`
	Commit     string `json:"commit"`
	// StoreDir is the path where services config are stored
	StoreDir string `json:"-"`
	// Path is the location of the config file
	Path   string      `json:"-"`
	Logger *log.Logger `json:"-"`
	Git    *Git        `json:"-"`
}

func (l *Locator) StorePath() string {
	return l.StoreDir + "/" + l.Branch + "-" + strings.Replace(l.Repository, "/", "-", -1)
}

func (l *Locator) RemoteURL() string {
	return fmt.Sprintf("git@%s.com:%s.git", l.Provider, l.Repository)
}

func (l *Locator) Read(file string) ([]byte, error) {
	configPath := l.StorePath()

	_, err := os.Stat(configPath)

	if errors.Is(err, fs.ErrNotExist) {
		err = l.Clone()
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	l.log(loggy.DebugLevel, "reading services config file", loggy.Fields{
		"tag":  "Config.read",
		"file": file,
	})

	return l.Git.ReadFile(configPath, l.Commit, file)
}

func (l *Locator) Save() error {
	contents, err := json.Marshal(l)
	if err != nil {
		return err
	}
	err = os.WriteFile(l.Path, contents, 0600)
	if err != nil {
		return err
	}

	l.log(loggy.InfoLevel, "updating config", loggy.Fields{
		"tag": "config.update",
	})

	return nil
}

func (l *Locator) LoadCommit(commit string) error {
	l.Commit = commit

	return l.Save()
}

func (l *Locator) Clone() error {
	_ = os.RemoveAll(l.StorePath())

	err := l.Git.Clone(l.RemoteURL(), l.StorePath(), l.Branch)

	if err != nil {
		return err
	}

	l.log(loggy.InfoLevel, "cloned config", loggy.Fields{
		"tag": "config.clone",
	})

	return nil
}

func (l *Locator) log(level loggy.Level, message string, fields loggy.Fields) {
	fields["commit"] = l.Commit
	fields["branch"] = l.Branch
	fields["location"] = l.RemoteURL()

	l.Logger.Print(loggy.NewEvent(level, message, fields))
}

func (l *Locator) Pull() error {
	_, err := l.Git.Pull(l.StorePath(), l.Branch)

	if err != nil {
		return err
	}

	l.log(loggy.InfoLevel, "pulled config", loggy.Fields{
		"tag": "config.pull",
	})

	return nil
}
