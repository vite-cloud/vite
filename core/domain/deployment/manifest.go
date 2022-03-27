package deployment

import (
	"encoding/json"
	"errors"
	"sync"
)

type LabeledValue struct {
	Label string
	Value any
}

var ErrValueNotFound = errors.New("value not found")

// deploymentJSON is the marshalable representation of a Manifest as it does not rely on sync.Map.
type deploymentJSON struct {
	ID        string
	Resources map[string][]LabeledValue
}

// Add adds a resource to the manifest under a given tag.
func (d *Deployment) Add(key, label string, value any) {
	v, ok := d.resources.Load(key)
	if !ok {
		d.resources.Store(key, []LabeledValue{{label, value}})
		return
	}

	d.resources.Store(key, append(v.([]LabeledValue), LabeledValue{label, value}))
}

// Get returns the resources associated with a given tag.
func (d *Deployment) Get(key string) ([]LabeledValue, error) {
	v, ok := d.resources.Load(key)
	if !ok {
		return nil, errors.New("no resources found matching given key")
	}

	return v.([]LabeledValue), nil
}

// MarshalJSON implements the json.Marshaler interface.
// It takes care of converting the resource map to a marshalable map.
func (d *Deployment) MarshalJSON() ([]byte, error) {
	v := make(map[string][]LabeledValue)

	d.resources.Range(func(key, value any) bool {
		// Add only accepts strings as key, therefore, it is
		// safe to assume that the key is a string.
		v[key.(string)] = value.([]LabeledValue)
		return true
	})

	return json.Marshal(deploymentJSON{
		ID:        d.ID,
		Resources: v,
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It takes care of converting the resources map to a sync.Map.
// Avoid using ints to the map, as they will be converted to float64s.
// Use strings instead.
func (d *Deployment) UnmarshalJSON(data []byte) error {
	var manifestJSON deploymentJSON

	err := json.Unmarshal(data, &manifestJSON)
	if err != nil {
		return err
	}

	d.ID = manifestJSON.ID
	d.resources = sync.Map{}

	for k, v := range manifestJSON.Resources {
		d.resources.Store(k, v)
	}

	return nil
}

func (d *Deployment) Find(key, label string) (any, error) {
	v, ok := d.resources.Load(key)
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
