package deployment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"os"
	"sync"
	"time"

	"github.com/vite-cloud/vite/core/domain/config"
	"github.com/vite-cloud/vite/core/domain/runtime"
)

// Deployment holds the information needed to deploy a service.
type Deployment struct {
	ID     string
	Docker *runtime.Client

	Bus       chan<- Event
	resources sync.Map
}

func (d *Deployment) RunHooks(ctx context.Context, containerID string, commands []string) error {
	for _, command := range commands {
		err := d.Docker.ContainerExec(ctx, containerID, command)
		if err != nil {
			return err
		}
	}

	return nil
}

// Deploy deploys a service.
func (d *Deployment) Deploy(ctx context.Context, events chan<- Event, service *config.Service) error {
	if service.IsTopLevel && len(service.Requires) > 0 {
		subnetter, err := runtime.NewSubnetManager()
		if err != nil {
			return err
		}

		subnet, err := subnetter.Next()
		if err != nil {
			return err
		}

		events <- Event{
			ID:      AcquireSubnet,
			Service: service,
			Data:    fmt.Sprintf("Assigned subnet %s to the service's network", subnet.String()),
		}

		networkID, err := d.Docker.NetworkCreate(ctx, fmt.Sprintf("%s_%s", service.Name, d.ID), runtime.NetworkCreateOptions{
			IPAM: &network.IPAM{
				Driver: "default",
				Config: []network.IPAMConfig{
					{
						Subnet: subnet.String(),
					},
				},
			},
		})
		if err != nil {
			return err
		}

		events <- Event{
			ID:      CreateNetwork,
			Service: service,
			Data:    fmt.Sprintf("Created network %s", networkID),
		}
		d.Add("network", service.Name, networkID)

		for _, require := range service.Requires {
			id, err := d.Find("container", require.Name)
			if err != nil {
				return err
			}

			if err = d.Docker.NetworkConnect(ctx, networkID, id.(string)); err != nil {
				return err
			}

			events <- Event{
				ID:      ConnectDependency,
				Service: service,
				Data:    fmt.Sprintf("Connected service %s to the service's network", require.Name),
			}

		}
	}

	err := d.Docker.ImagePull(ctx, service.Image, runtime.ImagePullOptions{
		Auth: service.Registry,
	})
	if err != nil {
		return err
	}

	events <- Event{
		ID:      PullImage,
		Service: service,
		Data:    service.Image,
	}

	var networking *network.NetworkingConfig

	// Connect the container to its network.
	if service.IsTopLevel && len(service.Requires) > 0 {
		// We can ignore the error as we know the network exists, as we created it above.
		net, _ := d.Find("network", service.Name)

		networking = &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				service.Name: {
					NetworkID: net.(string),
				},
			},
		}
	}

	ref, err := d.Docker.ContainerCreate(ctx, service.Image, runtime.ContainerCreateOptions{
		Name:     d.ID + "_" + service.Name,
		Env:      service.Env,
		Registry: service.Registry,
		Labels: map[string]string{
			"cloud.vite.service":    service.Name,
			"cloud.vite.deployment": d.ID,
		},
		Networking: networking,
	})
	if err != nil {
		return err
	}

	events <- Event{
		ID:      CreateContainer,
		Service: service,
	}
	d.Add("created_containers", service.Name, ref)

	err = d.RunHooks(ctx, ref.ID, service.Hooks.Prestart)
	if err != nil {
		return err
	}

	events <- Event{
		ID:      RunHook,
		Service: service,
		Data:    service.Hooks.Prestart,
	}

	err = d.Docker.ContainerStart(ctx, ref.ID)
	if err != nil {
		return err
	}

	events <- Event{
		ID:      StartContainer,
		Service: service,
	}

	err = d.RunHooks(ctx, ref.ID, service.Hooks.Poststart)
	if err != nil {
		return err
	}

	events <- Event{
		ID:      RunHook,
		Service: service,
		Data:    service.Hooks.Poststart,
	}

	err = d.EnsureContainerIsRunning(ctx, ref.ID)
	if err != nil {
		if err2 := d.Docker.ContainerStop(ctx, ref.ID); err2 != nil {
			return fmt.Errorf("%w (cleanup failed: %s)", err, err2)
		}
		return err
	}

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
func (d *Deployment) EnsureContainerIsRunning(ctx context.Context, containerID string) error {
	info, err := d.Docker.ContainerInspect(ctx, containerID)
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

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return ErrContainerTimeout
		case <-time.After(250 * time.Millisecond):
			info, err = d.Docker.ContainerInspect(ctx, containerID)
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

func (d *Deployment) Save() error {
	dir, err := Store.Dir()
	if err != nil {
		return err
	}

	contents, err := json.Marshal(d)
	if err != nil {
		return err
	}

	return os.WriteFile(dir+"/"+d.ID+".json", contents, 0644)
}

// List returns a list of all the manifests in the Store.
func List() ([]*Deployment, error) {
	var manifests []*Deployment

	dir, err := Store.Dir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			return nil, fmt.Errorf("manifest store is corrupted: %s is a directory", entry.Name())
		}

		f, err := Store.Open(entry.Name(), os.O_RDONLY, 0)
		if err != nil {
			return nil, err
		}

		defer f.Close()

		var m Deployment

		err = json.NewDecoder(f).Decode(&m)
		if err != nil {
			return nil, err
		}

		manifests = append(manifests, &m)
	}

	return manifests, nil
}

// Delete removes the manifest from the Store and returns an error if it fails.
// It does not return an error if the manifest does not exist.
func Delete(version string) error {
	dir, err := Store.Dir()
	if err != nil {
		return err
	}

	path := dir + "/" + version + ".json"

	if _, err = os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return err
	}

	return os.Remove(path)
}

// Get returns the manifest for a given version or os.ErrNotExist if it does not exist.
func Get(version string) (*Deployment, error) {
	f, err := Store.Open(version+".json", os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	var m Deployment

	err = json.NewDecoder(f).Decode(&m)
	return &m, err
}