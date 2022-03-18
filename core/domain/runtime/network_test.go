package runtime

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/vite-cloud/vite/core/domain/manifest"
	"gotest.tools/v3/assert"
	"strconv"
	"testing"
	"time"
)

func TestClient_NetworkCreate(t *testing.T) {
	t.Parallel()

	raw, err := client.NewClientWithOpts(client.FromEnv)
	assert.NilError(t, err)

	cli, err := NewClient()
	assert.NilError(t, err)

	ctx := context.Background()
	ctx = context.WithValue(ctx, manifest.ContextKey, &manifest.Manifest{})

	name := "test_" + strconv.Itoa(int(time.Now().UnixMilli()))

	res, err := cli.NetworkCreate(ctx, name, NetworkCreateOptions{})
	assert.NilError(t, err)

	ins, err := raw.NetworkInspect(ctx, name, types.NetworkInspectOptions{})
	assert.NilError(t, err)

	assert.Equal(t, res, ins.ID)
}
