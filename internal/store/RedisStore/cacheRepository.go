package RedisStore

import (
	"time"
)

type CacheRepository struct {
	store         *Store
	accessExpTime time.Duration
}

func (j *CacheRepository) Set() (string, error) {
	return "", nil
}

func (j *CacheRepository) Get() (string, error) {
	return "", nil
}
