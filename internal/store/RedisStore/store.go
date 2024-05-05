package RedisStore

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Store struct {
	client     *redis.Client
	ctx        context.Context
	expiration time.Duration
}

func New(client *redis.Client) *Store {
	return &Store{
		client: client,
		ctx:    context.Background(),
	}
}
