package runtime

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"gotest.tools/v3/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Stats(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1.41/containers/json" {
			containers := []types.Container{
				{
					ID:    "first",
					Names: []string{"/test_container"},
				},
			}

			err := json.NewEncoder(w).Encode(containers)
			assert.NilError(t, err)
			return
		}

		if r.URL.Path == "/v1.41/containers/first/stats" {
			w.Header().Add("Server", "docker/1.41 (linux)")
			stats := &types.Stats{
				MemoryStats: types.MemoryStats{
					MaxUsage: 6651904,
					Usage:    6537216,
					Limit:    67108864,
					Stats: map[string]uint64{
						"cache": 42,
					},
				},

				CPUStats: types.CPUStats{
					CPUUsage: types.CPUUsage{
						PercpuUsage: []uint64{
							8646879,
							24472255,
							36438778,
							30657443,
						},
						UsageInUsermode:   50000000,
						TotalUsage:        100215355,
						UsageInKernelmode: 30000000,
					},
					SystemUsage: 739306590000000,
				},

				PreCPUStats: types.CPUStats{
					CPUUsage: types.CPUUsage{
						PercpuUsage: []uint64{
							8646879,
							24350896,
							36438778,
							30657443,
						},
						UsageInUsermode:   50000000,
						TotalUsage:        100093996,
						UsageInKernelmode: 30000000,
					},
					SystemUsage: 9492140000000,
				},
			}

			err := json.NewEncoder(w).Encode(stats)
			assert.NilError(t, err)
			return
		}

		t.Fatal("unexpected request: ", r.URL.Path)
	}))

	raw, err := client.NewClientWithOpts(client.WithHost(server.URL))
	assert.NilError(t, err)

	cli, err := NewClient(WithDockerClient(raw))
	assert.NilError(t, err)

	stats, err := cli.Stats(context.Background(), types.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.Arg("name", "testclient_stats"),
		),
	})
	assert.NilError(t, err)

	assert.Equal(t, len(stats), 1)
	assert.Equal(t, stats[0].Name, "test_container")
	assert.Equal(t, stats[0].ID, "first")
	assert.Equal(t, stats[0].MemoryUsed, uint64(6537174))
	assert.Equal(t, stats[0].MemoryAvailable, uint64(67108864))
	assert.Equal(t, stats[0].MemoryUsage, (6537174.0/67108864.0)*100)
	assert.Equal(t, stats[0].CPUCount, 4)
	assert.Equal(t, stats[0].CPUDelta, uint64(100215355-100093996))
	assert.Equal(t, stats[0].CPUUsage, (121359.0/729814450000000.0)*100*4)
}
