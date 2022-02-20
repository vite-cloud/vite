package container

import (
	"fmt"
	"github.com/docker/docker/client"
	"github.com/mitchellh/go-homedir"
	"github.com/vite-cloud/vite/loggy"
	"github.com/vite-cloud/vite/service"
	"io"
	"log"
	"os"
	"strings"
)

type Option func(*Container) error

// FileWriter provides a minimal interface for Stdin.
type FileWriter interface {
	io.Writer
	Fd() uintptr
}

// FileReader provides a minimal interface for Stdout.
type FileReader interface {
	io.Reader
	Fd() uintptr
}

func WithConfig(config *service.Locator) Option {
	return func(ct *Container) error {
		ct.config = config

		return nil
	}
}

func WithStdio(stdin FileReader, stdout FileWriter, stderr io.Writer) Option {
	return func(ct *Container) error {
		ct.in = stdin
		ct.out = stdout
		ct.err = stderr

		return nil
	}
}

func WithServicesConfig(servicesConfig *service.Config) Option {
	return func(ct *Container) error {
		ct.servicesConfig = servicesConfig

		return nil
	}
}

func WithDefaultConfigHome() Option {
	return func(context *Container) error {
		for k, arg := range os.Args {
			if arg != "--config" && arg != "-c" {
				continue
			}

			if len(os.Args) <= k+1 {
				fmt.Fprintln(os.Stderr, "--config requires an argument")
				os.Exit(1)
			}

			context.home = strings.TrimRight(os.Args[k+1], "/")
			return nil
		}

		if os.Getenv("VITE_HOME") != "" {
			context.home = strings.TrimRight(os.Getenv("VITE_HOME"), "/")
			return nil
		}

		// otherwise, use the default
		userHome, err := homedir.Dir()
		if err != nil {
			return err
		}

		context.home = userHome + "/.vite"

		return nil
	}
}

func WithConfigHome(home string) Option {
	return func(context *Container) error {
		context.home = home
		return nil
	}
}

func WithDefaultInternalLogger() Option {
	return func(context *Container) error {
		context.logger = log.New(&loggy.FileLogger{
			Path: context.logFile(),
		}, "", 0)
		return nil
	}
}

func WithDefaultProxyLogger() Option {
	return func(context *Container) error {
		context.proxyLogger = log.New(loggy.CompositeLogger{
			Loggers: []io.Writer{
				&loggy.FileLogger{
					Path: context.proxyLogFile(),
				},
				&loggy.FileLogger{
					Writer: os.Stdout,
				},
			},
		}, "", 0)

		return nil
	}
}

func WithDefaultDockerClient() Option {
	return func(context *Container) error {
		client, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			return err
		}

		context.docker = client

		return nil
	}
}
