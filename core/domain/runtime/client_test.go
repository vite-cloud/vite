package runtime

import (
	"gotest.tools/v3/assert"
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
