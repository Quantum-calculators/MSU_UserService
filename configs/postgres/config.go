package postgres

import "fmt"

const defaultPort = "5432"
const defaultHost = "localhost"

type Config struct {
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DBName   string `toml:"dbname"`
	SSLMode  string `toml:"sslmode"`
}

func (c *Config) WithDefaults() {
	if c.Port == "" {
		c.Port = defaultPort
	}
	if c.Host == "" {
		c.Host = defaultHost
	}
}

func (c *Config) GetSQLaddr() string {
	return fmt.Sprintf(
		"host=%s port=%s password=%s dbname=%s sslmode=%s ",
		c.Host,
		c.Port,
		c.Password,
		c.DBName,
		c.SSLMode,
	)
}
