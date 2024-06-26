package SQLstore

import (
	"context"
	"database/sql"
	"strings"

	"github.com/Quantum-calculators/MSU_UserService/internal/model"
	"github.com/Quantum-calculators/MSU_UserService/internal/store"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(ctxb context.Context, u *model.User) error {
	ctx, cancel := context.WithTimeout(ctxb, r.store.QueryTimeout)
	defer cancel()
	if err := u.Validate(); err != nil {
		return err
	}
	if err := u.BeforeCreate(); err != nil {
		return model.ErrEncryptedPassword
	}
	err := r.store.db.QueryRowContext(ctx,
		"INSERT INTO users (email, encrypted_password, verification_token) VALUES ($1, $2, $3) RETURNING id",
		u.Email,
		u.EncryptedPassword,
		u.VerificationToken,
	).Scan(&u.ID)
	switch {
	case err == nil:
		return nil
	case strings.Contains(err.Error(), "duplicate key value violates unique constraint"):
		return store.ErrExistUserWithEmail
	case err != nil:
		return store.ErrUnidentified
	}
	return nil
}

func (r *UserRepository) FindByEmail(ctxb context.Context, email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctxb, r.store.QueryTimeout)
	defer cancel()
	u := &model.User{}
	if err := r.store.db.QueryRowContext(ctx,
		"SELECT id, email, encrypted_password, verify FROM users WHERE email = $1", email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
		&u.Verified,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) UpdateEmail(ctxb context.Context, newEmail string, u *model.User) error {
	ctx, cancel := context.WithTimeout(ctxb, r.store.QueryTimeout)
	defer cancel()
	if !model.ValidEmail(newEmail) {
		return model.ErrInvalidEmail
	}
	if err := r.store.db.QueryRowContext(ctx, "UPDATE users SET email = $1 WHERE email = $2", newEmail, u.Email).Err(); err != nil {
		return store.ErrUpdateEmailFailed
	}
	u.Email = newEmail
	return nil
}

func (r *UserRepository) UpdatePassword(ctxb context.Context, password string, u *model.User) error {
	ctx, cancel := context.WithTimeout(ctxb, r.store.QueryTimeout)
	defer cancel()
	if !model.ValidPassword(password) {
		return model.ErrInvalidPass
	}
	u.Password = password
	if err := u.BeforeCreate(); err != nil {
		return err
	}
	if err := r.store.db.QueryRowContext(ctx,
		"UPDATE users SET encrypted_password = $1 WHERE email = $2",
		u.EncryptedPassword,
		u.Email,
	).Err(); err != nil {
		return store.ErrUpdatePassFailed
	}
	return nil
}

func (r *UserRepository) GetUserByID(ctxb context.Context, UserID int) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctxb, r.store.QueryTimeout)
	defer cancel()
	u := model.User{}
	if err := r.store.db.QueryRowContext(ctx,
		"SELECT email, encrypted_password, verify FROM users WHERE id = $1",
		UserID,
	).Scan(
		&u.Email,
		&u.EncryptedPassword,
		&u.Verified,
	); err != nil {
		return &model.User{}, store.ErrRecordNotFound
	}
	return &u, nil
}

func (r *UserRepository) SetVerify(ctxb context.Context, Email string, verify bool) error {
	ctx, cancel := context.WithTimeout(ctxb, r.store.QueryTimeout)
	defer cancel()
	if err := r.store.db.QueryRowContext(ctx,
		"UPDATE users SET verify = $1 WHERE email = $2;",
		verify,
		Email,
	).Err(); err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) UpdateVerificationToken(ctxb context.Context, Email, token string) error {
	ctx, cancel := context.WithTimeout(ctxb, r.store.QueryTimeout)
	defer cancel()
	err := r.store.db.QueryRowContext(ctx,
		"UPDATE users SET verification_token = $1 WHERE email = $2;",
		token,
		Email,
	).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) CheckVerificationToken(ctxb context.Context, Email, token string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctxb, r.store.QueryTimeout)
	defer cancel()
	var dbToken string
	err := r.store.db.QueryRowContext(ctx,
		"SELECT verification_token FROM users WHERE email = $1",
		Email,
	).Scan(&dbToken)
	if err != nil {
		return false, err
	}
	if dbToken == token {
		return true, nil
	}
	return false, nil
}
