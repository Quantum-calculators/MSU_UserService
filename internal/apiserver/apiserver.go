package apiserver

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	messagebroker "github.com/Quantum-calculators/MSU_UserService/internal/messageBroker"
	"github.com/Quantum-calculators/MSU_UserService/internal/store/RedisStore"
	"github.com/Quantum-calculators/MSU_UserService/internal/store/SQLstore"
	"github.com/redis/go-redis/v9"
)

func Start(config *Config) error {
	SQLdb, err := newSQLdb(config.DatabaseURL)
	if err != nil {
		return err
	}
	defer SQLdb.Close()

	ctx := context.Background()
	rdb, err := newRedisdb(ctx, config.RedisAddr, config.RedisPas)
	if err != nil {
		return err
	}
	Rstore := RedisStore.New(rdb)
	sqlstore := SQLstore.New(SQLdb, config.ExpRefresh)
	broker := messagebroker.NewAMQL(config.AMQPaddr)
	srv := newServer(sqlstore, Rstore, broker)

	return http.ListenAndServe(config.ServerAddr, srv)
}

func newSQLdb(DatabaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", DatabaseURL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
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
