package store

import "github.com/Quantum-calculators/MSU_UserService/internal/model"

type UserRepository interface {
	// The function creates a record about the user in the database.
	//
	// 	Input params:
	// 		1. type User struct {
	// 			ID                int
	// 			Email             string
	// 			Password          string
	// 			RefreshToken      string
	// 			ExpRefreshToken   int
	// 			EncryptedPassword string
	// 		}
	// 	Output params:
	// 		1. error or nil
	Create(*model.User) error
	// The function searches for the user by his email.
	//
	// 	Input params:
	// 		1. email string
	// 	Output params:
	// 		1. type User struct {
	// 			ID                int
	// 			Email             string
	// 			Password          string
	// 			RefreshToken      string
	// 			ExpRefreshToken   int
	// 			EncryptedPassword string
	// 		}
	// 		2. error or nil
	FindByEmail(string) (*model.User, error)
	// Updates the user's email.
	// Accepts a new email and user model as input.
	//
	// 	Input params:
	//		1. email string
	// 	Output params:
	//		1. type User struct {
	// 			ID                int
	// 			Email             string
	// 			Password          string
	// 			RefreshToken      string
	// 			ExpRefreshToken   int
	// 			EncryptedPassword string
	// 		}
	// 		2. error or nil
	UpdateEmail(string, *model.User) error
	// Updates the user's password.
	// Accepts a new password and user model as input.
	//
	// 	Input params:
	//		1. UserID int
	// 		2. type User struct {
	// 			ID                int
	// 			Email             string
	// 			Password          string
	// 			RefreshToken      string
	// 			ExpRefreshToken   int
	// 			EncryptedPassword string
	// 		}
	// 	Output params:
	//		1. error or nil
	UpdatePassword(string, *model.User) error
	// Finds the user in the database.
	//
	// 	Input params:
	//		1. UserID int
	// 	Output params:
	//		1. type User struct {
	// 			ID                int
	// 			Email             string
	// 			Password          string
	// 			RefreshToken      string
	// 			ExpRefreshToken   int
	// 			EncryptedPassword string
	// 		}
	//		2. error or nil
	GetUserByID(int) (*model.User, error)
}

type CacheRepository interface {
	Set() (string, error)
	Get() (string, error)
}

type SessionRepository interface {
	// The function checks the existence of a session with the specified parameters.
	//
	// 	Input params:
	//		1. FingerPrint string
	//		2. RefreshToken string
	// 	Output param:
	//		1. type Session struct {
	//			ID           uint32
	//			UserId       uint32
	//			RefreshToken string
	//			Fingerprint  string
	//			ExpiresIn    int64
	//			CreatedAt    int64
	//		}
	//		2. error or nil
	VerifyRefreshToken(string, string) (*model.Session, error)
	//The function generates a Refresh Token and creates an entry in the session database with the specified fingerprint.
	//
	// 	Input params:
	//		1. UserID int
	//		2. FingerPrint string
	// 	Output param:
	//		1. type Session struct {
	//			ID           uint32
	//			UserId       uint32
	//			RefreshToken string
	//			Fingerprint  string
	//			ExpiresIn    int64
	//			CreatedAt    int64
	//		}
	//		2. error or nil
	CreateSession(uint32, string) (*model.Session, error)
	// Deletes the session with the specified fingerprint and Refresh Token.
	//
	// 	Input params:
	//		1. FingerPrint string
	//		2. RefreshToken string
	// 	Output param:
	//		1. error or nil
	DeleteSession(string, string) error
}
