package config

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestGitHubProvider_URL(t *testing.T) {
	provider := GitHubProvider{}

	assert.Equal(t, provider.URL(true, "foo/bar"), "git@github.com:foo/bar.git")
	assert.Equal(t, provider.URL(false, "foo/bar"), "https://github.com/foo/bar.git")
}

func TestGitLabProvider_URL(t *testing.T) {
	provider := GitLabProvider{}

	assert.Equal(t, provider.URL(true, "foo/bar"), "git@gitlab.com:foo/bar.git")
	assert.Equal(t, provider.URL(false, "foo/bar"), "https://gitlab.com/foo/bar.git")
}

func TestBitbucketProvider_URL(t *testing.T) {
	provider := BitbucketProvider{}

	assert.Equal(t, provider.URL(true, "foo/bar"), "git@bitbucket.com:foo/bar.git")
	assert.Equal(t, provider.URL(false, "foo/bar"), "https://bitbucket.org/foo/bar.git")
}
