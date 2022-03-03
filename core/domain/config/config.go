package config

// Config holds vite's configuration.
type Config struct {
}

// NewConfig creates a new Config using a config Locator.
func NewConfig(locator *Locator) (*Config, error) {
	_, err := locator.Read("vite.yaml")
	if err != nil {
		return nil, err
	}

	return &Config{}, nil
}
