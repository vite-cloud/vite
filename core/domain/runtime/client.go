package runtime

import (
	"github.com/docker/docker/client"
)

// Client is used as a proxy to interact with the underlying Docker client.
// It is mainly used to log actions performed by the daemon.
type Client struct {
	client *client.Client
}

type Opt func(*Client)

func WithDockerClient(dockerClient *client.Client) Opt {
	return func(c *Client) {
		c.client = dockerClient
	}
}

// NewClient creates a new docker client
func NewClient(opts ...Opt) (*Client, error) {
	clientInstance := &Client{}

	for _, opt := range opts {
		opt(clientInstance)
	}

	if clientInstance.client == nil {
		docker, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			return nil, err
		}

		clientInstance.client = docker
	}

	return clientInstance, nil
}
