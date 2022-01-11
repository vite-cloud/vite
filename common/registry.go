package common

import (
	"encoding/base64"
	"encoding/json"
)

type Registry struct {
	// The name of the registry.
	Name string `yaml:"name"`
	// The URL of the registry.
	Host string `yaml:"host"`
	// The username to use when authenticating with the registry.
	Username string `yaml:"username"`
	// The password to use when authenticating with the registry.
	Password string `yaml:"password"`
}

func (r Registry) IsDefault() bool {
	return r.Name == "" || r.Name == "default" || r.Name == "@"
}

func (r Registry) ToBase64() (string, error) {
	auth, err := json.Marshal(map[string]string{
		"username": r.Username,
		"password": r.Password,
	})

	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(auth), nil
}
