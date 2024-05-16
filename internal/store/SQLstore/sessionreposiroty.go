package SQLstore

import (
	"time"

	"github.com/Quantum-calculators/MSU_UserService/internal/model"
	"github.com/Quantum-calculators/MSU_UserService/internal/token_generator"
)

type SessionRepository struct {
	store *Store
}

func (s *SessionRepository) CreateSession(userId uint32, fingerpring string) (*model.Session, error) {
	refreshToken, err := token_generator.GenerateRandomString(128)
	if err != nil {
		return &model.Session{}, err
	}
	session := &model.Session{
		UserId:       userId,
		RefreshToken: refreshToken,
		Fingerprint:  fingerpring,
		ExpiresIn:    time.Now().Add(time.Minute * time.Duration(s.store.ExpRefreshToken)).Unix(),
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

func (s *SessionRepository) VerifyRefreshToken(fingerPrint, refreshToken string) (*model.Session, error) {
	var ID int
	var user_id int
	if err := s.store.db.QueryRow(
		"SELECT id, user_id FROM sessions WHERE fingerprint = $1 AND refresh_token = $2;",
		fingerPrint,
		refreshToken,
	).Scan(
		&ID,
		&user_id,
	); err != nil {
		return &model.Session{}, err
	}
	session, err := s.CreateSession(uint32(user_id), fingerPrint)
	if err != nil {
		return &model.Session{}, err
	} // TODO: добавить поле использован(t/f), чтобы проверять не украден ли токен.
	return session, nil
}

// TODO: добавить удаление сессии по истечению
// TODO: добавить /logout  (удалять сессию)
