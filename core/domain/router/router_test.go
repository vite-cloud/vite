package router

import (
	"encoding/json"
	"github.com/docker/docker/client"
	"github.com/vite-cloud/vite/core/domain/config"
	"github.com/vite-cloud/vite/core/domain/deployment"
	"github.com/vite-cloud/vite/core/domain/runtime"
	"gotest.tools/v3/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter_IPFor(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Assert(t, r.URL.Path == "/v1.41/containers/container-id/json")

		b, err := json.Marshal(map[string]any{
			"NetworkSettings": map[string]any{
				"IPAddress": "container-ip",
			},
		})

		assert.NilError(t, err)

		w.Header().Add("Content-Type", "application/json")
		_, err = w.Write(b)
		assert.NilError(t, err)
	}))

	cli, err := client.NewClientWithOpts(client.WithHost(server.URL))
	assert.NilError(t, err)

	docker, err := runtime.NewClient(runtime.WithDockerClient(cli))
	assert.NilError(t, err)

	router := New(&deployment.Deployment{
		Docker: docker,
		Config: &config.Config{
			Services: map[string]*config.Service{
				"test": {
					Name:  "test",
					Image: "nginx:1.15",
					Hosts: []string{"example.com"},
				},
			},
		},
	})
	router.deployment.Add("created_containers", "test", "container-id")

	ip, err := router.IPFor("example.com")
	assert.NilError(t, err)
	assert.Equal(t, ip, "container-ip")
}

func TestHostMatches(t *testing.T) {
	ok, err := hostMatches("example.com", "example.com")
	assert.Assert(t, ok)
	assert.NilError(t, err)

	ok, err = hostMatches("not-example.com", "example.com")
	assert.Assert(t, !ok)
	assert.NilError(t, err)

	ok, err = hostMatches("sub.example.com", "sub.*.example.com")
	assert.Assert(t, !ok)
	assert.NilError(t, err)

	ok, err = hostMatches("sub.smtg.example.com", "sub.*.example.com")
	assert.Assert(t, ok)
	assert.NilError(t, err)

	ok, err = hostMatches("", "\\((\x00")
	assert.Assert(t, !ok)
	assert.Assert(t, err != nil)
}
