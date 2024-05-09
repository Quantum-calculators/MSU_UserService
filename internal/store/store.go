package store

import "time"

type Store interface {
	User() UserRepository
}

type RedisStore interface {
	JWT(jwtSecretKey []byte, accessExpTime time.Duration) JWTRepository
}
