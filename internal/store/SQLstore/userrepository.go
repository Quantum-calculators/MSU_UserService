package SQLstore

import (
	"database/sql"

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
		return model.ErrEncryptedPassword
	}
	if err := r.store.db.QueryRow(
		"INSERT INTO users (email, encrypted_password) VALUES ($1, $2) RETURNING id",
		u.Email,
		u.EncryptedPassword,
	).Scan(&u.ID); err != nil {
		return store.ErrExistUserWithEmail
	}
	return nil
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
		return model.ErrInvalidEmail
	}
	if err := r.store.db.QueryRow("UPDATE users SET email = $1 WHERE email = $2", newEmail, u.Email).Err(); err != nil {
		return store.ErrUpdateEmailFailed
	}
	u.Email = newEmail
	return nil
}

func (r *UserRepository) UpdatePassword(password string, u *model.User) error {
	if !model.ValidPassword(password) {
		return model.ErrInvalidPass
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
		return store.ErrUpdatePassFailed
	}
	return nil
}

func (r *UserRepository) GetUserByID(UserID int) (*model.User, error) {
	u := model.User{}
	if err := r.store.db.QueryRow(
		"SELECT email, encrypted_password FROM users WHERE id = $1",
		UserID,
	).Scan(
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		return &model.User{}, store.ErrRecordNotFound
	}
	return &u, nil
}
