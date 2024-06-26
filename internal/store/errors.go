package store

import "errors"

var (
	ErrRecordNotFound      = errors.New("store: record not found")
	ErrExistUserWithEmail  = errors.New("store: a user with this email already exists")
	ErrUpdateEmailFailed   = errors.New("store: failed to update email")
	ErrUpdatePassFailed    = errors.New("store: failed to update password")
	ErrRefreshTokenExpired = errors.New("store: refreshToken has expired")
	ErrTimeOut             = errors.New("store: query timeout")
)
