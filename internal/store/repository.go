package store

import "github.com/Quantum-calculators/MSU_UserService/internal/model"

type UserRepository interface {
	Create(*model.User) error
	FindByEmail(string) (*model.User, error)
}

type JWTRepository interface {
	CreateAccessToken() (string, error)
	CreateRefeshToken() (string, error)
	UpdateAccessToken() (string, error)
	UpdateRefreshToken() (string, error)
}
