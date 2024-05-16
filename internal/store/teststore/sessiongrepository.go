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
	sessions map[uint32]*model.Session
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
	fmt.Println(session)
	s.sessions[userId] = session
	return session, nil
}

func (s *SessionRepository) VerifyRefreshToken(userID int, fingerPrint, refreshToken string) (string, error) {
	session := s.sessions[uint32(userID)]
	newRefreshToken, err := token_generator.GenerateRandomString(128)
	if err != nil {
		return "", err
	}
	if session.Fingerprint == fingerPrint && session.RefreshToken == refreshToken {
		return newRefreshToken, nil
	} else {
		return "", errors.New("invalid refresh token")
	}
}
