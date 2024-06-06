package model

import "errors"

var (
	ErrInvalidPass       = errors.New("model: invalid Password")
	ErrInvalidEmail      = errors.New("model: invalid Email")
	ErrEncryptedPassword = errors.New("model: the password could not be encrypted")
)
