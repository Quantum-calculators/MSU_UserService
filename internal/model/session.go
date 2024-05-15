package model

type Session struct {
	ID           uint32
	UserId       uint32
	RefreshToken string
	Fingerprint  string
	ExpiresIn    int64
	CreatedAt    int64
}
