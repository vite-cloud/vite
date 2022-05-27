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

// gatherer is an interface to collect metrics.
// It is used to mock in tests.
type gatherer interface {
	Gather() (*Metrics, error)
}

// Gatherer is the interface to collect metrics, which can be changed to mock in tests.
var Gatherer gatherer = &SystemGatherer{}

// Gather gathers metrics from docker and the system.
func Gather() (*Metrics, error) {
	return Gatherer.Gather()
}

// Metrics holds the metrics about docker and the system.
type Metrics struct {
	SystemMetrics     *SystemMetrics        `json:"system_metrics"`
	ContainersMetrics []*ContainerStatsJSON `json:"containers_metrics"`
}

// SystemMetrics holds the metrics about the system.
type SystemMetrics struct {
	Uptime time.Duration `json:"uptime"`

	MemoryTotal ByteSize `json:"memory_total"`
	MemoryUsed  ByteSize `json:"memory_used"`
	MemoryFree  ByteSize `json:"memory_free"`

	CPUCount int `json:"cpu_count"`

	CPUUser   float64 `json:"cpu_user"`
	CPUSystem float64 `json:"cpu_system"`
	CPUIdle   float64 `json:"cpu_idle"`
}

type ContainerStatsJSON struct {
	Name string `json:"name"`
	ID   string `json:"id"`

	MemoryUsed      uint64  `json:"memory_used"`
	MemoryAvailable uint64  `json:"memory_available"`
	MemoryUsage     float64 `json:"memory_usage"` // percentage

	CPUCount       int     `json:"cpu_count"`
	CPUDelta       uint64  `json:"cpu_delta"`
	CPUSystemDelta uint64  `json:"cpu_system_delta"`
	CPUUsage       float64 `json:"cpu_usage"` // percentage
}

// SystemGatherer collects metrics
type SystemGatherer struct{}

// TestGatherer is a mock gatherer for tests.
type TestGatherer struct {
	Metrics *Metrics
}

// Gather returns the metrics set while testing.
func (t TestGatherer) Gather() (*Metrics, error) {
	return t.Metrics, nil
}

// Gather gathers metrics from docker and the system.
func (s *SystemGatherer) Gather() (*Metrics, error) {
	metrics := &Metrics{
		SystemMetrics: &SystemMetrics{},
	}

	uptime, err := gouptime.Get()
	if err != nil {
		return nil, err
	}

	metrics.SystemMetrics.Uptime = uptime

	err = gatherMemory(metrics)
	if err != nil {
		return nil, err
	}

	err = gatherCPU(metrics)
	if err != nil {
		return nil, err
	}

	err = gatherDockerStats(metrics)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

func gatherMemory(metrics *Metrics) error {
	memory, err := gomemory.Get()
	if err != nil {
		return err
	}

	metrics.SystemMetrics.MemoryTotal = ByteSize(memory.Total)
	metrics.SystemMetrics.MemoryUsed = ByteSize(memory.Used)
	metrics.SystemMetrics.MemoryFree = ByteSize(memory.Free)

	return nil
}

func gatherCPU(metrics *Metrics) error {
	before, err := gocpu.Get()
	if err != nil {
		return err
	}

	time.Sleep(time.Second)

	after, err := gocpu.Get()
	if err != nil {
		return err
	}

	metrics.SystemMetrics.CPUCount = after.CPUCount
	metrics.SystemMetrics.CPUUser = float64(after.User-before.User) / float64(after.CPUCount)
	metrics.SystemMetrics.CPUSystem = float64(after.System-before.System) / float64(after.CPUCount)
	metrics.SystemMetrics.CPUIdle = float64(after.Idle-before.Idle) / float64(after.CPUCount)

	return nil
}

func gatherDockerStats(metrics *Metrics) error {
	client, err := runtime.NewClient()
	if err != nil {
		return err
	}

	stats, err := client.Stats(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return err
	}

	for _, stat := range stats {
		metrics.ContainersMetrics = append(metrics.ContainersMetrics, &ContainerStatsJSON{
			Name:            stat.Name,
			ID:              stat.ID,
			MemoryUsed:      stat.MemoryUsed,
			MemoryAvailable: stat.MemoryAvailable,
			MemoryUsage:     stat.MemoryUsage,
			CPUCount:        stat.CPUCount,
			CPUDelta:        stat.CPUDelta,
			CPUSystemDelta:  stat.CPUSystemDelta,
			CPUUsage:        stat.CPUUsage,
		})
	}

	return nil
}
