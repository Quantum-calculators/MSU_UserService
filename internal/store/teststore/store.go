package testStore

import (
	"github.com/Quantum-calculators/MSU_UserService/internal/model"
	"github.com/Quantum-calculators/MSU_UserService/internal/store"
)

type Store struct {
	userRepository    *UserRepository
	sessionRepository *SessionRepository
	ExpRefreshToken   int
}

func New() *Store {
	return &Store{}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
		users: make(map[string]*model.User),
	}

	return s.userRepository
}

func (s *Store) Session() store.SessionRepository {
	if s.sessionRepository != nil {
		return s.sessionRepository
	}

	s.sessionRepository = &SessionRepository{
		store:    s,
		sessions: make(map[string]*model.Session),
	}

	return s.sessionRepository
}
