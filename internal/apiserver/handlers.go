package apiserver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Quantum-calculators/MSU_UserService/internal/model"
	"github.com/Quantum-calculators/MSU_UserService/internal/store"
	token_generator "github.com/Quantum-calculators/MSU_UserService/internal/tokenGenerator"
	"github.com/golang-jwt/jwt"
)

// Перенести в конфигурацию
const VerificationURL = "127.0.0.1:8080/verification/"

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
			s.error(w, http.StatusMethodNotAllowed, ErrorOnlyPostMethod.Error())
			return
		}
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, http.StatusUnprocessableEntity, ErrorRequestFields.Error())
			return
		}

		expectedU, err := s.store.User().FindByEmail(r.Context(), req.Email)
		if err != nil {
			if err == store.ErrRecordNotFound {
				s.error(w, http.StatusNotFound, ErrorNotFoundUserWithEmail.Error())
				return
			}
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

			if err := s.store.User().UpdateVerificationToken(r.Context(), expectedU.Email, VerToken); err != nil {
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
		session, err := s.store.Session().CreateSession(r.Context(), uint32(expectedU.ID), GetFingerPrint(r))
		if err != nil {
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
		if err := s.store.Session().DeleteSession(r.Context(), GetFingerPrint(r), req.RefreshToken); err != nil {
			s.error(w, http.StatusInternalServerError, ErrorServer.Error())
			return
		}
	}
}

func (s *server) AccessToken() http.HandlerFunc {
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
			s.error(w, http.StatusMethodNotAllowed, ErrorOnlyGetMethod.Error())
			return
		}
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, http.StatusUnprocessableEntity, ErrorRequestFields.Error())
			return
		}
		session, err := s.store.Session().VerifyRefreshToken(r.Context(), GetFingerPrint(r), req.RefreshToken)
		if err != nil {
			s.error(w, http.StatusUnauthorized, ErrorUserUnauth.Error())
			return
		}
		user, err := s.store.User().GetUserByID(r.Context(), int(session.UserId))
		if err != nil {
			s.error(w, http.StatusInternalServerError, "")
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
		if err := s.store.User().Create(r.Context(), u); err != nil {
			s.error(w, http.StatusUnprocessableEntity, err.Error())
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
		if err := json.NewEncoder(w).Encode(u); err != nil {
			s.error(w, http.StatusInternalServerError, ErrorServer.Error())
			return
		}
		// w.WriteHeader(http.StatusCreated)
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
		if err != nil {
			s.error(w, http.StatusInternalServerError, ErrorServer.Error())
			return
		}
		err = s.store.User().SetVerify(r.Context(), email, ok)
		if err != nil {
			s.error(w, http.StatusInternalServerError, ErrorServer.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
