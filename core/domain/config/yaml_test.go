package config

import (
	"github.com/docker/docker/api/types"
	"gotest.tools/v3/assert"
	"testing"
)

func TestConfigYAML_ToConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		yaml    *configYAML
		want    *Config
		wantErr bool
	}{
		{
			name: "it sets the service's name",
			yaml: &configYAML{
				Services: map[string]*serviceYAML{
					"example": {
						Image: "nginx:latest",
					},
				},
			},
			want: &Config{
				Services: map[string]*Service{
					"example": {
						IsTopLevel: true,
						Name:       "example",
						Image:      "nginx:latest",
					},
				},
			},
		},
		{
			name: "it sets the service's hosts",
			yaml: &configYAML{
				Services: map[string]*serviceYAML{
					"example": {
						Hosts: []string{"example.com", "example.org"},
					},
				},
			},
			want: &Config{
				Services: map[string]*Service{
					"example": {
						IsTopLevel: true,
						Name:       "example",
						Hosts:      []string{"example.com", "example.org"},
					},
				},
			},
		},
		{
			name: "it expands hosts starting with a tilde",
			yaml: &configYAML{
				Services: map[string]*serviceYAML{
					"example": {
						Hosts: []string{"~example.com"},
					},
				},
			},
			want: &Config{
				Services: map[string]*Service{
					"example": {
						IsTopLevel: true,
						Name:       "example",
						Hosts: []string{
							"example.com",
							"www.example.com",
						},
					},
				},
			},
		},
		{
			name: "it sets the service's hooks",
			yaml: &configYAML{
				Services: map[string]*serviceYAML{
					"example": {
						Hooks: Hooks{
							Prestart:  []string{"prestart_hook"},
							Poststart: []string{"poststart_hook1", "poststart_hook2"},
							Prestop:   []string{"prestop_hook"},
							Poststop:  []string{"poststop_hook1", "poststop_hook2"},
						},
					},
				},
			},
			want: &Config{
				Services: map[string]*Service{
					"example": {
						IsTopLevel: true,
						Name:       "example",
						Hooks: Hooks{
							Prestart:  []string{"prestart_hook"},
							Poststart: []string{"poststart_hook1", "poststart_hook2"},
							Prestop:   []string{"prestop_hook"},
							Poststop:  []string{"poststop_hook1", "poststop_hook2"},
						},
					},
				},
			},
		},
		{
			name: "it sets the service's environment variables",
			yaml: &configYAML{
				Services: map[string]*serviceYAML{
					"example": {
						Env: []string{"FOO=bar", "BAR=baz"},
					},
				},
			},
			want: &Config{
				Services: map[string]*Service{
					"example": {
						IsTopLevel: true,
						Name:       "example",
						Env: []string{
							"FOO=bar",
							"BAR=baz",
						},
					},
				},
			},
		},
		{
			name: "it resolve the registry name to the actual registry",
			yaml: &configYAML{
				Services: map[string]*serviceYAML{
					"example": {
						Registry: "docker",
					},
				},
				Registries: map[string]*registryYAML{
					"docker": {
						Host: "docker.io",
					},
				},
			},
			want: &Config{
				Services: map[string]*Service{
					"example": {
						IsTopLevel: true,
						Name:       "example",
						Registry: &types.AuthConfig{
							ServerAddress: "docker.io",
						},
					},
				},
			},
		},
		{
			name: "it sets the service's registry",
			yaml: &configYAML{
				Services: map[string]*serviceYAML{
					"example": {
						Registry: &registryYAML{
							Host:     "registry.vite.cloud",
							Username: "foo",
							Password: "bar",
						},
					},
				},
			},
			want: &Config{
				Services: map[string]*Service{
					"example": {
						IsTopLevel: true,
						Name:       "example",
						Registry: &types.AuthConfig{
							ServerAddress: "registry.vite.cloud",
							Username:      "foo",
							Password:      "bar",
						},
					},
				},
			},
		},
		{
			name: "it fails if the registry name points to an unknown registry",
			yaml: &configYAML{
				Services: map[string]*serviceYAML{
					"example": {
						Registry: "nop",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := test.yaml.ToConfig()
			assert.Assert(t, (err != nil) == test.wantErr, "unexpected err: %v", err)
			assert.DeepEqual(t, test.want, got)
		})
	}
}
