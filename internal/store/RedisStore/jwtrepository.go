package RedisStore

type jwtRepository struct {
	store *Store
}

func (j *jwtRepository) CreateAccessToken() (string, error) {
	return "testCreateAccessToken", nil
}

func (j *jwtRepository) CreateRefeshToken() (string, error) {
	return "testCreateRefeshTokne", nil
}

func (j *jwtRepository) UpdateAccessToken() (string, error) {
	return "testUpdateAccessToken", nil
}

func (j *jwtRepository) UpdateRefreshToken() (string, error) {
	return "testUpdateRefreshToken", nil
}
