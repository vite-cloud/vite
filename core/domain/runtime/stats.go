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
	Name string
	ID   string

	MemoryUsed      uint64
	MemoryAvailable uint64
	MemoryUsage     float64 // percentage

	CPUCount       int
	CPUDelta       uint64
	CPUSystemDelta uint64
	CPUUsage       float64 // percentage
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

			var decoded *types.StatsJSON
			if err = json.NewDecoder(stat.Body).Decode(&decoded); err != nil {
				errs <- err
				return
			}

			var memoryCache uint64

			if _, ok := decoded.MemoryStats.Stats["cache"]; ok {
				memoryCache = decoded.MemoryStats.Stats["cache"]
			}

			memoryUsed := decoded.MemoryStats.Usage - memoryCache
			memoryAvailable := decoded.MemoryStats.Limit

			memoryUsage := float64(memoryUsed) / float64(memoryAvailable) * 100.0

			cpuDelta := decoded.CPUStats.CPUUsage.TotalUsage - decoded.PreCPUStats.CPUUsage.TotalUsage
			cpuSystemDelta := decoded.CPUStats.SystemUsage - decoded.PreCPUStats.SystemUsage
			cpuCount := len(decoded.CPUStats.CPUUsage.PercpuUsage)

			cpuUsage := float64(cpuDelta) / float64(cpuSystemDelta) * 100.0 * float64(cpuCount)

			metrics = append(metrics, &ContainerStats{
				Name:            container.Names[0][1:],
				ID:              container.ID,
				MemoryUsed:      memoryUsed,
				MemoryAvailable: memoryAvailable,
				MemoryUsage:     memoryUsage,

				CPUCount:       cpuCount,
				CPUDelta:       cpuDelta,
				CPUSystemDelta: cpuSystemDelta,
				CPUUsage:       cpuUsage,
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
