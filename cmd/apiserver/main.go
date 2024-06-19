package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/Quantum-calculators/MSU_UserService/internal/apiserver"
)

var (
	configPath       string
	PostgresConfPath string
	RabbitMQConfPath string
)

func init() {
	flag.StringVar(&configPath, "server-conf-path", "configs/apiserver.toml", "path to server config file")
	flag.StringVar(&PostgresConfPath, "postgres-conf-path", "configs/postgres/postgres.toml", "path to PostgreSQL config file")
	flag.StringVar(&RabbitMQConfPath, "rabbitmq-conf-path", "configs/rabbitMQ/rabbitMQ.toml", "path to RabbitMQ config file")
}

func main() {
	flag.Parse()

	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := apiserver.Start(config, PostgresConfPath, RabbitMQConfPath); err != nil {
		log.Fatal(err)
	}
}
