package redis

import "fmt"

const (
	defaultPort = "6379"
	defaultHost = "localhost"
)

type Config struct {
	Port     string `toml:"port"`
	Host     string `toml:"host"`
	Password string `toml:"password"`
}

func (c *Config) WithDefaults() {
	if c.Port == "" {
		c.Port = defaultPort
	}
	if c.Host == "" {
		c.Host = defaultHost
	}
}

func (c *Config) GenServerAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
