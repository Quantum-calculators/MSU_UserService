package SQLstore_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Quantum-calculators/MSU_UserService/internal/model"
	"github.com/Quantum-calculators/MSU_UserService/internal/store/SQLstore"
	token_generator "github.com/Quantum-calculators/MSU_UserService/internal/tokenGenerator"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := SQLstore.TestDB(t, databaseURL)
	defer teardown("users")
	s := SQLstore.New(db, 5, time.Millisecond*100)
	ctxb := context.Background()

	u := model.TestUser(t)
	err := s.User().Create(ctxb, u)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, teardown := SQLstore.TestDB(t, databaseURL)
	defer teardown("users")
	s := SQLstore.New(db, 5, time.Millisecond*100)
	ctxb := context.Background()

	email := "testuser@test.com"
	_, err := s.User().FindByEmail(ctxb, email)
	assert.Error(t, err)

	u := model.TestUser(t)
	u.Email = email
	s.User().Create(ctxb, u)

	u, err = s.User().FindByEmail(ctxb, email)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_UpdateEmail(t *testing.T) {
	db, teardown := SQLstore.TestDB(t, databaseURL)
	defer teardown("users")
	s := SQLstore.New(db, 5, time.Millisecond*100)
	ctxb := context.Background()

	u := model.TestUser(t)
	s.User().Create(ctxb, u)

	newEmail := "newEmail@test.com"
	err1 := s.User().UpdateEmail(ctxb, newEmail, u)
	assert.NoError(t, err1)

	newEmailIncorrerct := "incorrectEmail"
	err2 := s.User().UpdateEmail(ctxb, newEmailIncorrerct, u)
	assert.Error(t, err2)
}

func TestUserRepository_UpdatePassword(t *testing.T) {
	db, teardown := SQLstore.TestDB(t, databaseURL)
	defer teardown("users")
	s := SQLstore.New(db, 5, time.Millisecond*100)
	ctxb := context.Background()

	u := model.TestUser(t)
	s.User().Create(ctxb, u)

	newPassword := "CorrectPass"
	err1 := s.User().UpdatePassword(ctxb, newPassword, u)
	assert.NoError(t, err1)

	newPasswordIncor := "len<8"
	err2 := s.User().UpdatePassword(ctxb, newPasswordIncor, u)
	assert.Error(t, err2)
}

func TestUserRepository_GetUserByID(t *testing.T) {
	db, teardown := SQLstore.TestDB(t, databaseURL)
	defer teardown("users")
	s := SQLstore.New(db, 5, time.Millisecond*100)
	ctxb := context.Background()

	u := model.TestUser(t)
	err := s.User().Create(ctxb, u)
	assert.NoError(t, err)

	fmt.Println(u.ID)

	_, err1 := s.User().GetUserByID(ctxb, u.ID)
	assert.NoError(t, err1)

	_, err2 := s.User().GetUserByID(ctxb, 2)
	assert.Error(t, err2)
}

func TestUserRepository_SetVerify(t *testing.T) {
	db, teardown := SQLstore.TestDB(t, databaseURL)
	defer teardown("users")
	s := SQLstore.New(db, 5, time.Millisecond*100)
	ctxb := context.Background()

	u := model.TestUser(t)
	err := s.User().Create(ctxb, u)
	assert.NoError(t, err)

	NotVerifiedU, err1 := s.User().GetUserByID(ctxb, u.ID)
	assert.NoError(t, err1)
	assert.False(t, NotVerifiedU.Verified)

	err2 := s.User().SetVerify(ctxb, u.Email, true)
	assert.NoError(t, err2)

	VerifiedU, err3 := s.User().GetUserByID(ctxb, u.ID)
	assert.NoError(t, err3)
	assert.True(t, VerifiedU.Verified)
}

func TestUserRepository_CheckVerificationToken(t *testing.T) {
	db, teardown := SQLstore.TestDB(t, databaseURL)
	defer teardown("users")
	s := SQLstore.New(db, 5, time.Millisecond*100)
	ctxb := context.Background()

	u := model.TestUser(t)
	err := s.User().Create(ctxb, u)
	assert.NoError(t, err)

	pass, err := s.User().CheckVerificationToken(ctxb, u.Email, u.VerificationToken)
	assert.NoError(t, err)
	assert.True(t, pass)

	pass1, err := s.User().CheckVerificationToken(ctxb, u.Email, "not_valid_token")
	assert.NoError(t, err)
	assert.False(t, pass1)
}

func TestUserRepository_UpdateVerificationToken(t *testing.T) {
	db, teardown := SQLstore.TestDB(t, databaseURL)
	defer teardown("users")
	s := SQLstore.New(db, 5, time.Millisecond*100)
	ctxb := context.Background()

	u := model.TestUser(t)
	err := s.User().Create(ctxb, u)
	assert.NoError(t, err)

	pass, err := s.User().CheckVerificationToken(ctxb, u.Email, u.VerificationToken)
	assert.NoError(t, err)
	assert.True(t, pass)

	newVerToken, err := token_generator.GenerateRandomString(64)

	err = s.User().UpdateVerificationToken(ctxb, u.Email, newVerToken)
	assert.NoError(t, err)

	pass1, err := s.User().CheckVerificationToken(ctxb, u.Email, newVerToken)
	assert.NoError(t, err)
	assert.True(t, pass1)
}
