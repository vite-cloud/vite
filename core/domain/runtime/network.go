package runtime

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/vite-cloud/vite/core/domain/log"
	"github.com/vite-cloud/vite/core/domain/manifest"
)

// CreateNetworkManifestKey resource tag available in the manifest
const CreateNetworkManifestKey = "CreatedNetwork"

type NetworkCreateOptions struct {
	Driver string
	Labels map[string]string
	IPAM   *network.IPAM
}

func (c Client) NetworkCreate(ctx context.Context, name string, opts NetworkCreateOptions) (string, error) {
	res, err := c.client.NetworkCreate(ctx, name, types.NetworkCreate{
		CheckDuplicate: true,
		Driver:         opts.Driver,
		Labels:         opts.Labels,
		IPAM:           opts.IPAM,
	})
	if err != nil {
		return "", err
	}

	log.Log(log.DebugLevel, "created network", log.Fields{
		"name":   name,
		"id":     res.ID,
		"config": opts.IPAM,
	})

	ctx.Value(manifest.ContextKey).(*manifest.Manifest).Add(CreateNetworkManifestKey, res.ID)

	return res.ID, nil
}
