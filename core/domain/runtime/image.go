package runtime

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/vite-cloud/go-zoup"
	"github.com/vite-cloud/vite/core/domain/log"
	"io"
)

// ImagePullOptions is the set of options that can be used when pulling an image.
type ImagePullOptions struct {
	// Auth is the authentication settings for pulling the image on a custom registry.
	Auth *types.AuthConfig
	// Listener is an optional progress listener.
	// It gets called every time, the daemon sends a progress event.
	Listener func(status string)
}

func marshalAuth(auth *types.AuthConfig) (string, error) {
	if auth == nil {
		return "", nil
	}
	authJSON, err := json.Marshal(auth)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(authJSON), nil
}

// ImagePull pulls an image from a remote registry.
func (c Client) ImagePull(ctx context.Context, image string, options ImagePullOptions) error {
	opts := types.ImagePullOptions{}

	auth, err := marshalAuth(options.Auth)
	if err != nil {
		return err
	}

	opts.RegistryAuth = auth

	events, err := c.client.ImagePull(ctx, image, opts)
	if err != nil {
		return err
	}

	log.Log(zoup.DebugLevel, "pulling docker image", zoup.Fields{
		"image":     image,
		"with_auth": options.Auth != nil,
	})

	decoder := json.NewDecoder(events)

	var event *struct{ Status string }

	for {
		if err = decoder.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		if options.Listener != nil {
			options.Listener(event.Status)
		}
	}

	return nil
}
