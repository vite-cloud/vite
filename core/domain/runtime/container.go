package runtime

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/vite-cloud/vite/core/domain/log"
)

// resource tags available in the manifest
const (
	CreatedContainerManifestKey = "CreatedContainer"
	StartedContainerManifestKey = "StartedContainer"
	StoppedContainerManifestKey = "StoppedContainer"
	RemovedContainerManifestKey = "RemovedContainer"
)

// ContainerCreateOptions defines the options for creating a container
type ContainerCreateOptions struct {
	// Name of the container
	Name string
	// Registry to pull image from, if any
	Registry *types.AuthConfig
	// Env variables to set
	Env []string

	Labels map[string]string

	Networking *network.NetworkingConfig
}

// fullImageName returns the full image name, including registry if any
func fullImageName(image string, registry *types.AuthConfig) string {
	if registry == nil {
		return image
	}

	return fmt.Sprintf("%s/%s", registry.ServerAddress, image)
}

// ContainerCreate creates a container
func (c Client) ContainerCreate(ctx context.Context, image string, opts ContainerCreateOptions) (container.ContainerCreateCreatedBody, error) {
	res, err := c.client.ContainerCreate(ctx, &container.Config{
		Image:  fullImageName(image, opts.Registry),
		Env:    opts.Env,
		Labels: opts.Labels,
	}, &container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
	}, opts.Networking, nil, opts.Name)
	if err != nil {
		return container.ContainerCreateCreatedBody{}, err
	}

	log.Log(log.DebugLevel, "created container", log.Fields{
		"id":            res.ID,
		"image":         image,
		"with_registry": opts.Registry != nil,
	})

	return res, nil
}

// ContainerStart starts a container
func (c Client) ContainerStart(ctx context.Context, ID string) error {
	err := c.client.ContainerStart(ctx, ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	log.Log(log.DebugLevel, "started container", log.Fields{
		"id": ID,
	})

	return nil
}

// ContainerStop stops a container
func (c Client) ContainerStop(ctx context.Context, ID string) error {
	err := c.client.ContainerStop(ctx, ID, nil)
	if err != nil {
		return err
	}

	log.Log(log.DebugLevel, "stopped container", log.Fields{
		"id": ID,
	})

	return nil
}

// ContainerRemove removes a container
func (c Client) ContainerRemove(ctx context.Context, ID string) error {
	err := c.client.ContainerRemove(ctx, ID, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
	if err != nil {
		return err
	}

	log.Log(log.DebugLevel, "removed container", log.Fields{
		"id": ID,
	})

	return nil
}

func (c Client) ContainerExec(ctx context.Context, ID string, command string) error {
	ref, err := c.client.ContainerExecCreate(ctx, ID, types.ExecConfig{
		Cmd: []string{"sh", "-c", command},
	})
	if err != nil {
		return err
	}

	err = c.client.ContainerExecStart(ctx, ref.ID, types.ExecStartCheck{})
	if err != nil {
		return err
	}

	log.Log(log.DebugLevel, "exec command", log.Fields{
		"id":      ID,
		"command": command,
	})

	return nil
}

func (c Client) ContainerInspect(ctx context.Context, ID string) (types.ContainerJSON, error) {
	return c.client.ContainerInspect(ctx, ID)
}
