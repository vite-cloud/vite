package runtime

import (
	"context"
	"gotest.tools/v3/assert"
	"os/exec"
	"strings"
	"testing"
)

func TestClient_ImagePull(t *testing.T) {
	image := "alpine:latest"

	_ = exec.Command("docker", "rmi", image).Run()

	out, err := exec.Command("docker", "inspect", image).CombinedOutput()
	assert.Assert(t, err != nil)
	assert.Assert(t, strings.Contains(string(out), "No such object: alpine:latest"))

	client, err := NewClient()
	assert.NilError(t, err)

	err = client.ImagePull(context.Background(), image, ImagePullOptions{})
	assert.NilError(t, err)

	_, err = exec.Command("docker", "inspect", image).CombinedOutput()
	assert.NilError(t, err, "expected image to exist once pulled")
}
