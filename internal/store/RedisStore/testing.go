package RedisStore

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func TestRedis(t *testing.T, RedisAddr string) *redis.Client {
	t.Helper()

	rdb := redis.NewClient(&redis.Options{
		Addr: RedisAddr,
	})
	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		t.Fatal(err)
	}
	return rdb
}
