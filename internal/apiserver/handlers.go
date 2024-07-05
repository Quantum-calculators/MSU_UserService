package apiserver

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Quantum-calculators/MSU_UserService/internal/model"
	"github.com/Quantum-calculators/MSU_UserService/internal/store"
	token_generator "github.com/Quantum-calculators/MSU_UserService/internal/tokenGenerator"
	"github.com/golang-jwt/jwt"
)

// Перенести в конфигурацию
const VerificationURL = "127.0.0.1:8080/verification/"
const PasswordRecoveryURL = "127.0.0.1:8080/verification/"

type brokerMessage struct {
	Email string
	URL   string
}

func GetFingerPrint(r *http.Request) string {
	fingerPrint := r.Header.Get("Accept-Language") + r.Header.Get("Sec-Ch-Ua-Platform") + r.Header.Get("User-Agent")
	return fingerPrint
}

type errorResponse struct {
	Status_code int    `json:"status_code"`
	Message     string `json:"message"`
}

func (s *server) error(w http.ResponseWriter, statusCode int, message string) errorResponse {
	resp := errorResponse{
		Status_code: statusCode,
		Message:     message,
	}
	s.logger.Errorf("status code: %d, Error: %s", statusCode, message)
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Errorf("status code: %d, Error: %s", http.StatusInternalServerError, err.Error())
		return errorResponse{Status_code: http.StatusInternalServerError}
	}
	return resp
}

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
		// s.error(w, http.StatusMethodNotAllowed, "Only the POST method is allowed")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(HTML))
	}
}

func (s *server) Methods() http.HandlerFunc {
	type response struct {
		URI []string `json:"uri"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		resp := response{
			URI: []string{
				"/registration",
				"/login",
				"/GAT",
				"/logout",
				"/verification/{token}/{email}",
			},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			s.error(w, http.StatusUnprocessableEntity, ErrorServer.Error())
			return
		}
	}
}

func (s *server) Login() http.HandlerFunc {
	type request struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		FingerPrint string `json:"fingerPrint"`
	}
	type response struct {
		RefreshToken    string `json:"refreshToken"`
		ExpRefreshToken int    `json:"expRefreshToken"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			s.error(w, http.StatusMethodNotAllowed, ErrorOnlyPostMethod.Error())
			return
		}
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, http.StatusUnprocessableEntity, ErrorRequestFields.Error())
			return
		}

		expectedU, err := s.store.User().FindByEmail(r.Context(), req.Email)
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			s.error(w, http.StatusGatewayTimeout, ErrorServer.Error())
			return
		case errors.Is(err, store.ErrRecordNotFound): // <- плохо тянуть константы из зависимостей
			s.error(w, http.StatusForbidden, ErrorNotFoundUserWithEmail.Error())
			return
		case err != nil:
			s.error(w, http.StatusInternalServerError, ErrorServer.Error())
			return
		}
		if !expectedU.ComparePassword(req.Password) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if !expectedU.Verified {
			VerToken, err := token_generator.GenerateRandomString(64)
			if err != nil {
				s.error(w, http.StatusInternalServerError, ErrorServer.Error())
				return
			}

			message := brokerMessage{
				Email: expectedU.Email,
				URL:   VerificationURL + VerToken,
			}

			err = s.store.User().UpdateVerificationToken(r.Context(), expectedU.Email, VerToken)
			switch {
			case errors.Is(err, context.DeadlineExceeded):
				s.error(w, http.StatusGatewayTimeout, ErrorServer.Error())
				return
			case err != nil:
				s.error(w, http.StatusInternalServerError, ErrorServer.Error())
				return
			}

			body, err := json.Marshal(message)
			if err != nil {
				s.error(w, http.StatusInternalServerError, ErrorServer.Error())
				return
			}

			err = s.broker.Message().SendMessage(body, "/VerifyEmail")
			if err != nil {
				s.error(w, http.StatusUnprocessableEntity, "")
			}

			s.error(w, http.StatusUnauthorized, ErrorUserNotVerified.Error())
			return
		}

		// добавить проверку по полю Verified
		expectedU.Sanitize()
		session, err := s.store.Session().CreateSession(r.Context(), expectedU.Email, req.FingerPrint)
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			s.error(w, http.StatusGatewayTimeout, ErrorServer.Error())
			return
		case err != nil:
			s.error(w, http.StatusInternalServerError, ErrorServer.Error())
			return
		}
		resp := response{
			RefreshToken:    session.RefreshToken,
			ExpRefreshToken: int(session.ExpiresIn),
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			s.error(w, http.StatusUnprocessableEntity, ErrorServer.Error())
			return
		}
	}
}

func (s *server) Logout() http.HandlerFunc {
	type request struct {
		RefreshToken string `json:"refreshToken"`
		FingerPrint  string `json:"fingerPrint"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			s.error(w, http.StatusMethodNotAllowed, ErrorOnlyPostMethod.Error())
			return
		}
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, http.StatusUnprocessableEntity, ErrorServer.Error())
			return
		}
		err := s.store.Session().DeleteSession(r.Context(), req.FingerPrint, req.RefreshToken)
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			s.error(w, http.StatusGatewayTimeout, ErrorServer.Error())
			return
		case err != nil:
			s.error(w, http.StatusInternalServerError, ErrorServer.Error())
			return
		}
	}
}

func (s *server) AccessToken() http.HandlerFunc {
	type request struct {
		RefreshToken string `json:"refreshToken"`
		FingerPrint  string `json:"fingerPrint"`
	}
	type response struct {
		AccessToken     string `json:"accessToken"`
		RefreshToken    string `json:"refreshToken"`
		ExpRefreshToken int    `json:"expRefreshToken"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			s.error(w, http.StatusMethodNotAllowed, ErrorOnlyGetMethod.Error())
			return
		}
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, http.StatusUnprocessableEntity, ErrorRequestFields.Error())
			return
		}
		session, err := s.store.Session().VerifyRefreshToken(r.Context(), req.FingerPrint, req.RefreshToken)
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			s.error(w, http.StatusGatewayTimeout, ErrorServer.Error())
			return
		case err != nil:
			s.error(w, http.StatusUnauthorized, ErrorUserUnauth.Error())
			return
		}
		user, err := s.store.User().FindByEmail(r.Context(), session.Email)
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			s.error(w, http.StatusGatewayTimeout, ErrorServer.Error())
			return
		case err != nil:
			s.error(w, http.StatusInternalServerError, ErrorServer.Error())
			return
		}
		payload := jwt.MapClaims{
			"sub": user.Email,
			"exp": time.Now().Add(time.Duration(s.expAccess)).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
		accessToken, err := token.SignedString([]byte(s.jwtSecretKey))
		if err != nil {
			s.error(w, http.StatusInternalServerError, ErrorServer.Error())
			return
		}
		resp := response{
			AccessToken:     accessToken,
			RefreshToken:    session.RefreshToken,
			ExpRefreshToken: int(session.ExpiresIn),
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			s.error(w, http.StatusUnprocessableEntity, ErrorServer.Error())
			return
		}
	}
}

func (s *server) Registration() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			s.error(w, http.StatusMethodNotAllowed, ErrorOnlyPostMethod.Error())
			return
		}
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, http.StatusUnprocessableEntity, ErrorRequestFields.Error())
			return
		}
		VerToken, err := token_generator.GenerateRandomString(64)
		if err != nil {
			s.error(w, http.StatusInternalServerError, ErrorServer.Error())
			return
		}
		u := &model.User{
			Email:             req.Email,
			Password:          req.Password,
			VerificationToken: VerToken,
		}

		err = s.store.User().Create(r.Context(), u)
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			s.error(w, http.StatusGatewayTimeout, ErrorServer.Error())
			return
		case err != nil:
			s.error(w, http.StatusForbidden, err.Error())
			return
		}
		u.Sanitize()
		message := brokerMessage{
			Email: u.Email,
			URL:   VerificationURL + VerToken,
		}
		body, err := json.Marshal(message)
		if err != nil {
			s.error(w, http.StatusInternalServerError, ErrorServer.Error())
			return
		}

		err = s.broker.Message().SendMessage(body, "/VerifyEmail")
		if err != nil {
			s.error(w, http.StatusInternalServerError, ErrorServer.Error())
			return
		}
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(u); err != nil {
			s.error(w, http.StatusInternalServerError, ErrorServer.Error())
			return
		}
	}
}

func (s *server) Verification() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			s.error(w, http.StatusMethodNotAllowed, ErrorOnlyPostMethod.Error())
			return
		}
		token := r.PathValue("token")
		email := r.PathValue("email")

		ok, err := s.store.User().CheckVerificationToken(r.Context(), email, token)
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			s.error(w, http.StatusGatewayTimeout, ErrorServer.Error())
			return
		case errors.Is(err, store.ErrRecordNotFound):
			s.error(w, http.StatusNotFound, ErrNotFound.Error())
			return
		case err != nil:
			s.error(w, http.StatusInternalServerError, ErrorServer.Error())
			return
		}
		if ok {
			err = s.store.User().SetVerify(r.Context(), email, ok)
			switch {
			case errors.Is(err, context.DeadlineExceeded):
				s.error(w, http.StatusGatewayTimeout, ErrorServer.Error())
				return
			case errors.Is(err, store.ErrRecordNotFound):
				s.error(w, http.StatusBadRequest, ErrNotFound.Error())
				return
			case err != nil:
				s.error(w, http.StatusInternalServerError, ErrorServer.Error())
				return
			}
			w.WriteHeader(http.StatusOK)
		}
		w.WriteHeader(http.StatusNotAcceptable)
	}
}

func (s *server) СonfirmationPasswordRecovery() http.HandlerFunc {
	type request struct {
		NewPassword string `json:"new_password"`
		Token       string `json:"token"`
		Email       string `json:"email"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			s.error(w, http.StatusMethodNotAllowed, ErrorOnlyPostMethod.Error())
			return
		}
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, http.StatusUnprocessableEntity, ErrorRequestFields.Error())
			return
		}

		expectedToken, err := s.store.User().GetRecoveryPasswordToken(r.Context(), req.Email)
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			s.error(w, http.StatusGatewayTimeout, ErrorServer.Error())
			return
		case errors.Is(err, sql.ErrNoRows):
			s.error(w, http.StatusForbidden, ErrNotFound.Error())
			return
		}
		if expectedToken != req.Token {
			s.error(w, http.StatusForbidden, ErrorServer.Error())
			return
		}
		u := &model.User{
			Email: req.Email,
		}
		err = s.store.User().UpdatePassword(r.Context(), req.NewPassword, u)
		fmt.Println(req.NewPassword)
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			s.error(w, http.StatusGatewayTimeout, ErrorServer.Error())
			return
		case err != nil:
			s.error(w, http.StatusInternalServerError, err.Error())
			return
		}
		err = s.store.Session().DeleteAllSession(r.Context(), req.Email)
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			s.error(w, http.StatusGatewayTimeout, ErrorServer.Error())
			return
		case err != nil:
			s.error(w, http.StatusInternalServerError, ErrorServer.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (s *server) PasswordRecovery() http.HandlerFunc {
	type request struct {
		Email string `json:"email"`
	}
	type response struct {
		Email string `json:"email"`
		Token string `json:"token"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			s.error(w, http.StatusMethodNotAllowed, ErrorOnlyPostMethod.Error())
			return
		}
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, http.StatusUnprocessableEntity, ErrorRequestFields.Error())
			return
		}
		token, err := token_generator.GenerateRandomString(128)
		if err != nil {
			s.error(w, http.StatusInternalServerError, ErrorServer.Error())
			return
		}
		err = s.store.User().CreatePasswordRecoveryToken(r.Context(), req.Email, token)
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			s.error(w, http.StatusGatewayTimeout, ErrorServer.Error())
			return
		case err != nil:
			s.error(w, http.StatusInternalServerError, ErrorServer.Error())
			return
		}

		resp := response{
			Email: req.Email,
			Token: token,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			s.error(w, http.StatusUnprocessableEntity, ErrorServer.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (s *server) ChangePassword() http.HandlerFunc {
	type request struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		NewPassword string `json:"new_password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			s.error(w, http.StatusMethodNotAllowed, ErrorOnlyPostMethod.Error())
			return
		}
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, http.StatusUnprocessableEntity, ErrorRequestFields.Error())
			return
		}

		expectedU, err := s.store.User().FindByEmail(r.Context(), req.Email)
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			s.error(w, http.StatusGatewayTimeout, ErrorServer.Error())
			return
		case errors.Is(err, store.ErrRecordNotFound): // <- плохо тянуть константы из зависимостей
			s.error(w, http.StatusNotFound, ErrorNotFoundUserWithEmail.Error())
			return
		case err != nil:
			s.error(w, http.StatusInternalServerError, ErrorServer.Error())
			return
		}
		if !expectedU.ComparePassword(req.Password) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		err = s.store.User().UpdatePassword(r.Context(), req.NewPassword, expectedU)
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			s.error(w, http.StatusGatewayTimeout, ErrorServer.Error())
			return
		case err != nil:
			s.error(w, http.StatusInternalServerError, ErrorServer.Error())
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}
