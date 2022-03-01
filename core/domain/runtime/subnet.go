package runtime

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/c-robinson/iplib"
	"github.com/vite-cloud/vite/core/domain/datadir"
	"github.com/vite-cloud/vite/core/domain/log"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

// subnetManager handles all subnet related operations
type subnetManager struct {
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

func (sm *subnetManager) WithBlocks(blocks []iplib.Net4) *subnetManager {
	sm.blocks = blocks
	return sm
}

// Next returns the next available subnet from any of the blocks.
func (sm *subnetManager) Next() (*iplib.Net4, error) {
	next := make(chan *iplib.Net4)
	failed := 0

	for _, network := range sm.blocks {
		network := network
		go func() {
			subnets, _ := network.Subnet(24)

			for _, subnet := range subnets {
				ok, _ := sm.IsFree(&subnet)

				if ok {
					next <- &subnet
					return
				}
			}

			failed++
		}()
	}

	for {
		select {
		case subnet := <-next:
			if err := sm.Allocate(subnet); err != nil {
				return nil, err
			}
			return subnet, nil
		case <-time.After(time.Millisecond * 50):
			if failed == len(sm.blocks) {
				return nil, ErrNoAvailableSubnet
			}
		}
	}
}

// Allocate allocates a subnet.
func (sm *subnetManager) Allocate(subnet *iplib.Net4) error {
	ok, err := sm.IsFree(subnet)
	if err != nil {
		return err
	}

	if !ok {
		return ErrSubnetAlreadyAllocated
	}

	sm.mu.Lock()
	_, err = sm.used.WriteString(fmt.Sprintf("%s\n", subnet.String()))
	sm.mu.Unlock()

	if err != nil {
		return err
	}

	log.Log(log.DebugLevel, "subnet allocated", log.Fields{
		"subnet": subnet.String(),
	})

	return nil
}

// IsFree checks if a subnet is free for allocation.
func (sm *subnetManager) IsFree(subnet *iplib.Net4) (bool, error) {
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
		if cmp == subnet.String() {
			return false, nil
		}
	}

	return true, nil
}
