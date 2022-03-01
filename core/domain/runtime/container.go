package runtime

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/vite-cloud/vite/core/domain/log"
)

// ContainerCreateOptions defines the options for creating a container
type ContainerCreateOptions struct {
	// Name of the container
	Name string
	// Image to use
	Image string
	// Registry to pull image from, if any
	Registry *types.AuthConfig
	// Env variables to set
	Env []string
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
		Image: fullImageName(image, opts.Registry),
		Env:   opts.Env,
	}, nil, nil, nil, opts.Name)
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
func (c Client) ContainerStart(ctx context.Context, id string) error {
	err := c.client.ContainerStart(ctx, id, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	log.Log(log.DebugLevel, "started container", log.Fields{
		"id": id,
	})

	return nil
}

// ContainerStop stops a container
func (c Client) ContainerStop(ctx context.Context, id string) error {
	err := c.client.ContainerStop(ctx, id, nil)
	if err != nil {
		return err
	}

	log.Log(log.DebugLevel, "stopped container", log.Fields{
		"id": id,
	})

	return nil
}
