package teststore

import (
	"errors"

	"github.com/Quantum-calculators/MSU_UserService/internal/model"
	"github.com/Quantum-calculators/MSU_UserService/internal/store"
)

type UserRepository struct {
	store *Store
	users map[string]*model.User
}

func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
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
		return errors.New("not valid email")
	}
	newU := r.users[u.Email]
	newU.Email = newEmail
	r.users[u.Email] = newU
	u = newU
	return nil
}

func (r *UserRepository) UpdatePassword(password string, u *model.User) error {
	if !model.ValidPassword(password) {
		return errors.New("not valid password")
	}
	u.Password = password
	if err := u.BeforeCreate(); err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetUserByID(UserID int) (*model.User, error) {
	for _, j := range r.users {
		if j.ID == UserID {
			return j, nil
		}
	}
	return &model.User{}, errors.New("user not found")
}

// func (r *UserRepository) SetRefreshToken(refreshToken string, expRefreshToken int, u *model.User) error {
// 	u.RefreshToken = refreshToken
// 	u.ExpRefreshToken = expRefreshToken
// 	return nil
// }

// func (r *UserRepository) GetUserByRefreshToken(UserRefreshToken string) (*model.User, error) {
// 	return nil, nil
// }
