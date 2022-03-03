package config

import "fmt"

// Provider is an interface to retrieve the remote URL for a given Provider
type Provider interface {
	// URL returns the complete URL of the repository to clone.
	// Example for a GitHub Provider:
	// URL(true, hello/world) -> git@github.com:hello/world.git
	// URL(false, hello/world) -> https://github.com/hello/world.git
	URL(ssh bool, repository string) string
}

// GitHubProvider implements the Provider interface for GitHub
type GitHubProvider struct {
}

// URL returns the complete URL of the GitHub repository to clone.
func (g GitHubProvider) URL(ssh bool, repository string) string {
	if ssh {
		return fmt.Sprintf("git@github.com:%s.git", repository)
	}

	return fmt.Sprintf("https://github.com/%s.git", repository)
}

// GitLabProvider implements the Provider interface for GitLab
type GitLabProvider struct{}

// URL returns the complete URL of the GitLab repository to clone.
func (g GitLabProvider) URL(ssh bool, repository string) string {
	if ssh {
		return fmt.Sprintf("git@gitlab.com:%s.git", repository)
	}

	return fmt.Sprintf("https://gitlab.com/%s.git", repository)
}

// BitbucketProvider implements the Provider interface for Bitbucket
type BitbucketProvider struct{}

// URL returns the complete URL of the Bitbucket repository to clone.
func (b BitbucketProvider) URL(ssh bool, repository string) string {
	if ssh {
		return fmt.Sprintf("git@bitbucket.com:%s.git", repository)
	}

	return fmt.Sprintf("https://bitbucket.org/%s.git", repository)
}
