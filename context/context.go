package context

import (
	"github.com/redwebcreation/nest/config"
	"github.com/redwebcreation/nest/config/medic"
	"github.com/redwebcreation/nest/deploy"
	"io"
	"log"
	"os"
)

// Context is a struct that holds the context of the application
type Context struct {
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
	config *config.Config
	// servicesConfig contains the resolved server config from the config.
	// it is resolved once and cached.
	servicesConfig *config.ServicesConfig
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

	manifestManager *deploy.Manager
}

// Config returns the cached nest config or loads it if it hasn't been loaded yet.
func (c *Context) Config() (*config.Config, error) {
	if c.config == nil {
		cf, err := config.NewConfig(c.ConfigFile(), c.ConfigStoreDir(), c.Logger())
		if err != nil {
			return nil, err
		}

		c.config = cf
	}

	return c.config, nil
}

// ServicesConfig returns the cached services config or loads it if it hasn't been loaded yet.
func (c *Context) ServicesConfig() (*config.ServicesConfig, error) {
	nc, err := c.Config()
	if err != nil {
		return nil, err
	}

	if c.servicesConfig == nil {
		servicesConfig, err := nc.ServerConfig()
		if err != nil {
			return nil, err
		}

		c.servicesConfig = servicesConfig
	}

	err = medic.DiagnoseConfig(c.servicesConfig).MustPass()
	if err != nil {
		return nil, err
	}

	return c.servicesConfig, nil
}

func (c *Context) Out() FileWriter {
	if c.out == nil {
		c.out = os.Stdout
	}

	return c.out
}

func (c Context) In() FileReader {
	if c.in == nil {
		c.in = os.Stdin
	}

	return c.in
}

func (c Context) Err() io.Writer {
	if c.err == nil {
		c.err = os.Stderr
	}

	return c.err
}

func (c Context) Home() string {
	return c.home
}

func (c Context) ProxyLogger() *log.Logger {
	return c.proxyLogger
}

func (c Context) Logger() *log.Logger {
	return c.logger
}

func (c Context) ManifestManager() *deploy.Manager {
	if c.manifestManager == nil {
		c.manifestManager = &deploy.Manager{
			Path: c.ManifestsDir(),
		}
	}

	return c.manifestManager
}

func New(opts ...Option) (*Context, error) {
	ctx := &Context{}
	defaultOptions := []Option{
		WithDefaultConfigHome(),
		WithDefaultInternalLogger(),
		WithDefaultProxyLogger(),
	}

	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(ctx); err != nil {
			return nil, err
		}
	}

	return ctx, nil
}
