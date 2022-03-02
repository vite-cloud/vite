package manifest

import (
	"encoding/json"
	"fmt"
	"github.com/vite-cloud/vite/core/domain/datadir"
	"os"
	"sync"
)

const Store = datadir.Store("manifest")

type Manifest struct {
	Version string

	Resources sync.Map
}

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

func (m *Manifest) Add(key string, value any) {
	v, ok := m.Resources.Load(key)
	if !ok {
		m.Resources.Store(key, []any{value})
		return
	}

	m.Resources.Store(key, append(v.([]any), value))
}

func (m *Manifest) MarshalJSON() ([]byte, error) {
	v := make(map[string]any)

	valid := true

	m.Resources.Range(func(key, value any) bool {
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

	return json.Marshal(v)
}
