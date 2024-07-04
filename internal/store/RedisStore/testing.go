package redisStore

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func TestRedis(t *testing.T, redisaddr string) {
	t.Helper()

	rdb := redis.NewClient(&redis.Options{
		Addr: redisaddr,
	})
	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		t.Fatal(err)
	}
}
