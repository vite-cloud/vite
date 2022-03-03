package config

var EmptyConfig = Locator{
	Provider:   GitHubProvider{},
	Path:       "core/domain/config/testdata",
	UseHTTPS:   false,
	Branch:     "main",
	Commit:     "38bcc9d7194ec352647407ee7b053c9bf7ca4bab",
	Repository: "vite-cloud/vite",
}

