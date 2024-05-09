package RedisStore

import (
	"github.com/Quantum-calculators/MSU_UserService/internal/store"
	"github.com/redis/go-redis/v9"
)

type Store struct {
	client        *redis.Client
	JWTRepository *jwtRepository
}

func New(client *redis.Client) *Store {
	return &Store{
		client: client,
	}
}

func (s *Store) JWT() store.JWTRepository {
	if s.JWTRepository != nil {
		return s.JWTRepository
	}

	s.JWTRepository = &jwtRepository{
		store: s,
	}

	return s.JWTRepository
}
