package container

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func (c *Container) configStoreDir() string {
	return ensureDirExists(c.Home() + "/config-store")
}

func (c *Container) configFile() string {
	return ensureDirExists(c.Home() + "/config.json")
}

func (c *Container) manifestsDir() string {
	return ensureDirExists(c.Home() + "/manifests")
}

func (c *Container) certsDir() string {
	return ensureDirExists(c.Home() + "/certs")
}

func (c *Container) logFile() string {
	return ensureDirExists(c.Home() + "/logs/internal.log")
}

func (c *Container) proxyLogFile() string {
	return ensureDirExists(c.Home() + "/logs/proxy.log")
}

func (c *Container) cloudCredentialsFile() string {
	return c.Home() + "/.creds"
}

func (c *Container) subnetRegistryPath() string {
	return ensureDirExists(c.Home() + "/subnets.list")
}

// ensureDirExists creates all the directories in a given path if they don't exist.
func ensureDirExists(path string) string {
	// if the path contains a filename, create all its parent directories
	filename := filepath.Base(path)
	var isFilename bool
	if !strings.HasPrefix(filename, ".") && strings.Contains(filename, ".") {
		path = filepath.Dir(path)
		isFilename = true
	}

	_, err := os.Stat(path)

	switch {
	case errors.Is(err, fs.ErrNotExist):
		_ = os.MkdirAll(path, 0755)
	case err != nil:
		panic(err)
	}

	if isFilename {
		path += "/" + filename
	}

	return strings.TrimRight(path, "/") // remove trailing slash
}
