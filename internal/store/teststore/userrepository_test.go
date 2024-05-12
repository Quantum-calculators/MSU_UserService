package teststore_test

import (
	"testing"

	"github.com/Quantum-calculators/MSU_UserService/internal/model"
	"github.com/Quantum-calculators/MSU_UserService/internal/store"
	"github.com/Quantum-calculators/MSU_UserService/internal/store/teststore"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	s := teststore.New()
	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	s := teststore.New()
	email := "testuser@test.com"
	_, err := s.User().FindByEmail(email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	u := model.TestUser(t)
	u.Email = email
	s.User().Create(u)

	u, err = s.User().FindByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_UpdateEmail(t *testing.T) {
	s := teststore.New()
	u := model.TestUser(t)
	newEmail := "newemail@test.com"
	err := s.User().UpdatePassword(newEmail, u)
	assert.NoError(t, err)

	u2 := model.TestUser(t)
	newEmailIncorrect := "IncorrectEmail"
	err = s.User().UpdateEmail(newEmailIncorrect, u2)

	assert.Error(t, err)
}

func TestUserRepository_UpdatePassword(t *testing.T) {
	s := teststore.New()
	u := model.TestUser(t)
	newPas := "testPass12"
	err := s.User().UpdatePassword(newPas, u)
	assert.NoError(t, err)

	newPasIncorrect := "incor" // len < 8
	err = s.User().UpdatePassword(newPasIncorrect, u)
	assert.Error(t, err)
}
