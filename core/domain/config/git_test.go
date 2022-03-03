package config

import (
	"gotest.tools/v3/assert"
	"os"
	"testing"
)

func TestGit_Read(t *testing.T) {
	builder := newLocalRepo(t)
	builder.WriteFile("vite.yml", []byte("here is some content"), 0600)
	commit := builder.Commit()
	git := builder.Git()

	content, err := git.Read(commit, "vite.yml")
	assert.NilError(t, err)

	assert.Equal(t, string(content), "here is some content")
}

func TestGit_Read2(t *testing.T) {
	builder := newLocalRepo(t)
	commit := builder.
		WriteFile("something", []byte(""), 0600).
		Commit()
	git := builder.Git()

	_, err := git.Read(commit, "vite.yml")
	assert.ErrorContains(t, err, "fatal: path 'vite.yml' does not exist")
}

func TestGit_String(t *testing.T) {
	assert.Equal(t, Git("hello_world").String(), "hello_world")
}

func TestGit_Read3(t *testing.T) {
	_, err := Git("/tmp/does-not-exist").Read("fffffff", "vite.yml")
	assert.ErrorIs(t, err, ErrRepositoryNotFound)
}

func TestGit_RepoExists(t *testing.T) {
	ok := Git("/tmp/does-not-exist").RepoExists()
	assert.Assert(t, !ok)

	dir, err := os.MkdirTemp("", "git-repo")
	assert.NilError(t, err)

	ok = Git(dir).RepoExists()
	assert.Assert(t, ok)
}

func Test