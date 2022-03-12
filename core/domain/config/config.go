package config

import (
	"github.com/docker/docker/api/types"
	"github.com/vite-cloud/vite/core/domain/locator"
	"gopkg.in/yaml.v2"
)

// Config holds vite's configuration.
type Config struct {
	Services map[string]*Service

	Proxy struct {
		HTTP       string
		HTTPS      string
		SelfSigned bool
	}

	ControlPlane struct {
		Host string
	}
}

// Hooks contains a service's hooks.
type Hooks struct {
	Prestart  []string
	Poststart []string
	Prestop   []string
	Poststop  []string
}

// Service contains the configuration about a service.
type Service struct {
	Name string

	Image string

	Hosts []string

	Env []string

	Hooks Hooks

	// Requires is a list of services that must be running before this service
	Requires []*Service

	Registry *types.AuthConfig `yaml:"registry"`
}

// Get returns the Config given a config locator.Locator.
func Get(l *locator.Locator) (*Config, error) {
	contents, err := l.Read("vite.yaml")
	if err != nil {
		return nil, err
	}

	var c configYAML
	err = yaml.Unmarshal(contents, &c)
	if err != nil {
		return nil, err
	}

	return c.toConfig()
}
