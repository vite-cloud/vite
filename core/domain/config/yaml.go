package config

import (
	"fmt"
	"github.com/docker/docker/api/types"
)

// configYAML is the YAML representation of the config.
type configYAML struct {
	Services map[string]*serviceYAML `yaml:"services"`

	Registries map[string]*registryYAML `yaml:"registries"`

	Proxy struct {
		HTTP  string `yaml:"http"`
		HTTPS string `yaml:"https"`

		SelfSigned bool `yaml:"self_signed"`
	} `yaml:"proxy"`

	ControlPlaneHost string `yaml:"control_plane_host"`

	configRegistry map[string]*types.AuthConfig
	configServices map[string]*Service
}

// serviceYAML is the YAML representation of a service
type serviceYAML struct {
	Image string `yaml:"image"`

	Hosts []string `yaml:"hosts"`

	Env []string `yaml:"env"`

	Hooks struct {
		Prestart  []string `yaml:"prestart"`
		Poststart []string `yaml:"poststart"`
		Prestop   []string `yaml:"prestop"`
		Poststop  []string `yaml:"poststop"`
	} `yaml:"hooks"`

	Requires []string `yaml:"requires"`

	Registry any `yaml:"registry"`
}

// registryYAML is the YAML representation of a registry
type registryYAML struct {
	// Username is the username for the registry
	Username string `yaml:"username,omitempty"`
	// Password is the password to the registry
	Password string `yaml:"password,omitempty"`

	// Host is the address of the registry server.
	Host string `yaml:"host,omitempty"`
	// IdentityToken is used to authenticate the user and get
	// an access token for the registry.
	IdentityToken string `yaml:"identity_token,omitempty"`
	// RegistryToken is a bearer token to be sent to a registry
	RegistryToken string `yaml:"registry_token,omitempty"`
}

func (c configYAML) toConfig() (*Config, error) {
	config := &Config{}

	for name, service := range c.Services {
		converted, err := c.toConfigService(name, service)
		if err != nil {
			return nil, err
		}

		config.Services[name] = converted
	}

	return config, nil
}

func (c *configYAML) toConfigService(name string, s *serviceYAML) (*Service, error) {
	if _, ok := c.configServices[name]; ok {
		return c.configServices[name], nil
	}

	service := &Service{
		Name:  name,
		Image: s.Image,
		Hosts: s.Hosts,
		Env:   s.Env,
		Hooks: Hooks{
			Prestart:  s.Hooks.Prestart,
			Poststart: s.Hooks.Poststart,
			Prestop:   s.Hooks.Prestop,
			Poststop:  s.Hooks.Poststop,
		},
	}

	if s.Registry != nil {
		switch s.Registry.(type) {
		case string:
			if _, ok := c.Registries[s.Registry.(string)]; !ok {
				return nil, fmt.Errorf("registry %s not found", s.Registry.(string))
			}

			service.Registry = c.toConfigRegistry(c.Registries[s.Registry.(string)])
		case *registryYAML:
			registry := s.Registry.(*registryYAML)

			service.Registry = c.toConfigRegistry(registry)
		default:
			return nil, fmt.Errorf("invalid registry type %T (%v)", s.Registry, s.Registry)
		}
	}

	for _, require := range s.Requires {
		if _, ok := c.Services[require]; !ok {
			return nil, fmt.Errorf("service %s not found, %s can not depend on it", require, name)
		}

		converted, err := c.toConfigService(require, c.Services[require])
		if err != nil {
			return nil, err
		}

		service.Requires = append(service.Requires, converted)
	}

	c.configServices[name] = service

	return service, nil
}

func (c configYAML) toConfigRegistry(r *registryYAML) *types.AuthConfig {
	registry := &types.AuthConfig{
		Username:      r.Username,
		Password:      r.Password,
		ServerAddress: r.Host,
		IdentityToken: r.IdentityToken,
		RegistryToken: r.RegistryToken,
	}

	return registry
}
