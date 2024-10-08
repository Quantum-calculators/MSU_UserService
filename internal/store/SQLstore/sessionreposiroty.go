package SQLstore

import (
	"context"
	"database/sql"
	"time"

	"github.com/Quantum-calculators/MSU_UserService/internal/model"
	"github.com/Quantum-calculators/MSU_UserService/internal/store"
	token_generator "github.com/Quantum-calculators/MSU_UserService/internal/tokenGenerator"
)

type SessionRepository struct {
	store *Store
}

func (s *SessionRepository) CreateSession(ctxb context.Context, email string, fingerpring string) (*model.Session, error) {
	ctx, cancel := context.WithTimeout(ctxb, s.store.QueryTimeout)
	defer cancel()
	refreshToken, err := token_generator.GenerateRandomString(128)
	if err != nil {
		return &model.Session{}, err
	}
	session := &model.Session{
		Email:        email,
		RefreshToken: refreshToken,
		Fingerprint:  fingerpring,
		ExpiresIn:    time.Now().Add(time.Minute * time.Duration(s.store.ExpRefreshToken)).Unix(),
		CreatedAt:    time.Now().Unix(),
	}
	if err := s.store.db.QueryRowContext(
		ctx,
		"INSERT INTO sessions (email, refresh_token, fingerprint, expires_in, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id;",
		session.Email,
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

func (s *SessionRepository) VerifyRefreshToken(ctxb context.Context, fingerPrint, refreshToken string) (*model.Session, error) {
	ctx, cancel := context.WithTimeout(ctxb, s.store.QueryTimeout)
	defer cancel()
	session := &model.Session{}
	newRefreshToken, err := token_generator.GenerateRandomString(128)
	if err != nil {
		return &model.Session{}, err
	}
	transaction, err := s.store.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return &model.Session{}, err
	}
	if err := transaction.QueryRowContext(ctx,
		"SELECT id, email, expires_in FROM sessions WHERE fingerprint = $1 AND refresh_token = $2;",
		fingerPrint,
		refreshToken,
	).Scan(
		&session.ID,
		&session.Email,
		&session.ExpiresIn,
	); err != nil {
		return &model.Session{}, err
	}
	if session.ExpiresIn < time.Now().Unix() {
		return &model.Session{}, store.ErrRefreshTokenExpired
	}
	newSession := &model.Session{
		Email:        session.Email,
		RefreshToken: newRefreshToken,
		Fingerprint:  fingerPrint,
		ExpiresIn:    time.Now().Add(time.Minute * time.Duration(s.store.ExpRefreshToken)).Unix(),
		CreatedAt:    time.Now().Unix(),
	}
	if err := transaction.QueryRowContext(
		ctx,
		"INSERT INTO sessions (email, refresh_token, fingerprint, expires_in, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id;",
		newSession.Email,
		newRefreshToken,
		fingerPrint,
		newSession.ExpiresIn,
		newSession.CreatedAt,
	).Scan(
		&newSession.ID,
	); err != nil {
		return &model.Session{}, err
	}
	if err := transaction.Commit(); err != nil {
		return &model.Session{}, err
	}
	// TODO: добавить поле использован(t/f), чтобы проверять не украден ли токен.
	return newSession, nil
}

func (s *SessionRepository) DeleteSession(ctxb context.Context, fingerPrint, refreshToken string) error {
	ctx, cancel := context.WithTimeout(ctxb, s.store.QueryTimeout)
	defer cancel()
	if err := s.store.db.QueryRowContext(ctx,
		"DELETE FROM sessions WHERE fingerprint = $1 AND refresh_token = $2;",
		fingerPrint,
		refreshToken,
	).Err(); err != nil {
		return err
	}
	return nil
}

func (s *SessionRepository) DeleteAllSession(ctxb context.Context, email string) error {
	ctx, cancel := context.WithTimeout(ctxb, s.store.QueryTimeout)
	defer cancel()
	if err := s.store.db.QueryRowContext(ctx,
		"DELETE FROM sessions WHERE email = $1;",
		email,
	).Err(); err != nil {
		return err
	}
	return nil
}

// TODO: добавить удаление сессии по истечению
