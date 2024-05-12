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

func ValidPassword(s string) bool {
	return len(s) > 8 && len(s) < 100
}

func (u *User) Validate() error {
	if len(u.Email) == 0 && !ValidEmail(u.Email) {
		return errors.New("invalid Email")
	}
	if !ValidPassword(u.Password) {
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
