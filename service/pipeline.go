package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/vite-cloud/vite/docker"
	"strings"
	"time"
)

type Event struct {
	Service *Service
	Value   any
}

type Pipeline struct {
	Deployment      *Deployment
	Service         *Service
	HasDependencies bool
}

func (p *Pipeline) Log(v any) {
	p.Deployment.Events <- Event{p.Service, v}
}

func (p Pipeline) Run() error {
	if p.HasDependencies {
		net, err := p.CreateServiceNetwork()
		if err != nil {
			return err
		}

		err = p.ConnectRequiredServices(net)
		if err != nil {
			return err
		}
	}

	err := p.PullImage()
	if err != nil {
		return err
	}
	id, err := p.CreateContainer()
	if err != nil {
		return err
	}

	err = p.RunHooks(id, p.Service.Hooks.Prestart)
	if err != nil {
		return err
	}

	err = p.StartContainer(id)
	if err != nil {
		return err
	}

	err = p.RunHooks(id, p.Service.Hooks.Poststart)
	if err != nil {
		return err
	}

	err = p.EnsureContainerIsRunning(id)
	if err != nil {
		if err2 := p.Deployment.Docker.ContainerDelete(id); err2 != nil {
			return fmt.Errorf("%s (cleanup failed: %s)", err, err2)
		}

		return err
	}

	p.Log("deployment ended")

	return nil
}

func (p *Pipeline) PullImage() error {
	image := docker.Image(p.Service.Image)

	return p.Deployment.Docker.ImagePull(image, func(event *docker.PullEvent) {
		p.Log(event.Status)
	}, p.Deployment.ServicesConfig.Registries[p.Service.Registry])
}

func (p *Pipeline) CreateServiceNetwork() (string, error) {
	name := fmt.Sprintf("%s_%s", p.Service.Name, p.Deployment.ID)

	net, err := p.Deployment.Docker.NetworkCreate(name, map[string]string{
		"cloud.vite.service":       p.Service.Name,
		"cloud.vite.deployment_id": p.Deployment.ID,
	})

	if err != nil {
		return "", err
	}

	p.Deployment.Manifest.Networks[p.Service.Name] = net

	return net, nil
}

func (p *Pipeline) ConnectRequiredServices(networkID string) error {
	for _, require := range p.Service.Requires {
		err := p.Deployment.Docker.NetworkConnect(networkID, p.Deployment.Manifest.Containers[require], []string{require})

		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Pipeline) ContainerName() string {
	return "vite_" + p.Service.Name + "_" + strings.Replace(p.Service.Image, ":", "_", 1) + "_" + p.Deployment.ID
}

func (p *Pipeline) CreateContainer() (string, error) {
	containerName := p.ContainerName()

	var networking *network.NetworkingConfig

	if p.HasDependencies {
		networking = &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				p.Service.Name: {
					NetworkID: p.Deployment.Manifest.Networks[p.Service.Name],
				},
			},
		}
	}

	return p.Deployment.Docker.ContainerCreate(&container.Config{
		Image: p.Service.Image,
		Labels: map[string]string{
			"cloud.vite.service":       p.Service.Name,
			"cloud.vite.deployment_id": p.Deployment.ID,
		},
		Env: p.Service.Env.ForDocker(),
	}, &container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
	}, networking, containerName)
}

func (p *Pipeline) RunHooks(containerID string, commands []string) error {
	for _, command := range commands {
		err := p.Deployment.Docker.ContainerExec(containerID, command)
		if err != nil {
			return err
		}

		p.Log("ran command: " + command)
	}

	return nil
}

func (p *Pipeline) StartContainer(containerID string) error {
	err := p.Deployment.Docker.ContainerStart(containerID)
	if err != nil {
		return err
	}

	p.Deployment.Manifest.Containers[p.Service.Name] = containerID

	return nil
}

var (
	ErrContainerNotRunning = errors.New("container is not running")
	ErrContainerTimeout    = errors.New("container is not running (timeout)")
)

// EnsureContainerIsRunning will wait for the container to start and then return
// an error if the container is not running after either :
// - 10 seconds if the container has no health-check
// - Retries * (Interval + Timeout) if the container has a health-check
//
// todo(pipeline): return logs from failed container
func (p *Pipeline) EnsureContainerIsRunning(containerID string) error {
	info, err := p.Deployment.Docker.Client.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return err
	}

	var timeout time.Duration

	if info.Config.Healthcheck == nil {
		timeout = 10 * time.Second
	} else {
		seconds := float64(info.Config.Healthcheck.Retries) * (info.Config.Healthcheck.Interval + info.Config.Healthcheck.Timeout).Seconds()

		timeout = time.Duration(seconds) * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return ErrContainerTimeout
		case <-time.After(250 * time.Millisecond):
			info, err = p.Deployment.Docker.Client.ContainerInspect(ctx, containerID)
			if err != nil {
				return err
			}

			if info.RestartCount > 0 || info.State.ExitCode != 0 {
				return ErrContainerNotRunning
			}

			hasHealthchecks := info.State.Health != nil

			if !hasHealthchecks {
				if info.State.Status == "running" {
					return nil
				}

				return ErrContainerNotRunning
			}

			if info.State.Health.Status == types.Healthy {
				return nil
			}

			if info.State.Health.Status == types.Unhealthy {
				return ErrContainerNotRunning
			}
		}
	}
}
