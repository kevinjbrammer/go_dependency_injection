package config

// IConfig interface
type IConfig interface {
	GetDatabasePath() string
	GetServerPort() string
	ServiceEnabled() bool
}

// Config holds database and service configuration options
type Config struct {
	databasePath string
	port         string
	enabled      bool
}

// NewConfig constructs a Config struct
func NewConfig(databasePath, port string, enabled bool) *Config {
	return &Config{
		databasePath: databasePath,
		port:         port,
		enabled:      enabled,
	}
}

// GetDatabasePath returns the database path
func (c *Config) GetDatabasePath() string {
	return c.databasePath
}

// GetServerPort returns the service port
func (c *Config) GetServerPort() string {
	return c.port
}

// ServiceEnabled indicates whether the service is
// enabled or not
func (c *Config) ServiceEnabled() bool {
	return c.enabled
}
