package locator

import "fmt"

// Provider is a locator provider such as GitHub, GitLab...
type Provider string

// Name returns the provider name in lowercase and without spaces (e.g. github, gitlab...)
func (p *Provider) Name() string {
	return string(*p)
}

// URL returns the repository url given a protocol (ssh, https) and a repository name
func (p *Provider) URL(protocol, repository string) string {
	switch protocol {
	case "ssh":
		return fmt.Sprintf("git@%s.com:%s.git", p.Name(), repository)
	case "https":
		return fmt.Sprintf("https://%s.com/%s.git", p.Name(), repository)
	default:
		return fmt.Sprintf("git://%s.com/%s.git", p.Name(), repository)
	}
}
