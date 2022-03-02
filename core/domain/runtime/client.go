package runtime

import "github.com/docker/docker/client"

// Client is used as a proxy to interact with the underlying Docker client.
// It is mainly used to log actions performed by the daemon.
type Client struct {
	client *client.Client
}

var clientInstance *Client

// NewClient creates a new docker client
func NewClient() (*Client, error) {
	if clientInstance != nil {
		return clientInstance, nil
	}

	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	clientInstance = &Client{
		client: docker,
	}

	return clientInstance, nil
}
