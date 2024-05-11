package RedisStore

import (
	"time"

	"github.com/Quantum-calculators/MSU_UserService/internal/store"
	"github.com/redis/go-redis/v9"
)

type Store struct {
	client *redis.Client
	cache  *CacheRepository
}

func New(client *redis.Client) *Store {
	return &Store{
		client: client,
	}
}

func (s *Store) Cache(accessExpTime time.Duration) store.CacheRepository {
	if s.cache != nil {
		return s.cache
	}

	s.cache = &CacheRepository{
		store:         s,
		accessExpTime: accessExpTime,
	}

	return s.cache
}
