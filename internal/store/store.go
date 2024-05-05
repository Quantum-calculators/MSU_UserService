package store

type Store interface {
	User() UserRepository
}

type JWTStore interface {
}
