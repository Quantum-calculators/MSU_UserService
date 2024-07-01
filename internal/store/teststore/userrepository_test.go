package testStore_test

import (
	"testing"

	"github.com/Quantum-calculators/MSU_UserService/internal/model"
	"github.com/Quantum-calculators/MSU_UserService/internal/store"
	"github.com/Quantum-calculators/MSU_UserService/internal/store/testStore"
	token_generator "github.com/Quantum-calculators/MSU_UserService/internal/tokenGenerator"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	s := testStore.New()

	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(nil, u))
	assert.NotNil(t, u)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	s := testStore.New()
	email := "testuser@test.com"
	_, err := s.User().FindByEmail(nil, email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	u := model.TestUser(t)
	u.Email = email
	s.User().Create(nil, u)

	u, err = s.User().FindByEmail(nil, email)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_UpdateEmail(t *testing.T) {
	s := testStore.New()
	u := model.TestUser(t)
	newEmail := "newemail@test.com"
	err := s.User().UpdatePassword(nil, newEmail, u)
	assert.NoError(t, err)

	u2 := model.TestUser(t)
	newEmailIncorrect := "IncorrectEmail"
	err = s.User().UpdateEmail(nil, u2.Email, newEmailIncorrect)

	assert.Error(t, err)
}

func TestUserRepository_UpdatePassword(t *testing.T) {
	s := testStore.New()
	u := model.TestUser(t)
	newPas := "testPass12"
	err := s.User().UpdatePassword(nil, newPas, u)
	assert.NoError(t, err)

	newPasIncorrect := "incor" // len < 8
	err = s.User().UpdatePassword(nil, newPasIncorrect, u)
	assert.Error(t, err)
}

func TestUserRepository_GetUserByID(t *testing.T) {
	s := testStore.New()
	u := model.TestUser(t)
	s.User().Create(nil, u)

	_, err1 := s.User().GetUserByID(nil, u.ID)
	assert.NoError(t, err1)

	_, err2 := s.User().GetUserByID(nil, 2)
	assert.Error(t, err2)
}

func TestUserRepository_SetVerify(t *testing.T) {
	s := testStore.New()
	u := model.TestUser(t)
	s.User().Create(nil, u)

	NotVerifiedU, err1 := s.User().GetUserByID(nil, u.ID)
	assert.NoError(t, err1)
	assert.False(t, NotVerifiedU.Verified)

	err2 := s.User().SetVerify(nil, u.Email, true)
	assert.NoError(t, err2)

	VerifiedU, err3 := s.User().GetUserByID(nil, u.ID)
	assert.NoError(t, err3)
	assert.True(t, VerifiedU.Verified)
}

func TestUserRepository_CheckVerificationToken(t *testing.T) {
	s := testStore.New()
	u := model.TestUser(t)
	s.User().Create(nil, u)

	pass, err := s.User().CheckVerificationToken(nil, u.Email, u.VerificationToken)
	assert.NoError(t, err)
	assert.True(t, pass)

	pass1, err := s.User().CheckVerificationToken(nil, u.Email, "not_valid_token")
	assert.NoError(t, err)
	assert.False(t, pass1)
}

func TestUserRepository_UpdateVerificationToken(t *testing.T) {
	s := testStore.New()
	u := model.TestUser(t)
	err := s.User().Create(nil, u)
	assert.NoError(t, err)

	pass, err := s.User().CheckVerificationToken(nil, u.Email, u.VerificationToken)
	assert.NoError(t, err)
	assert.True(t, pass)

	newVerToken, err := token_generator.GenerateRandomString(64)
	assert.NoError(t, err)

	err = s.User().UpdateVerificationToken(nil, u.Email, newVerToken)
	assert.NoError(t, err)

	pass1, err := s.User().CheckVerificationToken(nil, u.Email, newVerToken)
	assert.NoError(t, err)
	assert.True(t, pass1)
}
