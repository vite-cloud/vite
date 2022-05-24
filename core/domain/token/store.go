package token

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vite-cloud/vite/core/domain/datadir"
	"os"
	"time"
)

const Store = datadir.Store("tokens")

type Token struct {
	Label      string
	Value      string
	CreatedAt  time.Time
	LastUsedAt time.Time
}

func (t *Token) Save() error {
	dir, err := Store.Dir()
	if err != nil {
		return err
	}

	contents, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return os.WriteFile(fmt.Sprintf("%s/%s.json", dir, t.Label), contents, 0644)
}

// List returns a list of all the manifests in the Store.
func List() ([]*Token, error) {
	var manifests []*Token

	dir, err := Store.Dir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			return nil, fmt.Errorf("manifest store is corrupted: %s is a directory", entry.Name())
		}

		f, err := Store.Open(entry.Name(), os.O_RDONLY, 0)
		if err != nil {
			return nil, err
		}

		defer f.Close()

		var m Token

		err = json.NewDecoder(f).Decode(&m)
		if err != nil {
			return nil, err
		}

		manifests = append(manifests, &m)
	}

	return manifests, nil
}

// Delete removes the manifest from the Store and returns an error if it fails.
// It does not return an error if the manifest does not exist.
func Delete(ID int64) error {
	dir, err := Store.Dir()
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%d.json", dir, ID)

	if _, err = os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return err
	}

	return os.Remove(path)
}

// Get returns the manifest for a given version or os.ErrNotExist if it does not exist.
func Get(ID int64) (*Token, error) {
	f, err := Store.Open(fmt.Sprintf("%d.json", ID), os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	var m Token

	err = json.NewDecoder(f).Decode(&m)
	return &m, err
}
