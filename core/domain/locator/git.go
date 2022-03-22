package locator

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/vite-cloud/vite/core/domain/log"
	"os"
	"os/exec"
	"strings"
)

// Git manages a given git repository.
// It holds the path to the repository path.
type Git string

// related errors
var (
	ErrRepositoryNotFound = errors.New("repository not found")
	ErrEmptyBranch        = errors.New("empty branch")
)

// Clone clones the given repository and branch.
func (g Git) Clone(remote, branch string) error {
	if branch == "" {
		return ErrEmptyBranch
	}

	cmd := exec.Command("git", "clone", "--branch", branch, remote, g.String())
	_, err := globalRun(cmd)
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

func globalRun(cmd *exec.Cmd) ([]byte, error) {
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	err := cmd.Run()
	out := buf.Bytes()

	log.Log(log.DebugLevel, "ran git command", log.Fields{
		"cmd":  cmd.Args,
		"err":  err,
		"code": cmd.ProcessState.ExitCode(),
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %s", err, out)
	}

	return out, nil
}

// run runs a git command with the given arguments.
func (g Git) run(args ...string) ([]byte, error) {
	if _, err := os.Stat(g.String()); errors.Is(err, os.ErrNotExist) {
		return nil, ErrRepositoryNotFound
	} else if err != nil {
		return nil, err
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = g.String()

	return globalRun(cmd)
}

// RepoExists returns true if a repository exists at the given path.
func (g Git) RepoExists() bool {
	_, err := os.Stat(g.String())
	return !errors.Is(err, os.ErrNotExist)
}

type CommitList []Commit

type Commit struct {
	Hash    string
	Message string
}

// Commits returns the list of commits in the given branch.
func (g Git) Commits(branch string) (CommitList, error) {
	if branch == "" {
		return nil, ErrEmptyBranch
	}

	out, err := g.run("log", "--pretty=%H%s", "--no-merges", branch)
	if err != nil {
		return nil, err
	}

	var commits CommitList

	for _, line := range strings.Split(string(out), "\n") {
		if line == "" {
			continue
		}

		hash := line[:40]
		message := line[40:]

		commits = append(commits, Commit{
			Hash:    hash,
			Message: message,
		})
	}

	return commits, nil
}

func (c CommitList) AsOptions() []string {
	var options []string

	for _, commit := range c {
		options = append(options, fmt.Sprintf("%s %s", commit.Hash, commit.Message))
	}

	return options
}
