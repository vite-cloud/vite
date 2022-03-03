package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// Git manages a given git repository.
// It holds the path to the repository path.
type Git string

var (
	ErrRepositoryNotFound = errors.New("repository not found")
)

// Clone clones the given repository and branch.
func (g Git) Clone(remote, branch string) error {
	if branch == "" {
		return fmt.Errorf("branch is empty")
	}

	_, err := g.run("clone", "--branch", branch, remote, g.String())
	return err
}

// Read reads a given file at a given revision and returns its contents.
func (g Git) Read(commit, path string) ([]byte, error) {
	return g.run("show", commit+":"+path)
}

// String returns the path to the repository.
func (g Git) String() string {
	return string(g)
}

func (g Git) run(args ...string) ([]byte, error) {
	if _, err := os.Stat(g.String()); errors.Is(err, os.ErrNotExist) {
		return nil, ErrRepositoryNotFound
	} else if err != nil {
		return nil, err
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = g.String()

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	err := cmd.Run()
	out := buf.Bytes()

	if err != nil {
		return nil, fmt.Errorf("%s: %s", err, out)
	}

	return out, nil
}

func (g Git) RepoExists() bool {
	_, err := os.Stat(g.String())
	return !errors.Is(err, os.ErrNotExist)
}
