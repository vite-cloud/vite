package runtime

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
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

	name := "test_" + strconv.Itoa(int(time.Now().UnixMilli()))

	res, err := cli.NetworkCreate(ctx, name, NetworkCreateOptions{})
	assert.NilError(t, err)

	ins, err := raw.NetworkInspect(ctx, name, types.NetworkInspectOptions{})
	assert.NilError(t, err)
	assert.Equal(t, res, ins.ID)

	//created, err := ctx.Value(manifest.ContextKey).(*manifest.Manifest).Get(CreatedNetworkManifestKey)
	//assert.NilError(t, err)

	//assert.Equal(t, len(created), 1)
	//assert.Equal(t, created[0], res)

	err = cli.NetworkRemove(ctx, res)
	assert.NilError(t, err)
}

func TestClient_NetworkCreate2(t *testing.T) {
	t.Parallel()

	cli, err := NewClient()
	assert.NilError(t, err)

	ctx := context.Background()

	name := "test_" + strconv.Itoa(int(time.Now().UnixMilli()))

	res, err := cli.NetworkCreate(ctx, name, NetworkCreateOptions{})
	assert.NilError(t, err)

	_, err = cli.NetworkCreate(ctx, name, NetworkCreateOptions{})
	assert.Error(t, err, fmt.Sprintf("Error response from daemon: network with name %s already exists", name), "err: %s", err)

	err = cli.NetworkRemove(ctx, res)
	assert.NilError(t, err)
}
