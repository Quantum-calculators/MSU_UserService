package SQLstore

import (
	"database/sql"
	"time"

	"github.com/Quantum-calculators/MSU_UserService/internal/store"
	_ "github.com/lib/pq"
)

type Store struct {
	db                *sql.DB
	userRepository    *UserRepository
	sessionRepository *SessionRepository
	ExpRefreshToken   int
	QueryTimeout      time.Duration
}

func New(db *sql.DB, expRefreshToken int, QueryTimeout time.Duration) *Store {
	return &Store{
		db:              db,
		ExpRefreshToken: expRefreshToken,
		QueryTimeout:    QueryTimeout,
	}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}

func (s *Store) Session() store.SessionRepository {
	if s.sessionRepository != nil {
		return s.sessionRepository
	}

	s.sessionRepository = &SessionRepository{
		store: s,
	}

	return s.sessionRepository
}
