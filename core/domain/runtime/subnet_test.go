package runtime

import (
	"github.com/c-robinson/iplib"
	"github.com/vite-cloud/vite/core/domain/datadir"
	"github.com/vite-cloud/vite/core/domain/log"
	"gotest.tools/v3/assert"
	"net"
	"os"
	"testing"
)

var logger *log.MemoryWriter

func TestMain(m *testing.M) {
	logger = log.UseTestLogger()
	os.Exit(m.Run())
}

func TestSubnetManager_IsFree(t *testing.T) {
	home, err := os.MkdirTemp("", "subnet-test")
	assert.NilError(t, err)

	datadir.SetHomeDir(home)

	manager, err := NewSubnetManager()
	assert.NilError(t, err)

	subnet := &iplib.Net4{
		IPNet: net.IPNet{
			IP:   net.IP{192, 168, 0, 0},
			Mask: net.IPMask{255, 255, 0, 0},
		},
	}

	ok, err := manager.IsFree(subnet)
	assert.NilError(t, err)
	assert.Assert(t, ok)

	err = manager.Allocate(subnet)
	assert.NilError(t, err)

	ok, err = manager.IsFree(subnet)
	assert.NilError(t, err)
	assert.Assert(t, !ok)

	// regression test for a case where IsFree would return true once.
	ok, err = manager.IsFree(subnet)
	assert.NilError(t, err)
	assert.Assert(t, !ok)
}

func TestSubnetManager_Allocate(t *testing.T) {
	home, err := os.MkdirTemp("", "subnet-test")
	assert.NilError(t, err)

	datadir.SetHomeDir(home)

	manager, err := NewSubnetManager()
	assert.NilError(t, err)

	subnet := &iplib.Net4{
		IPNet: net.IPNet{
			IP:   net.IP{192, 168, 0, 0},
			Mask: net.IPMask{255, 255, 0, 0},
		},
	}

	err = manager.Allocate(subnet)
	assert.NilError(t, err)

	dir, err := Store.Dir()
	assert.NilError(t, err)
	contents, err := os.ReadFile(dir + "/" + SubnetDataFile)
	assert.NilError(t, err)
	assert.Assert(t, len(contents) > 0)
	assert.Equal(t, string(contents), subnet.String()+"\n")
}

func TestSubnetManager_Allocate2(t *testing.T) {
	home, err := os.MkdirTemp("", "subnet-test")
	assert.NilError(t, err)

	datadir.SetHomeDir(home)

	manager, err := NewSubnetManager()
	assert.NilError(t, err)

	subnet := &iplib.Net4{
		IPNet: net.IPNet{
			IP:   net.IP{192, 168, 0, 0},
			Mask: net.IPMask{255, 255, 0, 0},
		},
	}

	err = manager.Allocate(subnet)
	assert.NilError(t, err)

	err = manager.Allocate(subnet)
	assert.ErrorIs(t, err, ErrSubnetAlreadyAllocated)
}

func TestSubnetManager_Next(t *testing.T) {
	home, err := os.MkdirTemp("", "subnet-test")
	assert.NilError(t, err)

	datadir.SetHomeDir(home)

	manager, err := NewSubnetManager()
	assert.NilError(t, err)

	// we're specifically giving only one subnet to the manager,
	// so we know for sure what the next subnet will be
	manager.WithBlocks([]iplib.Net4{
		iplib.NewNet4(net.IPv4(192, 168, 0, 0), 16),
	})

	subnet, err := manager.Next()
	assert.NilError(t, err)

	assert.Equal(t, subnet.String(), "192.168.0.0/24")

	subnet, err = manager.Next()
	assert.NilError(t, err)

	assert.Equal(t, subnet.String(), "192.168.1.0/24")
}

func TestSubnetManager_Next2(t *testing.T) {
	home, err := os.MkdirTemp("", "subnet-test")
	assert.NilError(t, err)

	datadir.SetHomeDir(home)

	manager, err := NewSubnetManager()
	assert.NilError(t, err)

	for i := 0; i < 257; i++ {
		_, err = manager.Next()
		assert.NilError(t, err)
	}
}

func TestSubnetManager_Allocate3(t *testing.T) {
	home, err := os.MkdirTemp("", "subnet-test")
	assert.NilError(t, err)

	datadir.SetHomeDir(home)

	manager, err := NewSubnetManager()
	assert.NilError(t, err)

	subnet := iplib.NewNet4(net.IPv4(192, 168, 0, 0), 24)
	err = manager.Allocate(&subnet)
	assert.NilError(t, err)

	assert.Assert(t, logger.Len() > 0)
	assert.Assert(t, logger.Last().Level == log.DebugLevel)
	assert.Assert(t, logger.Last().Message == "subnet allocated")
	assert.Assert(t, len(logger.Last().Fields) == 1)
	assert.Assert(t, logger.Last().Fields["subnet"] == subnet.String())
}
