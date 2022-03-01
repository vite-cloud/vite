package runtime

import "github.com/docker/docker/client"

// Client is used as a proxy to interact with the underlying Docker client.
// It is mainly used to log actions performed by the daemon.
type Client struct {
	client *client.Client
}

// NewClient creates a new docker client
func NewClient() (*Client, error) {
	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	return &Client{
		client: docker,
	}, nil
}
