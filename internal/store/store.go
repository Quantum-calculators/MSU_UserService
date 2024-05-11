package store

import "time"

type Store interface {
	User() UserRepository
}

type RedisStore interface {
	Cache(accessExpTime time.Duration) CacheRepository
}
