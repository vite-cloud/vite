package container

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/c-robinson/iplib"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/redwebcreation/nest/cloud"
	"github.com/redwebcreation/nest/docker"
	"github.com/redwebcreation/nest/loggy"
	"github.com/redwebcreation/nest/proxy"
	"github.com/redwebcreation/nest/proxy/plane"
	"github.com/redwebcreation/nest/service"
	"golang.org/x/crypto/acme/autocert"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"
)

// Container is a struct that holds the state of the application
// and creates the necessary components to run the application
type Container struct {
	// home is the path to the logger config for nest.
	//
	// It resolves to the following (in order):
	// - WithConfigHome option
	// - the --config/-c flag
	// - $NEST_HOME
	// - ~/.nest
	home string
	// config contains the path to nest's config file.
	// it is resolved once and cached.
	config *service.Locator
	// servicesConfig contains the resolved server config from the config.
	// it is resolved once and cached.
	servicesConfig *service.Config
	// out is a minimal interface to write to stdout.
	out FileWriter
	// in is a minimal interface to read from stdin.
	in FileReader
	// err is a minimal interface to write to stderr.
	err io.Writer
	// logger is nest's internal logger.
	// it is used to log any action that changes any kind of state.
	logger *log.Logger
	// proxyLogger is solely used to log proxy events such as a request coming in, an error in the proxy, etc.
	proxyLogger *log.Logger

	manifestManager *service.ManifestManager

	// dockerClient is our abstraction on top of the docker client.
	dockerClient *docker.Client
	// docker is the true docker client.
	docker *client.Client
}

func New(opts ...Option) (*Container, error) {
	ct := &Container{}
	defaultOptions := []Option{
		WithDefaultConfigHome(),
		WithDefaultInternalLogger(),
		WithDefaultProxyLogger(),
		WithDefaultDockerClient(),
	}

	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(ct); err != nil {
			return nil, err
		}
	}

	return ct, nil
}

// Config returns the cached nest config or loads it if it hasn't been loaded yet.
func (c *Container) Config() (*service.Locator, error) {
	if c.config == nil {
		contents, err := os.ReadFile(c.configFile())
		if err != nil {
			return nil, fmt.Errorf("run `nest setup` to setup nest")
		}

		cf := &service.Locator{
			StoreDir: c.configStoreDir(),
			Path:     c.configFile(),
			Logger:   c.Logger(),
			Git: &service.Git{
				Logger: c.logger,
			},
		}
		err = json.Unmarshal(contents, cf)
		if err != nil {
			return nil, err
		}

		c.config = cf
	}

	return c.config, nil
}

// ServicesConfig returns the cached services config or loads it if it hasn't been loaded yet.
func (c *Container) ServicesConfig() (*service.Config, error) {
	if c.servicesConfig == nil {
		cfg, err := c.Config()
		if err != nil {
			return nil, err
		}
		configFile := "nest.yaml"
		if cfg.Git.Exists(c.configStoreDir(), "nest.yml", cfg.Commit) {
			configFile = "nest.yml"
		}

		contents, err := cfg.Read(configFile)
		if err != nil {
			return nil, err
		}

		// todo: run medic
		servicesConfig := &service.Config{}
		err = yaml.Unmarshal(contents, servicesConfig)
		if err != nil {
			return nil, err
		}

		c.servicesConfig = servicesConfig
	}

	return c.servicesConfig, nil
}

func (c *Container) Out() FileWriter {
	if c.out == nil {
		c.out = os.Stdout
	}

	return c.out
}

func (c *Container) In() FileReader {
	if c.in == nil {
		c.in = os.Stdin
	}

	return c.in
}

func (c *Container) Err() io.Writer {
	if c.err == nil {
		c.err = os.Stderr
	}

	return c.err
}

func (c *Container) Home() string {
	return c.home
}

func (c *Container) ProxyLogger() *log.Logger {
	return c.proxyLogger
}

func (c *Container) Logger() *log.Logger {
	return c.logger
}

func (c *Container) ManifestManager() *service.ManifestManager {
	if c.manifestManager == nil {
		c.manifestManager = &service.ManifestManager{
			Path: c.manifestsDir(),
		}
	}

	return c.manifestManager
}

func (c *Container) CloudCredentials() (id, token string, err error) {
	bytes, err := os.ReadFile(c.cloudCredentialsFile())
	if err != nil {
		return "", "", err
	}

	credentials := string(bytes)

	if err = cloud.ValidateDsn(credentials); err != nil {
		return "", "", err
	}

	return cloud.ParseDsn(credentials)
}

func (c *Container) CloudClient() (*cloud.Client, error) {
	id, token, err := c.CloudCredentials()
	if err != nil {
		return nil, err
	}

	return cloud.NewClient(id, token), nil
}

func (c *Container) SetCloudCredentials(id string, token string) error {
	return ioutil.WriteFile(c.cloudCredentialsFile(), []byte(cloud.FormatDsn(id, token)), 0600)
}

func (c *Container) CertificateStore() autocert.DirCache {
	return autocert.DirCache(c.certsDir())
}

func (c *Container) NewConfig(provider, repository, branch string) *service.Locator {
	return &service.Locator{
		Provider:   provider,
		Repository: repository,
		Branch:     branch,
		Path:       c.configFile(),
		StoreDir:   c.configStoreDir(),
		Logger:     c.Logger(),
		Git: &service.Git{
			Logger: c.Logger(),
		},
	}
}

func (c *Container) DockerClient() *docker.Client {
	if c.dockerClient == nil {
		c.dockerClient = &docker.Client{
			Client: c.docker,
			Logger: c.Logger(),
			Subnetter: &docker.Subnetter{
				Lock:         &sync.Mutex{},
				RegistryPath: c.subnetRegistryPath(),
				Subnets: []iplib.Net4{
					iplib.NewNet4(net.IPv4(10, 0, 0, 0), 8),
				},
			},
		}
	}

	return c.dockerClient
}

func (c *Container) NewProxy(manifest *service.Manifest) (*proxy.Proxy, error) {
	servicesConfig, err := c.ServicesConfig()
	if err != nil {
		return nil, err
	}

	controlPlane, err := c.ControlPlane()
	if err != nil {
		return nil, err
	}

	p := &proxy.Proxy{
		ControlPlane:     controlPlane,
		Logger:           c.ProxyLogger(),
		CertificateStore: c.CertificateStore(),
		Config:           servicesConfig,
		HostToIP:         make(map[string]string),
	}

	for _, service := range servicesConfig.Services {
		for _, host := range service.Hosts {
			id := manifest.Containers[service.Name]

			ip, err := c.DockerClient().GetContainerIP(id)
			if err != nil {
				c.ProxyLogger().Print(loggy.NewEvent(loggy.ErrorLevel, "failed to get container ip", loggy.Fields{
					"error": err,
				}))

				continue
			}

			p.HostToIP[host] = ip
		}
	}

	p.CertificateManager = &autocert.Manager{
		Prompt: autocert.AcceptTOS,
		HostPolicy: func(ctx context.Context, host string) error {
			if _, ok := p.HostToIP[host]; ok {
				return nil
			}

			return fmt.Errorf("acme/autocert: host %s not configured", host)
		},
		Cache: c.CertificateStore(),
	}

	return p, nil
}

func (c *Container) ControlPlane() (*gin.Engine, error) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(func(ctx *gin.Context) {
		ctx.Set("nest", c)

		ctx.Next()
	})

	cfg, err := c.Config()
	if err != nil {
		return nil, err
	}

	servicesConfig, err := c.ServicesConfig()
	if err != nil {
		return nil, err
	}

	return plane.ControlPlane{
		ManifestManager: c.ManifestManager(),
		ServicesConfig:  servicesConfig,
		Config:          cfg,
		Docker:          c.DockerClient(),
	}.From(router), nil

}
