package store

type Store interface {
	User() UserRepository
}

type RedisStore interface {
	JWT() JWTRepository
}
