package config

import (
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/vite-cloud/vite/core/domain/locator"
	"gopkg.in/yaml.v2"
)

// Config holds vite's configuration.
type Config struct {
	Services map[string]*Service `json:"services"`

	Proxy struct {
		HTTP  string `json:"http"`
		HTTPS string `json:"https"`
	} `json:"proxy"`

	ControlPlane struct {
		Host string `json:"host"`
	} `json:"controlPlane"`

	Locator *locator.Locator `json:"locator"`
}

// Service contains the configuration about a service.
type Service struct {
	// IsTopLevel indicates whether this service is depended on by other services.
	IsTopLevel bool `json:"isTopLevel"`

	// Name is the service's name. It need not contain the registry host.
	Name string `json:"name"`

	// Image is the service's Docker image.
	Image string `json:"image"`

	// Hosts are a list of hosts to which the service answers to.
	Hosts []string `json:"hosts"`

	// Env is a list of environment variables to set.
	Env []string `json:"env"`

	// Hooks are the service's hooks: prestart, poststart, prestop, poststop.
	Hooks Hooks `json:"hooks"`

	// Requires is a list of services that must be running before this service
	Requires []*Service `json:"requires"`

	// Registry is the auth configuration for the service's registry.
	Registry *types.AuthConfig `yaml:"registry"`
}

var configCache = make(map[string]*Config)

// GetUsingDefaultLocator returns the Config given a config locator.Locator.
func GetUsingDefaultLocator() (*Config, error) {
	l, err := locator.LoadFromStore()
	if err != nil {
		return nil, err
	}

	return Get(l)
}

func Get(l *locator.Locator) (*Config, error) {
	if _, ok := configCache[l.Checksum()]; ok {
		return configCache[l.Checksum()], nil
	}

	contents, err := l.Read("vite.yaml")
	if errors.Is(err, locator.ErrInvalidCommit) {
		return nil, fmt.Errorf("could not read the config, no commit specified: run `vite use` to pick one")
	} else if err != nil {
		return nil, err
	}

	var c configYAML
	err = yaml.Unmarshal(contents, &c)
	if err != nil {
		return nil, err
	}

	converted, err := c.ToConfig()
	if err != nil {
		return nil, err
	}

	converted.Locator = l

	configCache[l.Checksum()] = converted

	return converted, nil
}
