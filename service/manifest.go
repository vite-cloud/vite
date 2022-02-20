package service

import (
	"encoding/json"
	"fmt"
	"os"
)

type Manifest struct {
	ID         string            `json:"id"`
	Containers map[string]string `json:"containers"`
	Networks   map[string]string `json:"networks"`
}

// ManifestManager contains the path to the manifest file and methods to manage manifests.
type ManifestManager struct {
	Path string
}

func (m ManifestManager) NewManifest(id string) *Manifest {
	return &Manifest{
		ID:         id,
		Containers: make(map[string]string),
		Networks:   make(map[string]string),
	}
}

func (m ManifestManager) LoadWithID(path string) (*Manifest, error) {
	bytes, err := os.ReadFile(m.Path + "/" + path)
	if err != nil {
		return nil, err
	}

	var manifest Manifest
	err = json.Unmarshal(bytes, &manifest)
	if err != nil {
		return nil, err
	}

	return &manifest, nil
}

func (m ManifestManager) Latest() (*Manifest, error) {
	manifests, err := os.ReadDir(m.Path)
	if err != nil {
		return nil, err
	}

	if len(manifests) == 0 {
		return nil, fmt.Errorf("no manifest found: run `vite deploy`")
	}

	latest := manifests[len(manifests)-1].Name()

	return m.LoadWithID(latest)
}

func (m ManifestManager) Save(manifest *Manifest) error {
	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return os.WriteFile(m.Path+"/"+manifest.ID+".json", bytes, 0600)
}
