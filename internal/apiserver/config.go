package apiserver

type Config struct {
	ServerAddr  string `toml:"ServerAddr"`
	LogLevel    string `toml:"LogLevel"`
	DatabaseURL string `toml:"SQLdb_url"`
	Redis       string `toml:"Redis"`
}

func NewConfig() *Config {
	return &Config{
		ServerAddr: "127.0.0.1:8080",
		LogLevel:   "Debug",
	}
}
