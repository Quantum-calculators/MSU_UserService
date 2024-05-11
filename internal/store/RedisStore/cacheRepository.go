package RedisStore

import (
	"time"
)

type CacheRepository struct {
	store         *Store
	accessExpTime time.Duration
}

func (j *CacheRepository) CreateAccessToken() (string, error) {
	return "testCreateAccessToken", nil
}

func (j *CacheRepository) CreateRefeshToken() (string, error) {
	return "testCreateRefeshTokne", nil
}

func (j *CacheRepository) UpdateAccessToken() (string, error) {
	return "testUpdateAccessToken", nil
}

func (j *CacheRepository) UpdateRefreshToken() (string, error) {
	return "testUpdateRefreshToken", nil
}
