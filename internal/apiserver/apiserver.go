package apiserver

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/BurntSushi/toml"
	apiserverConf "github.com/Quantum-calculators/MSU_UserService/configs/apiserver"
	postgres "github.com/Quantum-calculators/MSU_UserService/configs/postgres"
	rabbitmq "github.com/Quantum-calculators/MSU_UserService/configs/rabbitMQ"
	redisConf "github.com/Quantum-calculators/MSU_UserService/configs/redis"
	messagebroker "github.com/Quantum-calculators/MSU_UserService/internal/messageBroker/AMQPbroker"
	"github.com/Quantum-calculators/MSU_UserService/internal/store/SQLstore"
	RedisStore "github.com/Quantum-calculators/MSU_UserService/internal/store/redisStore"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

func Start(config *apiserverConf.Config, postgresConfPath, RabbitConfPath, RedisConfPath string) error {
	SQLdb, err := MakePostgres(postgresConfPath, config.DBMaxOpenConns)
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
	rdb, err := MakeRedis(ctx, RedisConfPath)
	if err != nil {
		return err
	}

	Rstore := RedisStore.New(rdb)
	sqlstore := SQLstore.New(SQLdb, config.ExpRefresh, time.Millisecond*time.Duration(config.QueryTimeOut))
	broker := messagebroker.New(RabbitChan)
	srv := newServer(sqlstore, Rstore, broker, config.ExpAccess, config.JwtSecretKey)
	loggerMiddleware := srv.Logging(srv)
	PanicHookMiddleware := srv.PanicRecoveryMiddleware(loggerMiddleware)

	return http.ListenAndServe(config.GenServerAddr(), PanicHookMiddleware)
}

func MakePostgres(configFilePath string, maxOpenConns int) (*sql.DB, error) {
	conf := &postgres.Config{}
	_, err := toml.DecodeFile(configFilePath, conf)
	if err != nil {
		log.Fatal(err)
	}
	conf.WithDefaults()
	db, err := sql.Open("postgres", conf.GetSQLaddr())
	db.SetMaxOpenConns(maxOpenConns)
	if err != nil {
		return nil, err
	}
	db.QueryRow("")
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
	conn, err := amqp.Dial(conf.GetAMQPaddr())
	if err != nil {
		return nil, errors.New("failed to connect to AMPQ")
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, errors.New("failed to connect to AMPQ")
	}
	return ch, nil
}

func MakeRedis(ctx context.Context, configFilePath string) (*redis.Client, error) {
	conf := &redisConf.Config{}
	_, err := toml.DecodeFile(configFilePath, conf)
	if err != nil {
		log.Fatal(err)
	}
	conf.WithDefaults()

	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.GenServerAddr(),
		Password: conf.Password,
	})
	err = rdb.Set(ctx, "key", "value", time.Minute*1).Err()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}
