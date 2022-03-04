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
		Provider:   GitHubProvider{},
	}

	home, err := os.MkdirTemp("", "locator_test")
	assert.NilError(t, err)

	datadir.SetHomeDir(home)

	dir, err := configStore.Dir()
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
