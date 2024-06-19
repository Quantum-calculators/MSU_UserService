package apiserver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/BurntSushi/toml"
	rabbitmq "github.com/Quantum-calculators/MSU_UserService/configs/RabbitMQ"
	postgres "github.com/Quantum-calculators/MSU_UserService/configs/postgres"
	messagebroker "github.com/Quantum-calculators/MSU_UserService/internal/messageBroker/AMQPbroker"
	"github.com/Quantum-calculators/MSU_UserService/internal/store/RedisStore"
	"github.com/Quantum-calculators/MSU_UserService/internal/store/SQLstore"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

func Start(config *Config, postgresConfPath, RabbitConfPath string) error {
	SQLdb, err := MakePostgres(postgresConfPath)
	if err != nil {
		return err
	}
	defer SQLdb.Close()

	RabbitChan, err := MakeBroker(RabbitConfPath)
	if err != nil {
		return err
	}
	defer RabbitChan.Close()

	ctx := context.Background()
	rdb, err := newRedisdb(ctx, config.RedisAddr, config.RedisPas)
	if err != nil {
		return err
	}

	Rstore := RedisStore.New(rdb)
	sqlstore := SQLstore.New(SQLdb, config.ExpRefresh)
	broker := messagebroker.New(RabbitChan)
	srv := newServer(sqlstore, Rstore, broker)

	return http.ListenAndServe(config.ServerAddr, srv)
}

func MakePostgres(configFilePath string) (*sql.DB, error) {
	conf := &postgres.Config{}
	_, err := toml.DecodeFile(configFilePath, conf)
	if err != nil {
		log.Fatal(err)
	}
	conf.WithDefaults()
	db, err := sql.Open("postgres", conf.GetSQLaddr())
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func MakeBroker(configFilePath string) (*amqp.Channel, error) {
	conf := &rabbitmq.Config{}
	_, err := toml.DecodeFile(configFilePath, conf)
	if err != nil {
		return nil, err
	}
	conf.WithDefaults()
	fmt.Println(conf.GetAMQPaddr())
	conn, err := amqp.Dial(conf.GetAMQPaddr())
	if err != nil {
		return nil, errors.New("Failed to connect to AMPQ")
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, errors.New("Failed to connect to AMPQ")
	}
	return ch, nil
}

func newRedisdb(ctx context.Context, Arrd, Password string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     Arrd,
		Password: Password,
	})
	err := rdb.Set(ctx, "key", "value", time.Minute*1).Err()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}
