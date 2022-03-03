package config

import (
	"gotest.tools/v3/assert"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type File struct {
	Name     string
	Contents []byte
	Perm     os.FileMode
}

type RepoBuilder struct {
	t     *testing.T
	path  string
	files []File
}

func newLocalRepo(t *testing.T) *RepoBuilder {
	repo, err := os.MkdirTemp("", "git-repo")
	assert.NilError(t, err)

	runGit(t, repo, "init")

	return &RepoBuilder{
		t:    t,
		path: repo,
	}
}

func runGit(t *testing.T, repo string, args ...string) []byte {
	cmd := exec.Command("git", args...)
	cmd.Dir = repo

	out, err := cmd.CombinedOutput()
	assert.NilError(t, err, string(out))

	return out
}

func (r *RepoBuilder) WriteFile(name string, contents []byte, perm os.FileMode) *RepoBuilder {
	r.files = append(r.files, File{
		Name:     name,
		Contents: contents,
		Perm:     perm,
	})
	return r
}

func (r *RepoBuilder) Commit() string {
	for _, file := range r.files {
		filePath := filepath.Join(r.path, file.Name)

		err := os.MkdirAll(filepath.Dir(filePath), 0755)
		assert.NilError(r.t, err)

		err = ioutil.WriteFile(filePath, file.Contents, file.Perm)
		assert.NilError(r.t, err)
	}

	r.files = nil

	name, email := string(runGit(r.t, r.path, "config", "user.name")), string(runGit(r.t, r.path, "config", "user.email"))

	if strings.TrimSpace(name) == "" {
		runGit(r.t, r.path, "config", "--local", "user.name", "Test User")
	}

	if strings.TrimSpace(email) == "" {
		runGit(r.t, r.path, "config", "--local", "user.email", "commiter@example.com")
	}

	runGit(r.t, r.path, "add", ".")

	runGit(r.t, r.path, "commit", "-m", "commit")

	// get last commit
	out := runGit(r.t, r.path, "log", "-1", "--pretty=format:%H")

	return strings.TrimSpace(string(out))
}

func (r *RepoBuilder) Git() Git {
	return Git(r.path)
}
