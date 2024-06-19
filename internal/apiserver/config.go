package apiserver

type Config struct {
	ServerAddr string `toml:"ServerAddr"`
	LogLevel   string `toml:"LogLevel"`
	RedisAddr  string `toml:"RedisAddr"`
	RedisPas   string `toml:"RedisPas"`
	ExpRefresh int    `toml:"ExpRefreshTokenInMin"`
	ExpAccess  int    `toml:"ExpAccessTokenInMin"`
}

func NewConfig() *Config {
	return &Config{
		ServerAddr: "127.0.0.1:8080",
		LogLevel:   "Debug",
	}
}
