package teststore

import (
	"github.com/Quantum-calculators/MSU_UserService/internal/model"
	"github.com/Quantum-calculators/MSU_UserService/internal/store"
)

type UserRepository struct {
	store *Store
	users map[string]*model.User
}

func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err // The Validate function can return two errors: 'invalid password' or 'invalid email'
	}

	if err := u.BeforeCreate(); err != nil {
		return model.ErrEncryptedPassword
	}
	_, ok := r.users[u.Email]
	if ok {
		return store.ErrExistUserWithEmail
	}
	r.users[u.Email] = u
	u.ID = len(r.users)
	return nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u, ok := r.users[email]
	if !ok {
		return nil, store.ErrRecordNotFound
	}
	return u, nil
}

func (r *UserRepository) UpdateEmail(newEmail string, u *model.User) error {
	if !model.ValidEmail(newEmail) {
		return model.ErrInvalidEmail
	}
	newU := r.users[u.Email]
	newU.Email = newEmail
	r.users[u.Email] = newU
	u = newU
	return nil
}

func (r *UserRepository) UpdatePassword(password string, u *model.User) error {
	if !model.ValidPassword(password) {
		return model.ErrInvalidPass
	}
	u.Password = password
	if err := u.BeforeCreate(); err != nil {
		return model.ErrEncryptedPassword
	}
	return nil
}

func (r *UserRepository) GetUserByID(UserID int) (*model.User, error) {
	for _, j := range r.users {
		if j.ID == UserID {
			return j, nil
		}
	}
	return &model.User{}, store.ErrRecordNotFound
}
