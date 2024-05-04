package store

type Config struct {
	DBUrl string `toml:"DB_url"`
}

func NewConfig() *Config {
	return &Config{}
}
