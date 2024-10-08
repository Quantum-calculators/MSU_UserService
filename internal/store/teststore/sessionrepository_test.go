package testStore_test

import (
	"testing"

	"github.com/Quantum-calculators/MSU_UserService/internal/model"
	"github.com/Quantum-calculators/MSU_UserService/internal/store/testStore"
	"github.com/stretchr/testify/assert"
)

func TestSessionRepository_CreateSession(t *testing.T) {
	s := testStore.New()

	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(nil, u))
	session := model.TestSession(t)
	_, err := s.Session().CreateSession(nil, u.Email, session.Fingerprint)
	assert.NoError(t, err)
}

func TestSessionRepository_VerifyRefreshToken(t *testing.T) {
	s := testStore.New()

	session := model.TestSession(t)
	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(nil, u))

	session, err := s.Session().CreateSession(nil, u.Email, session.Fingerprint)
	assert.NoError(t, err)

	_, err1 := s.Session().VerifyRefreshToken(nil, "invalidFingerprint", "invalidRefreshToken")
	assert.Error(t, err1)

	_, err2 := s.Session().VerifyRefreshToken(nil, session.Fingerprint, session.RefreshToken)
	assert.NoError(t, err2)
}

func TestSessionRepository_DeleteSession(t *testing.T) {
	s := testStore.New()

	session := model.TestSession(t)
	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(nil, u))
	session, err := s.Session().CreateSession(nil, u.Email, session.Fingerprint)
	assert.NoError(t, err)

	err1 := s.Session().DeleteSession(nil, session.Fingerprint, session.RefreshToken)
	assert.NoError(t, err1)

	err2 := s.Session().DeleteSession(nil, session.Fingerprint, "invalidResreshToken")
	assert.Error(t, err2)
}

func TestSessionRepository_DeleteAllSession(t *testing.T) {
	s := testStore.New()

	session := model.TestSession(t)
	u := model.TestUser(t)
	err := s.User().Create(nil, u)
	assert.NoError(t, err)
	session, err2 := s.Session().CreateSession(nil, u.Email, session.Fingerprint)
	assert.NoError(t, err2)
	_, err3 := s.Session().CreateSession(nil, u.Email, "Another fingerprint")
	assert.NoError(t, err3)

	err1 := s.Session().DeleteAllSession(nil, session.Email)
	assert.NoError(t, err1)
}
