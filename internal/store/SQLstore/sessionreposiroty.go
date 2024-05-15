package SQLstore

import (
	"time"

	"github.com/Quantum-calculators/MSU_UserService/internal/model"
	"github.com/Quantum-calculators/MSU_UserService/internal/token_generator"
)

type SessionRepository struct {
	store *Store
}

func (s *SessionRepository) CreateSession(userId uint32, fingerpring string, ExpRefreshToken int) (*model.Session, error) {
	refreshToken, err := token_generator.GenerateRandomString(128)
	if err != nil {
		return &model.Session{}, err
	}
	session := &model.Session{
		UserId:       userId,
		RefreshToken: refreshToken,
		Fingerprint:  fingerpring,
		ExpiresIn:    time.Now().Add(time.Duration(6 * 10e10)).Unix(),
		CreatedAt:    time.Now().Unix(),
	}
	if err := s.store.db.QueryRow(
		"INSERT INTO sessions (user_id, refresh_token, fingerprint, expires_in, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id;",
		session.UserId,
		session.RefreshToken,
		session.Fingerprint,
		session.ExpiresIn,
		session.CreatedAt,
	).Scan(
		&session.ID,
	); err != nil {
		return &model.Session{}, err
	}
	return session, nil
}

func (s *SessionRepository) VerifyRefreshToken(userID int, fingerPrint, refreshToken string) (string, error) {
	ID := 0
	if err := s.store.db.QueryRow(
		"SELECT id FROM sessions WHERE user_id = $1 AND fingerprint = $2 AND refresh_token = $3;",
		userID,
		fingerPrint,
		refreshToken,
	).Scan(
		&ID,
	); err != nil {
		return "", err
	}
	newRefreshToken, err := token_generator.GenerateRandomString(128)
	if err != nil {
		return "", err
	}
	return newRefreshToken, nil
}
