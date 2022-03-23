package locator

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestProvider_Name(t *testing.T) {
	p := Provider("github")

	assert.Equal(t, p.Name(), "github")
}

func TestProvider_URL(t *testing.T) {
	p := Provider("github")

	assert.Equal(t, p.URL("ssh", "foo/bar"), "git@github.com:foo/bar.git")
	assert.Equal(t, p.URL("https", "foo/bar"), "https://github.com/foo/bar.git")
	assert.Equal(t, p.URL("", "foo/bar"), "git://github.com/foo/bar.git")
}
