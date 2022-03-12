package metrics

import (
	gocpu "github.com/mackerelio/go-osstat/cpu"
	gomemory "github.com/mackerelio/go-osstat/memory"
	"gotest.tools/v3/assert"
	"testing"
	"time"
)

func TestGather(t *testing.T) {
	t.Parallel()

	metrics := &Metrics{
		SystemMetrics: &SystemMetrics{
			CPUCount: 123456789,
		},
	}
	testGatherer := &TestGatherer{metrics}
	Gatherer = testGatherer

	got, err := Gather()
	assert.NilError(t, err)
	assert.Equal(t, got, metrics)

	Gatherer = &SystemGatherer{}
}

func TestGatherMemory(t *testing.T) {
	t.Parallel()

	memory, err := gomemory.Get()
	assert.NilError(t, err)

	metrics := &Metrics{
		SystemMetrics: &SystemMetrics{},
	}
	err = gatherMemory(metrics)
	assert.NilError(t, err)

	assert.Assert(t, metrics.SystemMetrics.MemoryTotal != 0)
	assert.Assert(t, metrics.SystemMetrics.MemoryUsed != 0)
	assert.Assert(t, metrics.SystemMetrics.MemoryFree != 0)

	mb20 := MB * 20

	assert.Assert(t, metrics.SystemMetrics.MemoryTotal == ByteSize(memory.Total))
	assert.Assert(t, (metrics.SystemMetrics.MemoryUsed-metrics.SystemMetrics.MemoryUsed) <= mb20, "diff used memory (> 20MB): %f", metrics.SystemMetrics.MemoryUsed-metrics.SystemMetrics.MemoryUsed)
	assert.Assert(t, (metrics.SystemMetrics.MemoryFree-metrics.SystemMetrics.MemoryFree) <= mb20, "diff free memory (> 20MB): %f", metrics.SystemMetrics.MemoryFree-metrics.SystemMetrics.MemoryFree)

	assert.Assert(t, metrics.SystemMetrics.MemoryUsed+metrics.SystemMetrics.MemoryFree <= metrics.SystemMetrics.MemoryTotal, "used+free memory is larger than total memory")
}

func TestGatherCPU(t *testing.T) {
	t.Parallel()

	before, err := gocpu.Get()
	assert.NilError(t, err)

	time.Sleep(time.Second)

	after, err := gocpu.Get()
	assert.NilError(t, err)

	metrics := &Metrics{
		SystemMetrics: &SystemMetrics{},
	}

	err = gatherCPU(metrics)
	assert.NilError(t, err)

	assert.Equal(t, metrics.SystemMetrics.CPUCount, after.CPUCount)

	user := float64(after.User-before.User) / float64(after.CPUCount)
	system := float64(after.System-before.System) / float64(after.CPUCount)
	idle := float64(after.Idle-before.Idle) / float64(after.CPUCount)

	assert.Assert(t, metrics.SystemMetrics.CPUSystem < 100)
	assert.Assert(t, metrics.SystemMetrics.CPUUser < 100)
	assert.Assert(t, metrics.SystemMetrics.CPUIdle < 100)

	diff := 15.0

	assert.Assert(t, user-metrics.SystemMetrics.CPUUser <= diff, "user diff: %f", user-metrics.SystemMetrics.CPUUser)
	assert.Assert(t, system-metrics.SystemMetrics.CPUUser >= -diff, "user diff: %f", system-metrics.SystemMetrics.CPUUser)

	assert.Assert(t, system-metrics.SystemMetrics.CPUSystem <= diff, "system diff: %f", system-metrics.SystemMetrics.CPUSystem)
	assert.Assert(t, idle-metrics.SystemMetrics.CPUIdle >= -diff, "system diff: %f", system-metrics.SystemMetrics.CPUSystem)

	assert.Assert(t, idle-metrics.SystemMetrics.CPUIdle <= diff, "idle diff: %f", idle-metrics.SystemMetrics.CPUIdle)
	assert.Assert(t, user-metrics.SystemMetrics.CPUUser >= -diff, "idle diff: %f", idle-metrics.SystemMetrics.CPUIdle)

	assert.Assert(t, 100-(user+system+idle) <= diff, "total diff: %f", 100-(user+system+idle))
	assert.Assert(t, 100-(user+system+idle) >= -diff, "total diff: %f", 100-(user+system+idle))
}
