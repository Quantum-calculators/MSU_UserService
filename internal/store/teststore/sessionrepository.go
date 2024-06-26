package testStore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Quantum-calculators/MSU_UserService/internal/model"
	token_generator "github.com/Quantum-calculators/MSU_UserService/internal/tokenGenerator"
)

type SessionRepository struct {
	store    *Store
	sessions map[string]*model.Session
}

func (s *SessionRepository) CreateSession(cxt context.Context, userId uint32, fingerpring string) (*model.Session, error) {
	refreshToken, err := token_generator.GenerateRandomString(64)
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
	s.sessions[refreshToken] = session
	return session, nil
}

func (s *SessionRepository) VerifyRefreshToken(cxt context.Context, fingerPrint, refreshToken string) (*model.Session, error) {
	fmt.Println(refreshToken)
	session, ok := s.sessions[refreshToken]
	if ok && session.Fingerprint == fingerPrint && session.RefreshToken == refreshToken {
		newSession := &model.Session{
			UserId:       session.UserId,
			RefreshToken: refreshToken,
			Fingerprint:  fingerPrint,
			ExpiresIn:    time.Now().Add(time.Duration(6 * 10e10)).Unix(),
			CreatedAt:    time.Now().Unix(),
		}
		return newSession, nil
	} else {
		return &model.Session{}, errors.New("invalid refresh token")
	}
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
