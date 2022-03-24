package runtime

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/vite-cloud/vite/core/domain/log"
)

type NetworkCreateOptions struct {
	Labels map[string]string
	IPAM   *network.IPAM
}

func (c Client) NetworkCreate(ctx context.Context, name string, opts NetworkCreateOptions) (string, error) {
	res, err := c.client.NetworkCreate(ctx, name, types.NetworkCreate{
		CheckDuplicate: true,
		IPAM:           opts.IPAM,
		Labels:         opts.Labels,
	})
	if err != nil {
		return "", err
	}

	log.Log(log.DebugLevel, "created network", log.Fields{
		"name":   name,
		"id":     res.ID,
		"config": opts.IPAM,
	})

	return res.ID, nil
}

func (c Client) NetworkRemove(ctx context.Context, id string) error {
	err := c.client.NetworkRemove(ctx, id)
	if err != nil {
		return err
	}

	log.Log(log.DebugLevel, "removed network", log.Fields{
		"id": id,
	})

	return nil
}

func (c Client) NetworkConnect(ctx context.Context, networkId, containerID string) error {
	err := c.client.NetworkConnect(ctx, networkId, containerID, nil)
	if err != nil {
		return err
	}

	log.Log(log.DebugLevel, "connected network", log.Fields{
		"network":   networkId,
		"container": containerID,
	})

	return nil
}
