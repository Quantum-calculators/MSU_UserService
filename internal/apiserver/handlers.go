package apiserver

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"time"

	"github.com/Quantum-calculators/MSU_UserService/internal/model"
	"github.com/golang-jwt/jwt"
)

// Перенести в конфигурацию
const AccessTokenExp = 60   // min
const RefreshTokenExp = 720 // hours
const jwtSecretKey = "test"

// test handle
func (s *server) HandleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		HTML := `<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Document</title>
		</head>
		<body style="margin-left: 3vw; margin-top: 2vh;">
			<h1>Hello</h1>
		</body> 
		</html>
		`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(HTML))
		s.logger.Infof("%s\t%s", r.Method, r.URL)
	}
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}
	return string(ret), nil
}

// TODO: добавить описание ошибок в ответе пользователю
// TODO: Добавить проверку на метод Post
func (s *server) Login() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			s.logger.Warnf("%s\t%s\tError: %s", r.Method, r.URL, err.Error())
			return
		}

		expectedU, err := s.store.User().FindByEmail(req.Email)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s.logger.Errorf("%s\t%s\tError: %s", r.Method, r.URL, err.Error()) //сделать нормалныен ошибки, чтобы их можно было сообщать пользователю
			return
		}
		if !expectedU.ComparePassword(req.Password) {
			w.WriteHeader(http.StatusUnauthorized)
			s.logger.Infof("%s\t%s", r.Method, r.URL)
			return
		}
		expectedU.Sanitize()

		payload := jwt.MapClaims{
			"sub": req.Email,
			"exp": time.Now().Add(AccessTokenExp).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
		refreshToken, err := GenerateRandomString(128)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s.logger.Errorf("%s\t%s\tError: %s", r.Method, r.URL, err.Error())
			return
		}
		accessToken, err := token.SignedString([]byte(jwtSecretKey))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s.logger.Errorf("%s\t%s\tError: %s", r.Method, r.URL, err.Error())
			return
		}
		resp := response{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s.logger.Errorf("%s\t%s\tError: %s", r.Method, r.URL, err.Error())
			return
		}
		if err := s.store.User().SetRefreshToken(
			refreshToken,
			int(time.Now().Add(RefreshTokenExp).Unix()),
			expectedU,
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s.logger.Errorf("%s\t%s\tError: %s", r.Method, r.URL, err.Error())
			return
		}
		s.logger.Infof("%s\t%s", r.Method, r.URL)
	}
}

// TODO: Добавить проверку на метод Post
func (s *server) Registration() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			s.logger.Warnf("%s\t%s\tError: %s", r.Method, r.URL, err.Error())
			return
		}

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}

		if err := s.store.User().Create(u); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s.logger.Errorf("%s\t%s\tError: %s", r.Method, r.URL, err.Error()) //сделать нормалныен ошибки, чтобы их можно было сообщать пользователю
			return
		}
		u.Sanitize()
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(u); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s.logger.Errorf("%s\t%s\tError: %s", r.Method, r.URL, err.Error())
			return
		}
		s.logger.Infof("%s\t%s", r.Method, r.URL)
	}
}

func (s *server) UserCheck() http.HandlerFunc {
	type request struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
