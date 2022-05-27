package runtime

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/vite-cloud/go-zoup"
	"github.com/vite-cloud/vite/core/domain/log"
)

func (c *Client) RegistryLogin(ctx context.Context, auth types.AuthConfig) error {
	_, err := c.client.RegistryLogin(ctx, auth)

	c.client.HTTPClient()
	log.Log(zoup.DebugLevel, "login to registry", zoup.Fields{
		"host": auth.ServerAddress,
	})

	return err
}
