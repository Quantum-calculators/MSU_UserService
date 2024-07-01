package model

type Session struct {
	ID           int64
	Email        string
	RefreshToken string
	Fingerprint  string
	ExpiresIn    int64
	CreatedAt    int64
}
