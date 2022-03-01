package runtime

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestNewClient(t *testing.T) {
	_, err := NewClient()
	assert.NilError(t, err)
}
