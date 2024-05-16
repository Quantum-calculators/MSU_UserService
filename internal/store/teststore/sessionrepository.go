package teststore

import (
	"errors"
	"fmt"
	"time"

	"github.com/Quantum-calculators/MSU_UserService/internal/model"
	"github.com/Quantum-calculators/MSU_UserService/internal/token_generator"
)

type SessionRepository struct {
	store    *Store
	sessions map[string]*model.Session
}

func (s *SessionRepository) CreateSession(userId uint32, fingerpring string) (*model.Session, error) {
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

func (s *SessionRepository) VerifyRefreshToken(fingerPrint, refreshToken string) (*model.Session, error) {
	fmt.Println(refreshToken)
	session, ok := s.sessions[refreshToken]
	if ok && session.Fingerprint == fingerPrint && session.RefreshToken == refreshToken {
		newSession, err := s.CreateSession(uint32(session.UserId), fingerPrint)
		if err != nil {
			return &model.Session{}, err
		}
		return newSession, nil
	} else {
		return &model.Session{}, errors.New("invalid refresh token")
	}
}
