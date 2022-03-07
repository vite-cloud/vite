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
	Image string `yaml:"image"`

	Hosts []string `yaml:"hosts"`

	Env []string `yaml:"env"`

	Hooks struct {
		Prestart  []string `yaml:"prestart"`
		Poststart []string `yaml:"poststart"`
		Prestop   []string `yaml:"prestop"`
		Poststop  []string `yaml:"poststop"`
	} `yaml:"hooks"`
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
