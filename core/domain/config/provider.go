package config

import "fmt"

type provider interface {
	// URL returns the complete URL of the repository to clone.
	// Example for a GitHub provider:
	// URL(true, hello/world) -> git@github.com:hello/world.git
	// URL(false, hello/world) -> https://github.com/hello/world.git
	URL(ssh bool, repository string) string
}

type GitHubProvider struct {
}

func (g GitHubProvider) URL(ssh bool, repository string) string {
	if ssh {
		return fmt.Sprintf("git@github.com:%s.git", repository)
	}

	return fmt.Sprintf("https://github.com/%s.git", repository)
}

type GitLabProvider struct{}

func (g GitLabProvider) URL(ssh bool, repository string) string {
	if ssh {
		return fmt.Sprintf("git@gitlab.com:%s.git", repository)
	}

	return fmt.Sprintf("https://gitlab.com/%s.git", repository)
}

type BitbucketProvider struct{}

func (b BitbucketProvider) URL(ssh bool, repository string) string {
	if ssh {
		return fmt.Sprintf("git@bitbucket.com:%s.git", repository)
	}

	return fmt.Sprintf("https://bitbucket.org/%s.git", repository)
}
