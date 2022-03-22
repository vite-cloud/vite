package runtime

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/vite-cloud/vite/core/domain/log"
)

func (c *Client) RegistryLogin(ctx context.Context, auth types.AuthConfig) error {
	_, err := c.client.RegistryLogin(ctx, auth)

	log.Log(log.DebugLevel, "login to registry", log.Fields{
		"host": auth.ServerAddress,
	})

	return err
}
