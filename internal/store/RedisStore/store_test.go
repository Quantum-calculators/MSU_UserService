package RedisStore_test

import (
	"os"
	"testing"
)

var (
	RedisAddr string
)

func TestMain(m *testing.M) {
	RedisAddr = os.Getenv("Redis")
	if RedisAddr == "" {
		RedisAddr = "localhost:6379"
	}
	os.Exit(m.Run())
}
