package service

import (
	"github.com/redwebcreation/nest/docker"
	"gopkg.in/yaml.v2"
)

// Config represents nest's config
type Config struct {
	Services     ServiceMap  `yaml:"services" json:"services"`
	Registries   RegistryMap `yaml:"registries" json:"registries"`
	ControlPlane struct {
		Host string `yaml:"host" json:"host"`
	} `yaml:"control_plane" json:"controlPlane"`
	Proxy struct {
		HTTP       string `yaml:"http" json:"http"`
		HTTPS      string `yaml:"https" json:"https"`
		SelfSigned bool   `yaml:"self_signed" json:"selfSigned"`
	} `yaml:"proxy" json:"proxy"`
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain Config
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}

	if c.Proxy.HTTP == "" {
		c.Proxy.HTTP = "80"
	}

	if c.Proxy.HTTPS == "" {
		c.Proxy.HTTPS = "443"
	}

	return nil
}

func (c *Config) ExpandIncludes(config *Locator) error {
	for _, s := range c.Services {
		if s.Include == "" {
			continue
		}

		bytes, err := config.Read(s.Include)
		if err != nil {
			return err
		}

		var parsedService *Service

		err = yaml.Unmarshal(bytes, &parsedService)
		if err != nil {
			return err
		}

		parsedService.ApplyDefaults(s.Name)

		c.Services[s.Name] = parsedService
	}

	return nil
}

// RegistryMap maps registry names to their respective docker.Registry structs
type RegistryMap map[string]*docker.Registry

// UnmarshalYAML implements the yaml.Unmarshaler interface
func (r *RegistryMap) UnmarshalYAML(unmarshal func(any) error) error {
	var registries map[string]*docker.Registry
	if err := unmarshal(&registries); err != nil {
		return err
	}

	for name, registry := range registries {
		registry.Name = name
	}

	*r = registries

	return nil
}
