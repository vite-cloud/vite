package runtime

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"gotest.tools/v3/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_RegistryLogin(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.Path, "/v1.41/auth")

		body, err := io.ReadAll(r.Body)
		assert.NilError(t, err)

		var auth types.AuthConfig

		err = json.Unmarshal(body, &auth)
		assert.NilError(t, err)

		assert.Equal(t, auth.ServerAddress, "registry.example.com")
		assert.Equal(t, auth.Username, "nobody")
		assert.Equal(t, auth.Password, "password")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"identitytoken": "token"}`))
	}))

	raw, err := client.NewClientWithOpts(client.WithHost(server.URL))
	assert.NilError(t, err)

	cli, err := NewClient(WithDockerClient(raw))
	assert.NilError(t, err)

	err = cli.RegistryLogin(context.Background(), types.AuthConfig{
		ServerAddress: "registry.example.com",
		Username:      "nobody",
		Password:      "password",
	})
	assert.NilError(t, err)
}
