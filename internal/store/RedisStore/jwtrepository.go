package RedisStore

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type jwtRepository struct {
	store         *Store
	jwtSecretKey  []byte
	accessExpTime time.Duration
}

func (j *jwtRepository) CreateAccessToken(email string) (string, error) {
	payload := jwt.MapClaims{
		"sub": email,
		"exp": time.Now().Add(j.accessExpTime).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString(j.jwtSecretKey)
	if err != nil {
		return "", err
	}
	j.store.client.Set(*j.store.ctx, email, token, j.accessExpTime)
	return t, nil
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
