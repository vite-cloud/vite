package runtime

import (
	"gotest.tools/v3/assert"
	"gotest.tools/v3/env"
	"testing"
)

func TestNewClient(t *testing.T) {
	clientInstance = nil

	cli, err := NewClient()
	assert.NilError(t, err)

	assert.Assert(t, cli == clientInstance)

	old := clientInstance

	cli, err = NewClient()
	assert.NilError(t, err)

	assert.Assert(t, cli == old)
}

func TestNewClient2(t *testing.T) {
	clientInstance = nil

	defer env.Patch(t, "DOCKER_HOST", "invalid")()

	_, err := NewClient()
	assert.ErrorContains(t, err, "unable to parse docker host `invalid`")
}
