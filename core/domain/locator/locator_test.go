package locator

import (
	"github.com/vite-cloud/vite/core/domain/datadir"
	"gotest.tools/v3/assert"
	"os"
	"path/filepath"
	"testing"
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
	builder.
		WriteFile("vite.yaml", []byte("services:\n"), 0600).
		Commit()

	contents, err := locator.Read("vite.yaml")
	assert.NilError(t, err)

	assert.Equal(t, string(contents), "services:\n")
}
