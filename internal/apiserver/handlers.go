package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Quantum-calculators/MSU_UserService/internal/model"
	"github.com/golang-jwt/jwt"
)

// Перенести в конфигурацию
const jwtSecretKey = "test"
const AccessTokenExp = 10 // min

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

func GetFingerPrint(r *http.Request) string {
	fingerPrint := r.Header.Get("Accept-Language") + r.Header.Get("Sec-Ch-Ua-Platform") + r.Header.Get("User-Agent")
	return fingerPrint
}

// TODO: добавить описание ошибок в ответе пользователю
// TODO: Добавить проверку на метод Post
func (s *server) Login() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		RefreshToken    string `json:"refreshToken"`
		ExpRefreshToken int    `json:"expRefreshToken"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			s.logger.Warnf("%s\t%s\tError: %s", r.Method, r.URL, "MethodNotAllowed")
			return
		}
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
		session, err := s.store.Session().CreateSession(uint32(expectedU.ID), GetFingerPrint(r))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s.logger.Errorf("%s\t%s\tError: %s", r.Method, r.URL, err.Error())
			return
		}
		resp := response{
			RefreshToken:    session.RefreshToken,
			ExpRefreshToken: int(session.ExpiresIn),
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s.logger.Errorf("%s\t%s\tError: %s", r.Method, r.URL, err.Error())
			return
		}
		s.logger.Infof("%s\t%s", r.Method, r.URL)
	}
}

func (s *server) Logout() http.HandlerFunc {
	type request struct {
		RefreshToken string `json:"refreshToken"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			s.logger.Warnf("%s\t%s\tError: %s", r.Method, r.URL, "MethodNotAllowed")
			return
		}
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			s.logger.Warnf("%s\t%s\tError: %s", r.Method, r.URL, err.Error())
			return
		}
		if err := s.store.Session().DeleteSession(GetFingerPrint(r), req.RefreshToken); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s.logger.Errorf("%s\t%s\tError: %s", r.Method, r.URL, err.Error())
			return
		}
	}
}

func (s *server) GetAccessToken() http.HandlerFunc {
	type request struct {
		RefreshToken string `json:"refreshToken"`
	}
	type response struct {
		AccessToken     string `json:"accessToken"`
		RefreshToken    string `json:"refreshToken"`
		ExpRefreshToken int    `json:"expRefreshToken"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			s.logger.Warnf("%s\t%s\tError: %s", r.Method, r.URL, "MethodNotAllowed")
			return
		}
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			s.logger.Warnf("%s\t%s\tError: %s", r.Method, r.URL, err.Error())
			return
		}
		session, err := s.store.Session().VerifyRefreshToken(GetFingerPrint(r), req.RefreshToken)
		if err != nil {
			w.WriteHeader(http.StatusNonAuthoritativeInfo)
			s.logger.Errorf("%s\t%s\tError: %s", r.Method, r.URL, err.Error())
			return
		}
		fmt.Println(session)
		user, err := s.store.User().GetUserByID(int(session.UserId))
		fmt.Println(user)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			s.logger.Errorf("%s\t%s\tError: %s", r.Method, r.URL, err.Error())
			return
		}
		payload := jwt.MapClaims{
			"sub": user.Email,
			"exp": time.Now().Add(AccessTokenExp).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
		accessToken, err := token.SignedString([]byte(jwtSecretKey))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s.logger.Errorf("%s\t%s\tError: %s", r.Method, r.URL, err.Error())
			return
		}
		resp := response{
			AccessToken:     accessToken,
			RefreshToken:    session.RefreshToken,
			ExpRefreshToken: int(session.ExpiresIn),
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s.logger.Errorf("%s\t%s\tError: %s", r.Method, r.URL, err.Error())
			return
		}
		s.logger.Infof("%s\t%s", r.Method, r.URL)
	}
}

func (s *server) Registration() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			s.logger.Warnf("%s\t%s\tError: %s", r.Method, r.URL, "MethodNotAllowed")
			return
		}
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
