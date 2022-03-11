package deployment

import (
	"github.com/vite-cloud/vite/core/domain/config"
	"gotest.tools/v3/assert"
	"sort"
	"testing"
)

func TestFlatten(t *testing.T) {
	type Test struct {
		Services map[string][]string `json:"services"`
		Expected [][]string          `json:"sorted"`
		Cyclic   bool                `json:"cyclic"`
	}

	tests := []Test{
		{
			Services: map[string][]string{
				"example": {},
			},
			Expected: [][]string{
				{"example"},
			},
		},
		{
			Services: map[string][]string{
				"a": {"b"},
				"b": {"a"},
			},
			Cyclic: true,
		},
		{
			Services: map[string][]string{
				"a": {"b"},
				"b": {"c"},
				"c": {"a"},
			},
			Cyclic: true,
		},
		{
			Services: map[string][]string{
				"laravel": {"mysql", "redis", "elastic"},
				"redis":   {},
				"mysql":   {"fs"},
				"elastic": {"minio"},
				"minio":   {"fs"},
				"fs":      {},
			},
			Expected: [][]string{
				{"fs"},
				{"minio"},
				{"mysql", "redis", "elastic"},
				{"laravel"},
			},
		},
		{
			Services: map[string][]string{
				"example":  {"mysql"},
				"mysql":    {"fast-dfs", "logger"},
				"fast-dfs": {},
				"logger":   {},
			},
			Expected: [][]string{
				{"fast-dfs", "logger"},
				{"mysql"},
				{"example"},
			},
		},
	}

	for _, test := range tests {
		serviceMap := ServiceMap{}

		for name, dependencies := range test.Services {
			service := &config.Service{
				Requires: dependencies,
			}

			serviceMap[name] = service
		}

		sorted, err := serviceMap.Layered()

		if test.Cyclic {
			assert.ErrorContains(t, err, "circular dependency detected:")

			continue
		}

		assert.NilError(t, err)
		assert.Equal(t, len(sorted), len(test.Expected))

		for kl, layer := range test.Expected {
			assert.Equal(t, len(layer), len(sorted[kl]))
			sort.Strings(layer)

			for ks, service := range layer {
				assert.DeepEqual(t, sorted[kl][ks], serviceMap[service])
			}
		}
	}
}
