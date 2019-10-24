package config

// Config holds database and service configuration options
type Config struct {
	DatabasePath string
	Port         string
	Enabled      bool
}

// NewConfig constructs a Config struct
func NewConfig() *Config {
	return &Config{
		DatabasePath: "./spaceships.db",
		Port:         "8000",
		Enabled:      true,
	}
}
