package SQLstore

import (
	"database/sql"
	"errors"

	"github.com/Quantum-calculators/MSU_UserService/internal/model"
	"github.com/Quantum-calculators/MSU_UserService/internal/store"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}
	return r.store.db.QueryRow(
		"INSERT INTO users (email, encrypted_password) VALUES ($1, $2) RETURNING id",
		u.Email,
		u.EncryptedPassword,
	).Scan(&u.ID)
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, encrypted_password FROM users WHERE email = $1", email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) UpdateEmail(newEmail string, u *model.User) error {
	if !model.ValidEmail(newEmail) {
		return errors.New("not valid email")
	}
	if err := r.store.db.QueryRow("UPDATE users SET email = $1 WHERE email = $2", newEmail, u.Email).Err(); err != nil {
		return err
	}
	u.Email = newEmail
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
	if err := r.store.db.QueryRow(
		"UPDATE users SET encrypted_password = $1 WHERE email = $2",
		u.EncryptedPassword,
		u.Email,
	).Err(); err != nil {
		return err
	}
	return nil
}

// func (r *UserRepository) SetRefreshToken(refreshToken string, expRefreshToken int, u *model.User) error {
// 	if err := r.store.db.QueryRow(
// 		"UPDATE users SET refresh_token = $1, exp_refresh_token = $2 WHERE email = $3",
// 		refreshToken,
// 		expRefreshToken,
// 		u.Email,
// 	).Err(); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (r *UserRepository) GetUserByRefreshToken(UserRefreshToken string) (*model.User, error) {
// 	u := &model.User{}
// 	if err := r.store.db.QueryRow(
// 		"SELECT id, email FROM users WHERE refresh_token = $1",
// 		UserRefreshToken,
// 	).Scan(
// 		&u.ID,
// 		&u.Email,
// 	); err != nil {
// 		return nil, err
// 	}
// 	return u, nil
// }
