package locator

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
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

	dir, err := Store.Dir()
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
		Protocol:   "https",
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
		Protocol:   "https",
		Repository: "felixdorn/config-test",
	}

	_, err := locator.Read("does-not-exist.yaml")
	assert.ErrorContains(t, err, "path 'does-not-exist.yaml' does not exist")
}

func TestLocator_Read6(t *testing.T) {
	datadir.SetHomeDir("/nop")

	locator := Locator{
		Commit: "ffffffffffffffffffffffffffffffffffffffff",
	}

	_, err := locator.Read("vite.yaml")
	assert.ErrorIs(t, err, os.ErrPermission)
}

func TestLocator_Read7(t *testing.T) {
	datadir.UseTestHome(t)

	locator := Locator{
		Commit:   "8897a7d08a1e791418904afdce369818c19d2c3e",
		Branch:   "main",
		Provider: Provider("not.example"),
	}

	_, err := locator.Read("vite.yaml")
	assert.ErrorContains(t, err, "unable to look up not.example.com (port 9418) (Name or service not known)")
}

func TestLocator_Git(t *testing.T) {
	datadir.SetHomeDir("/nop")

	locator := Locator{}

	_, err := locator.git()
	assert.ErrorIs(t, err, os.ErrPermission)
}

func TestLocator_Save(t *testing.T) {
	datadir.SetHomeDir("/nop")

	locator := Locator{}

	err := locator.Save()
	assert.ErrorIs(t, err, os.ErrPermission)
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

func TestLoadFromStore2(t *testing.T) {
	datadir.SetHomeDir("/nop")

	_, err := LoadFromStore()
	assert.ErrorIs(t, err, os.ErrPermission)
}

func TestLoadFromStore3(t *testing.T) {
	datadir.UseTestHome(t)

	_, err := LoadFromStore()
	assert.Error(t, err, "config locator hasn't been configured yet, run `vite setup` first")
}

func TestLoadFromStore4(t *testing.T) {
	datadir.UseTestHome(t)

	f, err := Store.Open(ConfigFile, os.O_CREATE|os.O_WRONLY, 0600)
	assert.NilError(t, err)

	defer f.Close()

	_, err = f.WriteString("invalid JSON")
	assert.NilError(t, err)

	_, err = LoadFromStore()
	assert.Error(t, err, "invalid character 'i' looking for beginning of value")
}

func TestLocator_Checksum(t *testing.T) {
	t.Parallel()

	locator := Locator{
		Provider:   Provider("foo"),
		Protocol:   "ssh",
		Repository: "foo/bar",
		Branch:     "main",
		Commit:     "fffffff",
		Path:       "/sub/path",
	}

	data, err := json.Marshal(locator)
	assert.NilError(t, err)

	var values map[string]interface{}
	err = json.Unmarshal(data, &values)
	assert.NilError(t, err)

	sum, err := base64.StdEncoding.DecodeString(locator.Checksum())
	assert.NilError(t, err)

	for k, v := range values {
		if k == "protocol" {
			continue
		}

		assert.Assert(t, bytes.Contains(sum, []byte(v.(string))), "sum: %s key: %s value: %s", sum, k, v)
	}
}

func TestLocator_Commits(t *testing.T) {
	datadir.UseTestHome(t)

	locator := Locator{
		Branch:     "main",
		Repository: "foo/bar",
	}

	dir, err := Store.Dir()
	assert.NilError(t, err)

	repo := newLocalRepo(t, dir+"/main-foo-bar")

	first := repo.WriteFile("hello-world", []byte{}, 0600).Commit()
	second := repo.WriteFile("2hello-world", []byte{}, 0600).Commit()

	commits, err := locator.Commits()
	assert.NilError(t, err)

	assert.Equal(t, len(commits), 2)
	assert.Equal(t, commits[0].Hash, second)
	assert.Equal(t, commits[0].Message, "commit")
	assert.Equal(t, commits[1].Hash, first)
	assert.Equal(t, commits[1].Message, "commit")
}

func TestLocator_Commits2(t *testing.T) {
	datadir.SetHomeDir("/nop")

	locator := Locator{}

	_, err := locator.Commits()
	assert.ErrorIs(t, err, os.ErrPermission)
}
