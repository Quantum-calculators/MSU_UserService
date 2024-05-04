package model

import (
	"errors"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                int
	Email             string
	Password          string
	EncryptedPassword string
}

func ValidEmail(s string) bool {
	_, err := mail.ParseAddress(s)
	return err == nil
}

func (u *User) Validate() error {
	if len(u.Email) == 0 || !ValidEmail(u.Email) {
		return errors.New("invalid Email")
	}
	if len(u.Password) < 8 || len(u.Password) > 100 {
		if len(u.EncryptedPassword) == 0 {
			return errors.New("invalid Password")
		}
	}
	return nil
}

func (u *User) BeforeCreate() error {
	if len(u.Password) != 0 {
		enc, err := encryptString(u.Password)
		if err != nil {
			return err
		}
		u.EncryptedPassword = enc
	}
	return nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
