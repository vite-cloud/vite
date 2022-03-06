package runtime

import (
	"context"
	"github.com/docker/docker/client"
	"github.com/vite-cloud/vite/core/domain/log"
	"github.com/vite-cloud/vite/core/domain/manifest"
	"gotest.tools/v3/assert"
	"strconv"
	"testing"
	"time"
)

const testImage = "nginx:1.21.5"

func TestClient_ContainerCreate(t *testing.T) {
	logger := &log.MemoryWriter{}
	log.SetLogger(logger)

	ctx, cli, raw := createTestEnv(t)

	name := uniqueContainerName()

	err := cli.ImagePull(ctx, testImage, ImagePullOptions{})
	assert.NilError(t, err)

	body, err := cli.ContainerCreate(ctx, testImage, ContainerCreateOptions{
		Name: name,
		Env: []string{
			"A=B",
			"B=\"hello world\"",
		},
	})
	assert.NilError(t, err)

	ins, err := raw.ContainerInspect(ctx, body.ID)
	assert.NilError(t, err)

	assert.Equal(t, ins.Name, "/"+name)
	assert.Equal(t, ins.Config.Image, testImage)
	assert.Assert(t, len(ins.Config.Env) >= 2)
	assert.Assert(t, contains(ins.Config.Env, "A=B"))
	assert.Assert(t, contains(ins.Config.Env, "B=\"hello world\""))
	assert.Assert(t, ins.HostConfig.RestartPolicy.IsAlways())

	assert.Equal(t, logger.Len(), 1)
	assert.Equal(t, logger.Last().Message, "created container")
	assert.Equal(t, logger.Last().Level, log.DebugLevel)
	assert.Equal(t, logger.Last().Fields["id"], body.ID)
	assert.Equal(t, logger.Last().Fields["image"], testImage)
	assert.Equal(t, logger.Last().Fields["with_registry"], false)

	created, err := ctx.Value("manifest").(*manifest.Manifest).Get(CreatedContainerManifestKey)
	assert.NilError(t, err)

	assert.Equal(t, len(created), 1)
	assert.Equal(t, created[0], body.ID)
}

func TestClient_ContainerStart(t *testing.T) {
	logger := &log.MemoryWriter{}
	log.SetLogger(logger)

	ctx, cli, raw := createTestEnv(t)

	name := uniqueContainerName()

	err := cli.ImagePull(ctx, testImage, ImagePullOptions{})
	assert.NilError(t, err)

	body, err := cli.ContainerCreate(ctx, testImage, ContainerCreateOptions{
		Name: name,
	})
	assert.NilError(t, err)

	err = cli.ContainerStart(ctx, body.ID)
	assert.NilError(t, err)

	ins, err := raw.ContainerInspect(ctx, body.ID)
	assert.NilError(t, err)

	assert.Equal(t, ins.ID, body.ID)
	assert.Equal(t, ins.State.Status, "running")

	assert.Equal(t, logger.Len(), 2)
	assert.Equal(t, logger.Last().Message, "started container")
	assert.Equal(t, logger.Last().Level, log.DebugLevel)
	assert.Equal(t, logger.Last().Fields["id"], body.ID)

	started, err := ctx.Value("manifest").(*manifest.Manifest).Get(StartedContainerManifestKey)
	assert.NilError(t, err)

	assert.Equal(t, len(started), 1)
	assert.Equal(t, started[0], body.ID)
}

func TestClient_ContainerStop(t *testing.T) {
	logger := &log.MemoryWriter{}
	log.SetLogger(logger)

	ctx, cli, raw := createTestEnv(t)

	name := uniqueContainerName()

	err := cli.ImagePull(ctx, testImage, ImagePullOptions{})
	assert.NilError(t, err)

	body, err := cli.ContainerCreate(ctx, testImage, ContainerCreateOptions{
		Name: name,
	})
	assert.NilError(t, err)

	err = cli.ContainerStart(ctx, body.ID)
	assert.NilError(t, err)

	err = cli.ContainerStop(ctx, body.ID)
	assert.NilError(t, err)

	ins, err := raw.ContainerInspect(ctx, body.ID)
	assert.NilError(t, err)

	assert.Equal(t, ins.ID, body.ID)
	assert.Equal(t, ins.State.Status, "exited")

	assert.Equal(t, logger.Len(), 3)
	assert.Equal(t, logger.Last().Message, "stopped container")
	assert.Equal(t, logger.Last().Level, log.DebugLevel)
	assert.Equal(t, logger.Last().Fields["id"], body.ID)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func createTestEnv(t *testing.T) (context.Context, *Client, *client.Client) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "manifest", &manifest.Manifest{})

	cli, err := NewClient()
	assert.NilError(t, err)

	raw, err := client.NewClientWithOpts(client.FromEnv)
	assert.NilError(t, err)

	return ctx, cli, raw
}

func uniqueContainerName() string {
	return "test_" + strconv.Itoa(int(time.Now().UnixMilli()))
}
