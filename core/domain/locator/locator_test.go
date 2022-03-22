package locator

import (
	"os"
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

func TestLocator_Read2(t *testing.T) {
	t.Parallel()

	locator := Locator{}

	_, err := locator.Read("vite.yaml")
	assert.ErrorIs(t, err, ErrInvalidCommit)
}

func TestLocator_Read3(t *testing.T) {
	datadir.UseTestHome(t)

	locator := Locator{
		Commit:     "8897a7d08a1e791418904afdce369818c19d2c3e",
		Provider:   Provider("github"),
		Protocol:   "ssh",
		Repository: "felixdorn/config-test",
	}

	_, err := locator.Read("vite.yaml")
	assert.Error(t, err, "could not clone repository git@github.com:felixdorn/config-test.git: no branch specified (run `vite setup` again)")
}

func TestLocator_Read4(t *testing.T) {
	datadir.UseTestHome(t)

	locator := Locator{
		Commit:     "8897a7d08a1e791418904afdce369818c19d2c3e",
		Provider:   Provider("github"),
		Branch:     "main",
		Protocol:   "ssh",
		Repository: "felixdorn/config-test",
	}

	contents, err := locator.Read("vite.yaml")
	assert.NilError(t, err)
	assert.Equal(t, string(contents), "services:\n  example:\n    # will fail as it requires some variables on startup\n    image: nginx:1.21.5\n    hosts:\n      - ~example.com\ncontrol_plane:\n  host: vite.example.com\n")
}

func TestLocator_Read5(t *testing.T) {
	datadir.UseTestHome(t)

	locator := Locator{
		Commit:     "8897a7d08a1e791418904afdce369818c19d2c3e",
		Provider:   Provider("github"),
		Branch:     "main",
		Protocol:   "ssh",
		Repository: "felixdorn/config-test",
	}

	_, err := locator.Read("does-not-exist.yaml")
	assert.ErrorContains(t, err, "path 'does-not-exist.yaml' does not exist")
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

func TestLocator_Git(t *testing.T) {
	datadir.SetHomeDir("/nop")

	locator := Locator{}

	_, err := locator.git()
	assert.ErrorIs(t, err, os.ErrPermission)
}
