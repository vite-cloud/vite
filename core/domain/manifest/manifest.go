package manifest

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/vite-cloud/vite/core/domain/datadir"
)

// Store is the manifest store.
const Store = datadir.Store("manifest")

type contextKey string

type LabeledValue struct {
	Label string
	Value any
}

// ContextKey is the key used to store the manifest in a context.
// It prevents overlapping with other libraries that might set a "manifest" key
// in the same context.
var ContextKey = contextKey("manifest")

var ErrValueNotFound = errors.New("value not found")

// Manifest is a set of tagged resources for a given deployment.
// The zero Manifest is empty and ready for use.
type Manifest struct {
	// Version is the version of the manifest's deployment.
	Version string

	// resources is a map of tags associated with resources.
	resources sync.Map
}

// manifestJSON is the marshalable representation of a Manifest as it does not rely on sync.Map.
type manifestJSON struct {
	Version   string
	Resources map[string][]LabeledValue
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
func (m *Manifest) Add(key, label string, value any) {
	v, ok := m.resources.Load(key)
	if !ok {
		m.resources.Store(key, []LabeledValue{{label, value}})
		return
	}

	m.resources.Store(key, append(v.([]LabeledValue), LabeledValue{label, value}))
}

// Get returns the resources associated with a given tag.
func (m *Manifest) Get(key string) ([]LabeledValue, error) {
	v, ok := m.resources.Load(key)
	if !ok {
		return nil, errors.New("no resources found matching given key")
	}

	return v.([]LabeledValue), nil
}

// MarshalJSON implements the json.Marshaler interface.
// It takes care of converting the resource map to a marshalable map.
func (m *Manifest) MarshalJSON() ([]byte, error) {
	v := make(map[string][]LabeledValue)

	m.resources.Range(func(key, value any) bool {
		// Add only accepts strings as key, therefore, it is
		// safe to assume that the key is a string.
		v[key.(string)] = value.([]LabeledValue)
		return true
	})

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

func (m *Manifest) Find(key, label string) (any, error) {
	v, ok := m.resources.Load(key)
	if !ok {
		return nil, ErrValueNotFound
	}

	for _, lv := range v.([]LabeledValue) {
		if lv.Label == label {
			return lv.Value, nil
		}
	}

	return nil, ErrValueNotFound
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
