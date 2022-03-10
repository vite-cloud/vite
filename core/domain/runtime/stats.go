package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"sync"
)

// ContainerStats holds information about the memory and cpu usage of a container
type ContainerStats struct {
	Name            string
	ID              string
	MemoryUsed      uint64
	MemoryAvailable uint64
	CPUCount        uint64
	CPUDelta        uint64
	CPUSystemDelta  uint64
}

// Stats returns the stats of many containers
func (c *Client) Stats(ctx context.Context, opts types.ContainerListOptions) ([]*ContainerStats, error) {
	var metrics []*ContainerStats

	containers, err := c.client.ContainerList(ctx, opts)
	if err != nil {
		return nil, err
	}

	wg := &sync.WaitGroup{}
	errs := make(chan error)

	wg.Add(len(containers))

	for _, container := range containers {
		container := container

		go func() {
			defer wg.Done()

			stat, err := c.client.ContainerStats(ctx, container.ID, false)
			if err != nil {
				errs <- err
				return
			}

			if stat.OSType != "linux" {
				errs <- fmt.Errorf("unsupported OSType: %s", stat.OSType)
				return
			}

			var decoded *types.StatsJSON
			if err = json.NewDecoder(stat.Body).Decode(&decoded); err != nil {
				errs <- err
				return
			}

			metrics = append(metrics, &ContainerStats{
				Name:            container.Names[0][1:],
				ID:              container.ID,
				MemoryUsed:      decoded.MemoryStats.Usage,
				MemoryAvailable: decoded.MemoryStats.Limit,
				CPUCount:        decoded.CPUStats.CPUUsage.TotalUsage,
				CPUDelta:        decoded.CPUStats.CPUUsage.UsageInUsermode,
				CPUSystemDelta:  decoded.CPUStats.CPUUsage.UsageInKernelmode,
			})
		}()
	}

	go func() {
		wg.Wait()
		close(errs)
	}()

	var errors error

	for err = range errs {
		errors = fmt.Errorf("%s\n%w", errors, err)
	}

	return metrics, errors
}
