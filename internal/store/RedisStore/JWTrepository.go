package RedisStore

import "time"

type JWTrepository struct {
	rdb *Store
}

func (j *JWTrepository) Set(key string, value interface{}) error {
	j.rdb.client.Set(j.rdb.ctx, key, value, time.Hour*1)
	return nil
}
