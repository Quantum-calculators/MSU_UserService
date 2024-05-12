package store

import "github.com/Quantum-calculators/MSU_UserService/internal/model"

type UserRepository interface {
	Create(*model.User) error
	FindByEmail(string) (*model.User, error)
	UpdateEmail(string, *model.User) error
	UpdatePassword(string, *model.User) error
}

type CacheRepository interface {
	Set() (string, error)
	Get() (string, error)
}
