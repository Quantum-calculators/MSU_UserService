package SQLstore_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Quantum-calculators/MSU_UserService/internal/model"
	"github.com/Quantum-calculators/MSU_UserService/internal/store/SQLstore"
	"github.com/stretchr/testify/assert"
)

func TestSessionRepository_CreateSession(t *testing.T) {
	db, teardown := SQLstore.TestDB(t, databaseURL)
	defer teardown("users", "sessions")
	s := SQLstore.New(db, 5, time.Millisecond*100)
	ctxb := context.Background()

	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(ctxb, u))
	session := model.TestSession(t)
	fmt.Println(u)
	_, err := s.Session().CreateSession(ctxb, u.Email, session.Fingerprint)
	assert.NoError(t, err)
}

func TestSessionRepository_VerifyRefreshToken(t *testing.T) {
	db, teardown := SQLstore.TestDB(t, databaseURL)
	defer teardown("users", "sessions")
	s := SQLstore.New(db, 5, time.Millisecond*100)
	ctxb := context.Background()

	session := model.TestSession(t)
	u := model.TestUser(t)
	fmt.Println(u)
	assert.NoError(t, s.User().Create(ctxb, u))
	session, err := s.Session().CreateSession(ctxb, u.Email, session.Fingerprint)
	assert.NoError(t, err)

	_, err1 := s.Session().VerifyRefreshToken(ctxb, "invalidFingerprint", "invalidRefreshToken")
	assert.Error(t, err1)

	_, err2 := s.Session().VerifyRefreshToken(ctxb, "invalidFingerprint", "invalidRefreshToken") // ??
	assert.Error(t, err2)

	Fingerprint := "ru-RU.Chromium.macOS.Mozilla/5.0"

	_, err3 := s.Session().VerifyRefreshToken(ctxb, Fingerprint, session.RefreshToken)
	assert.NoError(t, err3)
}

func TestSessionRepository_DeleteSession(t *testing.T) {
	db, teardown := SQLstore.TestDB(t, databaseURL)
	defer teardown("users", "sessions")
	s := SQLstore.New(db, 5, time.Millisecond*100)
	ctxb := context.Background()

	session := model.TestSession(t)
	u := model.TestUser(t)
	err := s.User().Create(ctxb, u)
	assert.NoError(t, err)
	session, err2 := s.Session().CreateSession(ctxb, u.Email, session.Fingerprint)
	assert.NoError(t, err2)

	err1 := s.Session().DeleteSession(ctxb, session.Fingerprint, session.RefreshToken)
	assert.NoError(t, err1)
}

func TestSessionRepository_DeleteAllSession(t *testing.T) {
	db, teardown := SQLstore.TestDB(t, databaseURL)
	defer teardown("users", "sessions")
	s := SQLstore.New(db, 5, time.Millisecond*100)
	ctxb := context.Background()

	session := model.TestSession(t)
	u := model.TestUser(t)
	err := s.User().Create(ctxb, u)
	assert.NoError(t, err)
	session, err2 := s.Session().CreateSession(ctxb, u.Email, session.Fingerprint)
	assert.NoError(t, err2)
	_, err3 := s.Session().CreateSession(ctxb, u.Email, "Another fingerprint")
	assert.NoError(t, err3)

	err1 := s.Session().DeleteAllSession(ctxb, session.Email)
	assert.NoError(t, err1)
}
