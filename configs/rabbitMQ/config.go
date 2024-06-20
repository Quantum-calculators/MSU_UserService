package rabbitmq

import "fmt"

const defaultPort = "5672"
const defaultHost = "localhost"

type Config struct {
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	User     string `toml:"user"`
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

func (c *Config) GetAMQPaddr() string {
	if c.User != "" {
		if c.Password != "" {
			return fmt.Sprintf("amqp://%s:%s@%s:%s/", c.User, c.Password, c.Host, c.Port)
		}
		return fmt.Sprintf("amqp://%s@%s:%s/", c.User, c.Host, c.Port)
	}
	return fmt.Sprintf("amqp://%s:%s/", c.Host, c.Port)

}
