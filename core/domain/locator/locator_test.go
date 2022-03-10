package locator

import (
	"path/filepath"
	"testing"

	"github.com/vite-cloud/vite/core/domain/datadir"
	"gotest.tools/v3/assert"
)

func TestLocator_Read(t *testing.T) {
	locator := Locator{
		Branch:     "main",
		Repository: "foo/bar",
		Provider:   Provider("github"),
	}

	datadir.UseTestHome(t)

	dir, err := ConfigStore.Dir()
	assert.NilError(t, err)

	builder := newLocalRepo(t, filepath.Join(dir, "main-foo-bar"))
	commit := builder.
		WriteFile("vite.yaml", []byte("services:\n"), 0600).
		Commit()

	locator.Commit = commit

	contents, err := locator.Read("vite.yaml")
	assert.NilError(t, err)

	assert.Equal(t, string(contents), "services:\n")
}

func TestLoadFromStore(t *testing.T) {
	datadir.UseTestHome(t)

	locator := Locator{
		Branch:     "main",
		Repository: "foo/bar",
		Provider:   Provider("github"),
		Commit:     "ffffffffffffffffffffffffffffffffffffffff",
		Path:       "/sub/path",
		Protocol:   "https",
	}
	err := locator.Save()
	assert.NilError(t, err)

	l, err := LoadFromStore()
	assert.NilError(t, err)

	assert.Equal(t, l.Branch, "main")
	assert.Equal(t, l.Repository, "foo/bar")
	assert.Equal(t, l.Commit, "ffffffffffffffffffffffffffffffffffffffff")
	assert.Equal(t, l.Path, "/sub/path")
	assert.Equal(t, l.Protocol, "https")
	assert.Equal(t, l.Provider.Name(), "github")
}
