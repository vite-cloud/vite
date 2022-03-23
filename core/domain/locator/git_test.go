package locator

import (
	"io/fs"
	"os"
	"syscall"
	"testing"

	"gotest.tools/v3/assert"
)

func TestGit_Read(t *testing.T) {
	dir, err := os.MkdirTemp("", "git_test")
	assert.NilError(t, err)

	builder := newLocalRepo(t, dir)
	builder.WriteFile("vite.yml", []byte("here is some content"), 0600)
	commit := builder.Commit()
	git := builder.Git()

	content, err := git.Read(commit, "vite.yml")
	assert.NilError(t, err)

	assert.Equal(t, string(content), "here is some content")
}

func TestGit_Read2(t *testing.T) {
	dir, err := os.MkdirTemp("", "git_test")
	assert.NilError(t, err)

	builder := newLocalRepo(t, dir)
	commit := builder.
		WriteFile("something", []byte(""), 0600).
		Commit()
	git := builder.Git()

	_, err = git.Read(commit, "vite.yml")
	assert.ErrorContains(t, err, "'vite.yml' does not exist")
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

func TestGit_Clone(t *testing.T) {
	dir, err := os.MkdirTemp("", "git-repo")
	assert.NilError(t, err)

	git := Git(dir)
	err = git.Clone("https://github.com/vite-cloud/vite.git", "main")
	assert.NilError(t, err)
}

func TestGit_Clone2(t *testing.T) {
	dir, err := os.MkdirTemp("", "git-repo")
	assert.NilError(t, err)

	git := Git(dir)
	err = git.Clone("nop does not exist", "main")
	assert.ErrorContains(t, err, "fatal: repository 'nop does not exist' does not exist")
}

func TestCommitList_AsOptions(t *testing.T) {
	list := CommitList{
		{Hash: "123", Message: "hello"},
		{Hash: "456", Message: "world"},
	}

	assert.DeepEqual(t, list.AsOptions(), []string{"123 hello", "456 world"})
}

func TestGit_Commits(t *testing.T) {
	_, err := Git("/nop").Commits("main")
	assert.ErrorIs(t, err, ErrRepositoryNotFound)
}

func TestGit_Commits2(t *testing.T) {
	_, err := Git("/nop").Commits("")
	assert.ErrorIs(t, err, ErrEmptyBranch)
}

func TestGit_Run(t *testing.T) {
	_, err := Git("\x00").run("init")
	assert.DeepEqual(t, err, &fs.PathError{Op: "stat", Err: syscall.Errno(0x16), Path: "\x00"})
}
