package RedisStore

import (
	"context"
	"time"

	"github.com/Quantum-calculators/MSU_UserService/internal/store"
	"github.com/redis/go-redis/v9"
)

type Store struct {
	client        *redis.Client
	JWTRepository *jwtRepository
	ctx           *context.Context
}

func New(client *redis.Client) *Store {
	return &Store{
		client: client,
		ctx:    &ctx,
	}
}

func (s *Store) JWT(jwtSecretKey []byte, accessExpTime time.Duration) store.JWTRepository {
	if s.JWTRepository != nil {
		return s.JWTRepository
	}

	s.JWTRepository = &jwtRepository{
		store:         s,
		jwtSecretKey:  jwtSecretKey,
		accessExpTime: accessExpTime,
	}

	return s.JWTRepository
}
