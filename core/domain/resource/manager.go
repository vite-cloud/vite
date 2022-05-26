package resource

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vite-cloud/vite/core/domain/datadir"
	"os"
)

func Save[T any](store datadir.Store, result T, name func(T) string) error {
	dir, err := store.Dir()
	if err != nil {
		return err
	}

	contents, err := json.Marshal(result)
	if err != nil {
		return err
	}

	return os.WriteFile(fmt.Sprintf("%s/%s.json", dir, name(result)), contents, 0644)
}

// List returns a list of all the results in the store.
func List[T any](store datadir.Store) ([]*T, error) {
	var results []*T

	dir, err := store.Dir()
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

		f, err := store.Open(entry.Name(), os.O_RDONLY, 0)
		if err != nil {
			return nil, err
		}

		defer f.Close()

		var m T

		err = json.NewDecoder(f).Decode(&m)
		if err != nil {
			return nil, err
		}

		results = append(results, &m)
	}

	return results, nil
}

// Delete removes the manifest from the Store and returns an error if it fails.
// It does not return an error if the manifest does not exist.
func Delete[T any](store datadir.Store, entry T, name func(T) string) error {
	dir, err := store.Dir()
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%v.json", dir, name(entry))

	if _, err = os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return err
	}

	return os.Remove(path)
}

// Get returns the manifest for a given version or os.ErrNotExist if it does not exist.
func Get[T any](store datadir.Store, name any) (*T, error) {
	f, err := store.Open(fmt.Sprintf("%v.json", name), os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	var e T

	err = json.NewDecoder(f).Decode(&e)
	return &e, err
}
