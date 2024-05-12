package SQLstore_test

import (
	"testing"

	"github.com/Quantum-calculators/MSU_UserService/internal/model"
	"github.com/Quantum-calculators/MSU_UserService/internal/store/SQLstore"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := SQLstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := SQLstore.New(db)
	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, teardown := SQLstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := SQLstore.New(db)
	email := "testuser@test.com"
	_, err := s.User().FindByEmail(email)
	assert.Error(t, err)

	u := model.TestUser(t)
	u.Email = email
	s.User().Create(u)

	u, err = s.User().FindByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_UpdateEmail(t *testing.T) {
	db, teardown := SQLstore.TestDB(t, databaseURL)
	defer teardown("users")
	s := SQLstore.New(db)
	u := model.TestUser(t)
	s.User().Create(u)

	newEmail := "newEmail@test.com"
	err1 := s.User().UpdateEmail(newEmail, u)
	assert.NoError(t, err1)

	newEmailIncorrerct := "incorrectEmail"
	err2 := s.User().UpdateEmail(newEmailIncorrerct, u)
	assert.Error(t, err2)
}

func TestUserRepository_UpdatePassword(t *testing.T) {
	db, teardown := SQLstore.TestDB(t, databaseURL)
	defer teardown("users")
	s := SQLstore.New(db)
	u := model.TestUser(t)
	s.User().Create(u)

	newPassword := "CorrectPass"
	err1 := s.User().UpdatePassword(newPassword, u)
	assert.NoError(t, err1)

	newPasswordIncor := "len<8"
	err2 := s.User().UpdatePassword(newPasswordIncor, u)
	assert.Error(t, err2)
}
