package config

import (
	"github.com/vite-cloud/vite/core/domain/locator"
	"gopkg.in/yaml.v2"
)

// Config holds vite's configuration.
type Config struct {
	Services map[string]*Service `yaml:"services"`

	Proxy struct {
		HTTP       string `yaml:"http"`
		HTTPS      string `yaml:"https"`
		SelfSigned bool   `yaml:"self_signed"`
	} `yaml:"proxy"`

	ControlPlane struct {
		Host string `yaml:"host"`
	} `yaml:"control_plane"`
}

type Service struct {
	Name string

	Image string

	Hosts []string

	Env []string

	Hooks struct {
		Prestart  []string
		Poststart []string
		Prestop   []string
		Poststop  []string
	}

	// Requires is a list of services that must be running before this service
	Requires []string
}

func (s *Service) String() string {
	return s.Name
}

func Get(l *locator.Locator) (*Config, error) {
	contents, err := l.Read("vite.yaml")
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(contents, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
