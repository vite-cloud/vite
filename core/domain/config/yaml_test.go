package config

import (
	"testing"

	"github.com/docker/docker/api/types"
	"gotest.tools/v3/assert"
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
		{
			name: "it fails if the registry is an invalid type",
			yaml: &configYAML{
				Services: map[string]*serviceYAML{
					"example": {
						Registry: 45.4,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "a service can require another one",
			yaml: &configYAML{
				Services: map[string]*serviceYAML{
					"first": {
						Requires: []string{"second"},
					},
					"second": {},
				},
			},
			want: &Config{
				Services: map[string]*Service{
					"first": {
						IsTopLevel: true,
						Name:       "first",
						Requires: []*Service{
							{
								IsTopLevel: false,
								Name:       "second",
							},
						},
					},
					"second": {
						IsTopLevel: false,
						Name:       "second",
					},
				},
			},
		},
		{
			name: "it fails if service.Requires contains an unknown service",
			yaml: &configYAML{
				Services: map[string]*serviceYAML{
					"example": {
						Requires: []string{"something"},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "it fails if a dependency is wrong",
			yaml: &configYAML{
				Services: map[string]*serviceYAML{
					"a": {
						Requires: []string{"b"},
					},
					"b": {
						Requires: []string{"c"}, // c does not exist
					},
				},
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		test := test

		t.Run("service only:"+test.name, func(t *testing.T) {
			t.Parallel()

			got, err := test.yaml.ToConfig()
			assert.Assert(t, (err != nil) == test.wantErr, "unexpected err: %v", err)
			if err == nil {
				assert.DeepEqual(t, test.want.Services, got.Services)
			}
		})
	}
}

func TestConfigYAML_ToConfig2(t *testing.T) {
	// it can handle circular dependency
	config := &configYAML{
		Services: map[string]*serviceYAML{
			"a": {
				Requires: []string{"b"},
			},
			"b": {
				Requires: []string{"a"},
			},
		},
	}
	got, err := config.ToConfig()
	assert.NilError(t, err)

	assert.Assert(t, len(got.Services) == 2)

	assert.Equal(t, got.Services["a"].Name, "a")
	assert.Equal(t, len(got.Services["a"].Requires), 1)
	assert.Equal(t, got.Services["a"].Requires[0], got.Services["b"])

	assert.Equal(t, got.Services["b"].Name, "b")
	assert.Equal(t, len(got.Services["b"].Requires), 1)
	assert.Equal(t, got.Services["b"].Requires[0], got.Services["a"])
}

func TestConfigYAML_ToConfig3(t *testing.T) {
	// it can require itself (a special case of circular dependency handling)
	config := &configYAML{
		Services: map[string]*serviceYAML{
			"a": {
				Requires: []string{"a"},
			},
		},
	}

	got, err := config.ToConfig()
	assert.NilError(t, err)

	assert.Assert(t, len(got.Services) == 1)

	assert.Equal(t, got.Services["a"].Name, "a")
	assert.Equal(t, len(got.Services["a"].Requires), 1)
	assert.Equal(t, got.Services["a"].Requires[0], got.Services["a"])

}
