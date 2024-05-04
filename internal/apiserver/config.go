package apiserver

import "github.com/PepsiKingIV/Lib_REST_API_server/internal/store"

type Config struct {
	ServerAddr string `toml:"ServerAddr"`
	LogLevel   string `toml:"LogLevel"`
	Store      *store.Config
}

func NewConfig() *Config {
	return &Config{
		ServerAddr: "127.0.0.1:8080",
		LogLevel:   "Debug",
		Store:      store.NewConfig(),
	}
}
