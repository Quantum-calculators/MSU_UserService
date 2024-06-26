package apiserver

import "fmt"

const (
	defaultPort         = "8080"
	defaultHost         = "localhost"
	defaultLogLevel     = "Debug"
	defaultExpRefresh   = 30
	defaultExpAccess    = 5
	defaultQueryTimeOut = 500 //millisecond
)

type Config struct {
	Port         string `toml:"port"`
	Host         string `toml:"host"`
	LogLevel     string `toml:"loglevel"`
	ExpRefresh   int    `toml:"exp_refresh_token_in_min"`
	ExpAccess    int    `toml:"exp_access_token_in_min"`
	JwtSecretKey string `toml:"jwt_secret_key"`
	QueryTimeOut int    `toml:"query_timeout"` //millisecond
}

func (c *Config) WithDefaults() {
	if c.Port == "" {
		c.Port = defaultPort
	}
	if c.Host == "" {
		c.Host = defaultHost
	}
	if c.LogLevel == "" {
		c.LogLevel = defaultLogLevel
	}
	if c.ExpAccess == 0 {
		c.ExpAccess = defaultExpAccess
	}
	if c.ExpRefresh == 0 {
		c.ExpRefresh = defaultExpAccess
	}
	if c.QueryTimeOut == 0 {
		c.QueryTimeOut = defaultQueryTimeOut
	}
}

func (c *Config) GenServerAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
