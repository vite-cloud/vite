package metrics

import (
	"context"
	"github.com/docker/docker/api/types"
	gocpu "github.com/mackerelio/go-osstat/cpu"
	gomemory "github.com/mackerelio/go-osstat/memory"
	gouptime "github.com/mackerelio/go-osstat/uptime"
	"github.com/vite-cloud/vite/core/domain/runtime"
	"time"
)

// Metrics holds the metrics about docker and the system.
type Metrics struct {
	SystemMetrics     *SystemMetrics
	ContainersMetrics []*runtime.ContainerStats
}

// SystemMetrics holds the metrics about the system.
type SystemMetrics struct {
	// Uptime is the system uptime in seconds.
	Uptime float64

	MemoryTotal uint64
	MemoryUsed  uint64
	MemoryFree  uint64

	CPUCount int

	CPUUser   float64
	CPUSystem float64
	CPUIdle   float64
}

// Gather gathers metrics from docker and the system.
func Gather() (*Metrics, error) {
	metrics := &Metrics{
		SystemMetrics: &SystemMetrics{},
	}

	memory, err := gomemory.Get()
	if err != nil {
		return nil, err
	}

	metrics.SystemMetrics.MemoryTotal = memory.Total
	metrics.SystemMetrics.MemoryUsed = memory.Used
	metrics.SystemMetrics.MemoryFree = memory.Free

	uptime, err := gouptime.Get()
	if err != nil {
		return nil, err
	}

	metrics.SystemMetrics.Uptime = uptime.Seconds()

	before, err := gocpu.Get()
	if err != nil {
		return nil, err
	}

	time.Sleep(time.Second)

	after, err := gocpu.Get()
	if err != nil {
		return nil, err
	}

	metrics.SystemMetrics.CPUCount = after.CPUCount
	metrics.SystemMetrics.CPUUser = float64(after.User-before.User) / float64(after.CPUCount) * 100
	metrics.SystemMetrics.CPUSystem = float64(after.System-before.System) / float64(after.CPUCount) * 100
	metrics.SystemMetrics.CPUIdle = float64(after.Idle-before.Idle) / float64(after.CPUCount) * 100

	client, err := runtime.NewClient()
	if err != nil {
		return nil, err
	}

	stats, err := client.Stats(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	metrics.ContainersMetrics = stats

	return metrics, nil
}
