package runtime

import (
	"gotest.tools/v3/assert"
	"gotest.tools/v3/env"
	"testing"
)

func TestNewClient(t *testing.T) {
	_, err := NewClient()
	assert.NilError(t, err)
}

func TestNewClient2(t *testing.T) {
	defer env.Patch(t, "DOCKER_HOST", "invalid")()

	_, err := NewClient()
	assert.ErrorContains(t, err, "unable to parse docker host `invalid`")
}
