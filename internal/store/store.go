package store

import "time"

type Store interface {
	User() UserRepository
	Session() SessionRepository
}

type RedisStore interface {
	Cache(accessExpTime time.Duration) CacheRepository
}
