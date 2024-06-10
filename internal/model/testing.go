package model

import (
	"testing"
	"time"
)

func TestUser(t *testing.T) *User {
	return &User{
		Email:             "user@example.org",
		Password:          "examplePassword",
		Verified:          false,
		VerificationToken: "34mq0pcxum3q048xjt438pxjgmp30jx89mp30483jcxqm3hpcgx89h3mcg93q04x",
	}
}

func TestSession(t *testing.T) *Session {
	return &Session{
		Fingerprint: "ru-RU.Chromium.macOS.Mozilla/5.0",
		ExpiresIn:   1516239022,
		CreatedAt:   time.Now().Unix(),
	}
}
