package SQLstore_test

import (
	"fmt"
	"testing"

	"github.com/Quantum-calculators/MSU_UserService/internal/model"
	"github.com/Quantum-calculators/MSU_UserService/internal/store/SQLstore"
	"github.com/stretchr/testify/assert"
)

func TestSessionRepository_CreateSession(t *testing.T) {
	db, teardown := SQLstore.TestDB(t, databaseURL)
	defer teardown("users", "sessions")
	s := SQLstore.New(db, 5)

	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(u))
	session := model.TestSession(t)
	fmt.Println(u)
	_, err := s.Session().CreateSession(uint32(u.ID), session.Fingerprint)
	assert.NoError(t, err)
}

func TestSessionRepository_VerifyRefreshToken(t *testing.T) {
	db, teardown := SQLstore.TestDB(t, databaseURL)
	defer teardown("users", "sessions")
	s := SQLstore.New(db, 5)

	session := model.TestSession(t)
	u := model.TestUser(t)
	fmt.Println(u)
	assert.NoError(t, s.User().Create(u))
	session, err := s.Session().CreateSession(uint32(u.ID), session.Fingerprint)
	assert.NoError(t, err)

	_, err1 := s.Session().VerifyRefreshToken("invalidFingerprint", "invalidRefreshToken")
	assert.Error(t, err1)

	_, err2 := s.Session().VerifyRefreshToken("invalidFingerprint", "invalidRefreshToken")
	assert.Error(t, err2)

	Fingerprint := "ru-RU.Chromium.macOS.Mozilla/5.0"

	_, err3 := s.Session().VerifyRefreshToken(Fingerprint, session.RefreshToken)
	assert.NoError(t, err3)
}
