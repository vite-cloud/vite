package metrics

import (
	gocpu "github.com/mackerelio/go-osstat/cpu"
	gomemory "github.com/mackerelio/go-osstat/memory"
	"gotest.tools/v3/assert"
	"testing"
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

	assert.Assert(t, metrics.SystemMetrics.MemoryTotal == ByteSize(memory.Total))
	assert.Assert(t, metrics.SystemMetrics.MemoryUsed+metrics.SystemMetrics.MemoryFree <= metrics.SystemMetrics.MemoryTotal, "used+free memory is larger than total memory")
}

func TestGatherCPU(t *testing.T) {
	t.Parallel()

	cpu, err := gocpu.Get()
	assert.NilError(t, err)

	metrics := &Metrics{
		SystemMetrics: &SystemMetrics{},
	}

	err = gatherCPU(metrics)
	assert.NilError(t, err)

	assert.Equal(t, metrics.SystemMetrics.CPUCount, cpu.CPUCount)
	assert.Assert(t, metrics.SystemMetrics.CPUSystem < 100)
	assert.Assert(t, metrics.SystemMetrics.CPUUser < 100)
	assert.Assert(t, metrics.SystemMetrics.CPUIdle < 100)
}
