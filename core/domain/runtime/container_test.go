package runtime

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/vite-cloud/vite/core/domain/log"
	"gotest.tools/v3/assert"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"
)

const testImage = "nginx:1.21.5"

type testCtx struct {
	t             *testing.T
	ctx           context.Context
	cli           *Client
	raw           *client.Client
	logger        *log.MemoryWriter
	containerName string
}

func TestClient(t *testing.T) {
	tests := []struct {
		name string
		test func(ctx *testCtx)
	}{
		{
			name: "it can start containers",
			test: testContainerStart,
		},
		{
			name: "it can create containers",
			test: testContainerCreate,
		},
		{
			name: "it can stop containers",
			test: testContainerStop,
		},
		{
			name: "it can remove containers",
			test: testContainerRemove,
		},
	}

	cli, err := NewClient()
	assert.NilError(t, err)

	raw, err := client.NewClientWithOpts(client.FromEnv)
	assert.NilError(t, err)

	err = cli.ImagePull(context.Background(), testImage, ImagePullOptions{})
	assert.NilError(t, err)

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			logger := &log.MemoryWriter{}
			log.SetLogger(logger)

			test.test(&testCtx{
				t:             t,
				cli:           cli,
				raw:           raw,
				ctx:           context.Background(),
				logger:        logger,
				containerName: "test_" + strconv.Itoa(int(time.Now().UnixMilli())),
			})
		})
	}

	testContainers, err := raw.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.Arg("name", "test_"),
		),
	})
	assert.NilError(t, err)

	var wg sync.WaitGroup
	for _, container := range testContainers {
		container := container

		wg.Add(1)

		go func() {
			defer wg.Done()

			_ = raw.ContainerStop(context.Background(), container.ID, nil)

			err = raw.ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{
				Force: true,
			})
			assert.NilError(t, err)
		}()
	}
	wg.Wait()
}

func testContainerCreate(tc *testCtx) {
	body, err := tc.cli.ContainerCreate(tc.ctx, testImage, ContainerCreateOptions{
		Name: tc.containerName,
		Env: []string{
			"A=B",
			"B=\"hello world\"",
		},
	})
	assert.NilError(tc.t, err)

	ins, err := tc.raw.ContainerInspect(tc.ctx, body.ID)
	assert.NilError(tc.t, err)

	assert.Equal(tc.t, ins.Name, "/"+tc.containerName)
	assert.Equal(tc.t, ins.Config.Image, testImage)
	assert.Assert(tc.t, len(ins.Config.Env) >= 2)

	sort.Strings(ins.Config.Env)

	assert.Equal(tc.t, ins.Config.Env[0], "A=B")
	assert.Equal(tc.t, ins.Config.Env[1], "B=\"hello world\"")
	assert.Assert(tc.t, ins.HostConfig.RestartPolicy.IsAlways())

	assert.Equal(tc.t, tc.logger.Len(), 1)
	assert.Equal(tc.t, tc.logger.Last().Message, "created container")
	assert.Equal(tc.t, tc.logger.Last().Level, log.DebugLevel)
	assert.Equal(tc.t, tc.logger.Last().Fields["id"], body.ID)
	assert.Equal(tc.t, tc.logger.Last().Fields["image"], testImage)
	assert.Equal(tc.t, tc.logger.Last().Fields["with_registry"], false)

	//created, err := tc.ctx.Value(manifest.ContextKey).(*manifest.Manifest).Get(CreatedContainerManifestKey)
	//assert.NilError(tc.t, err)
	//
	//assert.Equal(tc.t, len(created), 1)
	//assert.Equal(tc.t, created[0], body.ID)
}

func testContainerStart(tc *testCtx) {
	body, err := tc.cli.ContainerCreate(tc.ctx, testImage, ContainerCreateOptions{
		Name: tc.containerName,
	})
	assert.NilError(tc.t, err)

	err = tc.cli.ContainerStart(tc.ctx, body.ID)
	assert.NilError(tc.t, err)

	ins, err := tc.raw.ContainerInspect(tc.ctx, body.ID)
	assert.NilError(tc.t, err)

	assert.Equal(tc.t, ins.ID, body.ID)
	assert.Equal(tc.t, ins.State.Status, "running")

	assert.Equal(tc.t, tc.logger.Len(), 2)
	assert.Equal(tc.t, tc.logger.Last().Message, "started container")
	assert.Equal(tc.t, tc.logger.Last().Level, log.DebugLevel)
	assert.Equal(tc.t, tc.logger.Last().Fields["id"], body.ID)

	//started, err := tc.ctx.Value(manifest.ContextKey).(*manifest.Manifest).Get(StartedContainerManifestKey)
	//assert.NilError(tc.t, err)
	//
	//assert.Equal(tc.t, len(started), 1)
	//assert.Equal(tc.t, started[0], body.ID)
}

func testContainerStop(tc *testCtx) {
	body, err := tc.cli.ContainerCreate(tc.ctx, testImage, ContainerCreateOptions{
		Name: tc.containerName,
	})
	assert.NilError(tc.t, err)

	err = tc.cli.ContainerStart(tc.ctx, body.ID)
	assert.NilError(tc.t, err)

	err = tc.cli.ContainerStop(tc.ctx, body.ID)
	assert.NilError(tc.t, err)

	ins, err := tc.raw.ContainerInspect(tc.ctx, body.ID)
	assert.NilError(tc.t, err)

	assert.Equal(tc.t, ins.ID, body.ID)
	assert.Equal(tc.t, ins.State.Status, "exited")

	assert.Equal(tc.t, tc.logger.Len(), 3)
	assert.Equal(tc.t, tc.logger.Last().Message, "stopped container")
	assert.Equal(tc.t, tc.logger.Last().Level, log.DebugLevel)
	assert.Equal(tc.t, tc.logger.Last().Fields["id"], body.ID)

	//stopped, err := tc.ctx.Value(manifest.ContextKey).(*manifest.Manifest).Get(StoppedContainerManifestKey)
	//assert.NilError(tc.t, err)
	//
	//assert.Equal(tc.t, len(stopped), 1)
	//assert.Equal(tc.t, stopped[0], body.ID)
}

func testContainerRemove(tc *testCtx) {
	body, err := tc.cli.ContainerCreate(tc.ctx, testImage, ContainerCreateOptions{
		Name: tc.containerName,
	})
	assert.NilError(tc.t, err)

	err = tc.cli.ContainerStart(tc.ctx, body.ID)
	assert.NilError(tc.t, err)

	err = tc.cli.ContainerStop(tc.ctx, body.ID)
	assert.NilError(tc.t, err)

	err = tc.cli.ContainerRemove(tc.ctx, body.ID)
	assert.NilError(tc.t, err)

	_, err = tc.raw.ContainerInspect(tc.ctx, body.ID)
	assert.Error(tc.t, err, "Error: No such container: "+body.ID)

	assert.Equal(tc.t, tc.logger.Len(), 4)
	assert.Equal(tc.t, tc.logger.Last().Message, "removed container")
	assert.Equal(tc.t, tc.logger.Last().Level, log.DebugLevel)
	assert.Equal(tc.t, tc.logger.Last().Fields["id"], body.ID)

	//removed, err := tc.ctx.Value(manifest.ContextKey).(*manifest.Manifest).Get(RemovedContainerManifestKey)
	//assert.NilError(tc.t, err)
	//
	//assert.Equal(tc.t, len(removed), 1)
	//assert.Equal(tc.t, removed[0].(string), body.ID)
}
