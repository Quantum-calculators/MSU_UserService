package apiserver

type Config struct {
	ServerAddr  string `toml:"ServerAddr"`
	LogLevel    string `toml:"LogLevel"`
	DatabaseURL string `toml:"SQLdb_url"`
	RedisAddr   string `toml:"RedisAddr"`
	RedisPas    string `toml:"RedisPas"`
	ExpRefresh  int    `toml:"ExpRefreshTokenInMin"`
	ExpAccess   int    `toml:"ExpAccessTokenInMin"`
	AMQPaddr    string `toml:"AMQPaddr"`
}

func NewConfig() *Config {
	return &Config{
		ServerAddr: "127.0.0.1:8080",
		LogLevel:   "Debug",
	}
}
