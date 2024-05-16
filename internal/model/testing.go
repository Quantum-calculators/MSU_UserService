package model

import (
	"testing"
	"time"
)

func TestUser(t *testing.T) *User {
	return &User{
		Email:    "user@example.org",
		Password: "examplePassword",
	}
}

func TestSession(t *testing.T) *Session {
	return &Session{
		Fingerprint: "ru-RU.Chromium.macOS.Mozilla/5.0",
		ExpiresIn:   1516239022,
		CreatedAt:   time.Now().Unix(),
	}
}
