package locator

import "fmt"

type Provider string

func (p *Provider) Name() string {
	return string(*p)
}

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

var GitHubProvider = Provider("github")

var GitLabProvider = Provider("gitlab")

var BitbucketProvider = Provider("bitbucket")
