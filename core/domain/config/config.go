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

type Hooks struct {
	Prestart  []string
	Poststart []string
	Prestop   []string
	Poststop  []string
}

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

func (s *Service) String() string {
	return s.Name
}

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
