package config

import (
	"github.com/vite-cloud/vite/core/domain/datadir"
	"github.com/vite-cloud/vite/core/domain/locator"
	"gotest.tools/v3/assert"
	"testing"
)

func TestConfig_Reload(t *testing.T) {
	t.Parallel()

	datadir.UseTestHome(t)

	l := locator.Locator{
		Commit:     "8897a7d08a1e791418904afdce369818c19d2c3e",
		Branch:     "main",
		Provider:   locator.Provider("github"),
		Protocol:   "ssh",
		Repository: "felixdorn/config-test",
	}

	err := l.Save()
	assert.NilError(t, err)

	conf, err := Get()
	assert.NilError(t, err)

	assert.Equal(t, len(conf.Services["example"].Hosts), 2)
	assert.Equal(t, conf.Services["example"].Hosts[0], "example.com")
	assert.Equal(t, conf.Services["example"].Hosts[1], "www.example.com")

	l.Branch = "something"
	l.Commit = "4e56859590586c235cb3d09cc04e8f0ad1d9f1d8"

	err = l.Save()
	assert.NilError(t, err)

	err = conf.Reload()
	assert.NilError(t, err)

	assert.Equal(t, len(conf.Services["example"].Hosts), 1)
	assert.Equal(t, conf.Services["example"].Hosts[0], "not.example.com")
}
