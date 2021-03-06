package runtime

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/vite-cloud/go-zoup"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/c-robinson/iplib"
	"github.com/vite-cloud/vite/core/domain/datadir"
	"github.com/vite-cloud/vite/core/domain/log"
)

// subnetManager handles all subnet related operations
type subnetManager struct {
	// mu ensures that we
	mu     *sync.Mutex
	used   *os.File
	blocks []iplib.Net4
}

// related errors
var (
	ErrNoAvailableSubnet      = errors.New("no available subnet")
	ErrSubnetAlreadyAllocated = errors.New("subnet already allocated")
)

// Store is the unique name of the subnet store
// It will be located under the data directory and must absolutely be unique
// Changing it is a major breaking change.
const Store = datadir.Store("subnets")

// SubnetDataFile is the name of the file that stores created subnets.
const SubnetDataFile = "subnet.dat"

// DefaultSubnetBlocks is the list of private ipv4 blocks.
// It respects https://datatracker.ietf.org/doc/html/rfc1918.
// It contains the following blocks:
// - 10.0.0.0/8
// - 172.16.0.0/12
// - 192.168.0.0/16
var DefaultSubnetBlocks = []iplib.Net4{
	iplib.NewNet4(net.IPv4(10, 0, 0, 0), 8),
	iplib.NewNet4(net.IPv4(172, 16, 0, 0), 12),
	iplib.NewNet4(net.IPv4(192, 168, 0, 0), 16),
}

// NewSubnetManager creates a new subnet manager.
func NewSubnetManager() (*subnetManager, error) {
	file, err := Store.Open(SubnetDataFile, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}

	return &subnetManager{
		mu:     &sync.Mutex{},
		used:   file,
		blocks: DefaultSubnetBlocks,
	}, nil
}

// WithBlocks sets the list of subnet blocks to use.
func (sm *subnetManager) WithBlocks(blocks []iplib.Net4) *subnetManager {
	sm.blocks = blocks
	return sm
}

// Next returns the next available subnet from any of the blocks.
func (sm *subnetManager) Next() (*iplib.Net4, error) {
	next := make(chan iplib.Net4)
	failed := 0

	for _, network := range sm.blocks {
		go func(network iplib.Net4) {
			subnets, _ := network.Subnet(24)

			for _, subnet := range subnets {
				ok, _ := sm.IsFree(subnet.String())

				if ok {
					next <- subnet
					return
				}
			}

			sm.mu.Lock()
			failed++
			sm.mu.Unlock()
		}(network)
	}

	for {
		select {
		case subnet := <-next:
			if err := sm.Allocate(subnet.String()); err != nil {
				return nil, err
			}
			return &subnet, nil
		case <-time.After(time.Millisecond * 50):
			sm.mu.Lock()
			if failed == len(sm.blocks) {
				sm.mu.Unlock()
				return nil, ErrNoAvailableSubnet
			}
			sm.mu.Unlock()
		}
	}
}

// Allocate allocates a subnet.
func (sm *subnetManager) Allocate(subnet string) error {
	ok, err := sm.IsFree(subnet)
	if err != nil {
		return err
	}

	if !ok {
		return ErrSubnetAlreadyAllocated
	}

	sm.mu.Lock()
	_, err = sm.used.WriteString(fmt.Sprintf("%s\n", subnet))
	sm.mu.Unlock()

	if err != nil {
		return err
	}

	log.Log(zoup.DebugLevel, "subnet allocated", zoup.Fields{
		"subnet": subnet,
	})

	return nil
}

// IsFree checks if a subnet is free for allocation.
func (sm *subnetManager) IsFree(subnet string) (bool, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	currentOffset, err := sm.used.Seek(0, io.SeekCurrent)
	if err != nil {
		return false, err
	}

	_, err = sm.used.Seek(0, io.SeekStart)
	if err != nil {
		return false, err
	}

	defer func(used *os.File, offset int64) {
		_, err = used.Seek(offset, io.SeekStart)
		if err != nil {
			panic(err)
		}
	}(sm.used, currentOffset)

	scanner := bufio.NewScanner(sm.used)

	for scanner.Scan() {
		cmp := scanner.Text()
		if cmp == subnet {
			return false, nil
		}
	}

	return true, nil
}
