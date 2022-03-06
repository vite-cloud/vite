package manifest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vite-cloud/vite/core/domain/datadir"
	"os"
	"sync"
)

// Store is the manifest store.
const Store = datadir.Store("manifest")

// Manifest is a set of tagged resources for a given deployment.
type Manifest struct {
	// Version is the version of the manifest's deployment.
	Version string

	// resources is a map of tags associated with resources.
	resources sync.Map
}

// manifestJSON is the marshalable representation of a Manifest as it does not rely on sync.Map.
type manifestJSON struct {
	Version   string
	Resources map[string]any
}

// Save writes the manifest to the Store.
func (m *Manifest) Save() error {
	f, err := Store.Open(m.Version+".json", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	contents, err := json.Marshal(m)
	if err != nil {
		return err
	}

	_, err = f.Write(contents)
	if err != nil {
		return err
	}

	return f.Close()
}

// Add adds a resource to the manifest under a given tag.
func (m *Manifest) Add(key string, value any) {
	v, ok := m.resources.Load(key)
	if !ok {
		m.resources.Store(key, []any{value})
		return
	}

	m.resources.Store(key, append(v.([]any), value))
}

// Get returns the resources associated with a given tag.
func (m *Manifest) Get(key string) ([]any, error) {
	v, ok := m.resources.Load(key)
	if !ok {
		return nil, errors.New("no resources found")
	}

	return v.([]any), nil
}

// MarshalJSON implements the json.Marshaler interface.
// It takes care of converting the resource map to a marshalable map.
func (m *Manifest) MarshalJSON() ([]byte, error) {
	v := make(map[string]any)

	valid := true

	m.resources.Range(func(key, value any) bool {
		k, ok := key.(string)
		if !ok {
			valid = false
			return false
		}

		v[k] = value
		return true
	})

	if !valid {
		return nil, fmt.Errorf("manifest.MarshalJSON: invalid manifest key (must be string)")
	}

	return json.Marshal(manifestJSON{
		Version:   m.Version,
		Resources: v,
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It takes care of converting the resources map to a sync.Map.
// Avoid using ints to the map, as they will be converted to float64s.
// Use strings instead.
func (m *Manifest) UnmarshalJSON(data []byte) error {
	var manifestJSON manifestJSON

	err := json.Unmarshal(data, &manifestJSON)
	if err != nil {
		return err
	}

	m.Version = manifestJSON.Version
	m.resources = sync.Map{}

	for k, v := range manifestJSON.Resources {
		m.resources.Store(k, v)
	}

	return nil
}

// List returns a list of all the manifests in the Store.
func List() ([]*Manifest, error) {
	var manifests []*Manifest

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

		var manifest Manifest

		err = json.NewDecoder(f).Decode(&manifest)
		if err != nil {
			return nil, err
		}

		manifests = append(manifests, &manifest)
	}

	return manifests, nil
}

// Delete removes the manifest from the Store and returns an error if it fails.
// It does not return an error if the manifest does not exist.
func Delete(version string) error {
	dir, err := Store.Dir()
	if err != nil {
		return err
	}

	path := dir + "/" + version + ".json"

	if _, err = os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return err
	}

	return os.Remove(path)
}

// Get returns the manifest for a given version or os.ErrNotExist if it does not exist.
func Get(version string) (*Manifest, error) {
	f, err := Store.Open(version+".json", os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	var manifest Manifest

	err = json.NewDecoder(f).Decode(&manifest)
	return &manifest, err
}
