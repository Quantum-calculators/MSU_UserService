package model

import (
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                int    `json:"id"`
	Email             string `json:"email"`
	Password          string `json:"password,omitempty"`
	EncryptedPassword string `json:"-"`
	Verified          bool   `json:"verified,omitempty"`
	VerificationToken string `json:"verification_token,omitempty"`
}

func ValidEmail(s string) bool {
	_, err := mail.ParseAddress(s)
	return err == nil
}

func ValidPassword(s string) bool {
	return len(s) > 8 && len(s) < 72
}

func (u *User) Validate() error {
	if len(u.Email) == 0 || !ValidEmail(u.Email) {
		return ErrInvalidEmail
	}
	if !ValidPassword(u.Password) {
		if len(u.EncryptedPassword) == 0 {
			return ErrInvalidPass
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

// Returns False if the passwords do not match
func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}

func (u *User) Sanitize() {
	u.Password = ""
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
