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

// Gatherer is the interface to collect metrics, which can be changed to mock in tests.
var Gatherer gatherer = &SystemGatherer{}

// Metrics holds the metrics about docker and the system.
type Metrics struct {
	SystemMetrics     *SystemMetrics
	ContainersMetrics []*runtime.ContainerStats
}

// SystemMetrics holds the metrics about the system.
type SystemMetrics struct {
	Uptime time.Duration

	MemoryTotal ByteSize
	MemoryUsed  ByteSize
	MemoryFree  ByteSize

	CPUCount int

	CPUUser   float64
	CPUSystem float64
	CPUIdle   float64
}

type gatherer interface {
	Gather() (*Metrics, error)
}

type SystemGatherer struct{}

// Gather gathers metrics from docker and the system.
func (s *SystemGatherer) Gather() (*Metrics, error) {
	metrics := &Metrics{
		SystemMetrics: &SystemMetrics{},
	}

	memory, err := gomemory.Get()
	if err != nil {
		return nil, err
	}

	metrics.SystemMetrics.MemoryTotal = ByteSize(memory.Total)
	metrics.SystemMetrics.MemoryUsed = ByteSize(memory.Used)
	metrics.SystemMetrics.MemoryFree = ByteSize(memory.Free)

	uptime, err := gouptime.Get()
	if err != nil {
		return nil, err
	}

	metrics.SystemMetrics.Uptime = uptime

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

type TestGatherer struct {
	Metrics *Metrics
}

func (t TestGatherer) Gather() (*Metrics, error) {
	return t.Metrics, nil
}

func Gather() (*Metrics, error) {
	return Gatherer.Gather()
}
