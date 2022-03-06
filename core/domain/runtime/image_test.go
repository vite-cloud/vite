package runtime

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"gotest.tools/v3/assert"
	"os/exec"
	"strings"
	"testing"
)

func TestClient_ImagePull(t *testing.T) {
	image := "alpine:latest"

	_ = exec.Command("docker", "rmi", image).Run()

	out, err := exec.Command("docker", "inspect", image).CombinedOutput()
	assert.Assert(t, err != nil) // assert that image does not exist
	assert.Assert(t, strings.Contains(string(out), "No such object: alpine:latest"))

	client, err := NewClient()
	assert.NilError(t, err)

	err = client.ImagePull(context.Background(), image, ImagePullOptions{})
	assert.NilError(t, err)

	out, err = exec.Command("docker", "inspect", image).CombinedOutput()
	assert.NilError(t, err, "expected image to exist once pulled", string(out))
}

func TestClient_ImagePull2(t *testing.T) {
	image := "nginx:1.21.5"

	client, err := NewClient()
	assert.NilError(t, err)

	err = client.ImagePull(context.Background(), image, ImagePullOptions{})
	assert.NilError(t, err)

	var statuses []string

	expected := []string{
		"Pulling from library/nginx",
		"Digest: sha256:0d17b565c37bcbd895e9d92315a05c1c3c9a29f762b011a10c54a66cd53c9b31",
		"Status: Image is up to date for nginx:1.21.5",
	}

	err = client.ImagePull(context.Background(), image, ImagePullOptions{
		Listener: func(status string) {
			statuses = append(statuses, status)
		},
	})
	assert.NilError(t, err)
	assert.Equal(t, len(statuses), 3)

	for i, status := range statuses {
		assert.Equal(t, status, expected[i])
	}
}

func TestMarshalAuth(t *testing.T) {
	auth := &types.AuthConfig{
		Username: "foo",
		Password: "bar",

		ServerAddress: "https://example.com",
	}

	authJSON, err := json.Marshal(auth)
	assert.NilError(t, err)

	got, err := marshalAuth(auth)
	assert.NilError(t, err)

	assert.Equal(t, base64.URLEncoding.EncodeToString(authJSON), got)
}
