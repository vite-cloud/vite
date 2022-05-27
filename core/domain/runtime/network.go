package runtime

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/vite-cloud/go-zoup"
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

	log.Log(zoup.DebugLevel, "created network", zoup.Fields{
		"name":   name,
		"id":     res.ID,
		"config": opts.IPAM,
	})

	return res.ID, nil
}

func (c Client) NetworkRemove(ctx context.Context, ID string) error {
	err := c.client.NetworkRemove(ctx, ID)
	if err != nil {
		return err
	}

	log.Log(zoup.DebugLevel, "removed network", zoup.Fields{
		"id": ID,
	})

	return nil
}

func (c Client) NetworkConnect(ctx context.Context, networkID, containerID string) error {
	err := c.client.NetworkConnect(ctx, networkID, containerID, nil)
	if err != nil {
		return err
	}

	log.Log(zoup.DebugLevel, "connected network", zoup.Fields{
		"network":   networkID,
		"container": containerID,
	})

	return nil
}
