package docker

import (
	"github.com/docker/docker/client"
	"github.com/redwebcreation/nest/loggy"
	"log"
)

type Client struct {
	// Client is the underlying docker client
	// todo: make this private
	Client    *client.Client
	Logger    *log.Logger
	Subnetter *Subnetter
}

func (c Client) Log(level loggy.Level, message string, fields loggy.Fields) {
	c.Logger.Print(loggy.NewEvent(level, message, fields))
}
