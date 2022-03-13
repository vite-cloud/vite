package config

// Hooks contains a service's hooks.
// It is used both my the configYAML and the Config
type Hooks struct {
	// Commands to run before the container is started.
	Prestart []string `json:"prestart" yaml:"prestart"`
	// Commands to run after the container is started.
	Poststart []string `json:"poststart" yaml:"poststart"`
	// Commands to run before the container is stopped.
	Prestop []string `json:"prestop" yaml:"prestop"`
	// Commands to run after the container is stopped.
	Poststop []string `json:"poststop" yaml:"poststop"`
}
