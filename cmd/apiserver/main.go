package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	apiserverConf "github.com/Quantum-calculators/MSU_UserService/configs/apiserver"
	"github.com/Quantum-calculators/MSU_UserService/internal/apiserver"
)

var (
	configPath       string
	PostgresConfPath string
	RabbitMQConfPath string
	RedisConfPath    string
)

func init() {
	flag.StringVar(&configPath, "server-conf-path", "configs/apiserver/apiserver.toml", "path to server config file")
	flag.StringVar(&PostgresConfPath, "postgres-conf-path", "configs/postgres/postgres.toml", "path to PostgreSQL config file")
	flag.StringVar(&RabbitMQConfPath, "rabbitmq-conf-path", "configs/rabbitMQ/rabbitMQ.toml", "path to RabbitMQ config file")
	flag.StringVar(&RedisConfPath, "redis-conf-path", "configs/redis/redis.toml", "path to Radis config file")
}

func main() {
	flag.Parse()
	conf := &apiserverConf.Config{}
	_, err := toml.DecodeFile(configPath, conf)
	if err != nil {
		log.Fatal(err)
	}
	conf.WithDefaults()

	if err := apiserver.Start(
		conf,
		PostgresConfPath,
		RabbitMQConfPath,
		RedisConfPath,
	); err != nil {
		log.Fatal(err)
	}
}
