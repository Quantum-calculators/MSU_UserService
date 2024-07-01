package testStore

import (
	"context"
	"errors"
	"time"

	"github.com/Quantum-calculators/MSU_UserService/internal/model"
	token_generator "github.com/Quantum-calculators/MSU_UserService/internal/tokenGenerator"
)

type SessionRepository struct {
	store    *Store
	sessions map[string]*model.Session
}

func (s *SessionRepository) CreateSession(cxt context.Context, email string, fingerpring string) (*model.Session, error) {
	refreshToken, err := token_generator.GenerateRandomString(64)
	if err != nil {
		return &model.Session{}, err
	}
	session := &model.Session{
		Email:        email,
		RefreshToken: refreshToken,
		Fingerprint:  fingerpring,
		ExpiresIn:    time.Now().Add(time.Duration(6 * 10e10)).Unix(),
		CreatedAt:    time.Now().Unix(),
	}
	s.sessions[refreshToken] = session
	return session, nil
}

func (s *SessionRepository) VerifyRefreshToken(cxt context.Context, fingerPrint, refreshToken string) (*model.Session, error) {
	for i := range s.sessions {
		if s.sessions[i].Fingerprint == fingerPrint && s.sessions[i].RefreshToken == refreshToken {
			newSession := &model.Session{
				Email:        s.sessions[i].Email,
				RefreshToken: refreshToken,
				Fingerprint:  fingerPrint,
				ExpiresIn:    time.Now().Add(time.Duration(6 * 10e10)).Unix(),
				CreatedAt:    time.Now().Unix(),
			}
			s.sessions[newSession.Email] = newSession
			return newSession, nil
		} else {
			return &model.Session{}, errors.New("invalid refresh token")
		}
	}
	return &model.Session{}, nil
}

func (s *SessionRepository) DeleteSession(cxt context.Context, fingerPrint, refreshToken string) error {
	_, ok := s.sessions[refreshToken]
	if !ok {
		return errors.New("session not found")
	}
	delete(s.sessions, refreshToken)
	_, ok = s.sessions[refreshToken]
	if ok {
		return errors.New("deleting error")
	}
	return nil
}

func (s *SessionRepository) DeleteAllSession(ctxb context.Context, email string) error {
	delete(s.sessions, email)
	_, ok := s.sessions[email]
	if ok {
		errors.New("error deleting records")
	}
	return nil
}
