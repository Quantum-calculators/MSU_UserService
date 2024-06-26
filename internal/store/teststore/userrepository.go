package testStore

import (
	"context"
	"errors"

	"github.com/Quantum-calculators/MSU_UserService/internal/model"
	"github.com/Quantum-calculators/MSU_UserService/internal/store"
)

type UserRepository struct {
	store *Store
	users map[string]*model.User
}

func (r *UserRepository) Create(ctxb context.Context, u *model.User) error {
	if err := u.Validate(); err != nil {
		return err // The Validate function can return two errors: 'invalid password' or 'invalid email'
	}

	if err := u.BeforeCreate(); err != nil {
		return model.ErrEncryptedPassword
	}
	u.Verified = false
	_, ok := r.users[u.Email]
	if ok {
		return store.ErrExistUserWithEmail
	}
	r.users[u.Email] = u
	u.ID = len(r.users)
	return nil
}

func (r *UserRepository) FindByEmail(ctxb context.Context, email string) (*model.User, error) {
	u, ok := r.users[email]
	if !ok {
		return nil, store.ErrRecordNotFound
	}
	return u, nil
}

func (r *UserRepository) UpdateEmail(ctxb context.Context, newEmail string, u *model.User) error {
	if !model.ValidEmail(newEmail) {
		return model.ErrInvalidEmail
	}
	newU := r.users[u.Email]
	newU.Email = newEmail
	r.users[u.Email] = newU
	u = newU
	return nil
}

func (r *UserRepository) UpdatePassword(ctxb context.Context, password string, u *model.User) error {
	if !model.ValidPassword(password) {
		return model.ErrInvalidPass
	}
	u.Password = password
	if err := u.BeforeCreate(); err != nil {
		return model.ErrEncryptedPassword
	}
	return nil
}

func (r *UserRepository) GetUserByID(ctxb context.Context, UserID int) (*model.User, error) {
	for _, j := range r.users {
		if j.ID == UserID {
			return j, nil
		}
	}
	return &model.User{}, store.ErrRecordNotFound
}

func (r *UserRepository) SetVerify(ctxb context.Context, Email string, verify bool) error {
	_, ok := r.users[Email]
	if !ok {
		return errors.New("user not found")
	}
	r.users[Email].Verified = verify
	return nil
}

func (r *UserRepository) CheckVerificationToken(ctxb context.Context, Email, token string) (bool, error) {
	_, ok := r.users[Email]
	if !ok {
		return false, errors.New("user not found")
	}

	if r.users[Email].VerificationToken == token {
		return true, nil
	}
	return false, nil
}

func (r *UserRepository) UpdateVerificationToken(ctxb context.Context, Email, token string) error {
	user, ok := r.users[Email]
	if !ok {
		return errors.New("user not found")
	}
	user.VerificationToken = token
	return nil
}
