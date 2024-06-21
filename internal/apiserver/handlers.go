package apiserver

import (
	"encoding/json"
	"errors"
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

type errorResponse struct {
	Status_code int    `json:"status_code"`
	Message     string `json:"message"`
}

func (s *server) error(w http.ResponseWriter, statusCode int, message string) errorResponse {
	resp := errorResponse{
		Status_code: statusCode,
		Message:     message,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Errorf("Error:%s", err)
		return errorResponse{Status_code: http.StatusInternalServerError}
	}
	return resp
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
			s.error(w, http.StatusMethodNotAllowed, "Only the POST method is allowed")
			s.logger.Warnf("%s\t%s\tError: %s", r.Method, r.URL, "MethodNotAllowed")
			return
		}
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, http.StatusUnprocessableEntity, "Incorrect request fields")
			s.logger.Warnf("%s\t%s\t  %d\tError: %s", r.Method, r.URL, http.StatusUnprocessableEntity, err.Error())
			return
		}

		expectedU, err := s.store.User().FindByEmail(req.Email)
		if err != nil {
			if err == store.ErrRecordNotFound {
				s.error(w, http.StatusNotFound, "There is no user with this email address")
				s.logger.Warnf("%s\t%s\t  %d\tError: %s", r.Method, r.URL, http.StatusNotFound, err.Error()) //сделать нормалныен ошибки, чтобы их можно было сообщать пользователю
				return
			}
			s.error(w, http.StatusInternalServerError, "Server error - the user could not be found")
			s.logger.Errorf("%s\t%s\t  %d\tError: %s", r.Method, r.URL, http.StatusNotFound, err.Error()) //сделать нормалныен ошибки, чтобы их можно было сообщать пользователю
			return
		}
		if !expectedU.ComparePassword(req.Password) {
			w.WriteHeader(http.StatusUnauthorized)
			s.logger.Infof("%s\t%s", r.Method, r.URL)
			return
		}
		if !expectedU.Verified {
			VerToken, err := token_generator.GenerateRandomString(64)
			if err != nil {
				s.error(w, http.StatusInternalServerError, "User registration error")
				s.logger.Errorf("%s\t%s\tError: %s", r.Method, r.URL, err.Error())
				return
			}

			message := brokerMessage{
				Email: expectedU.Email,
				URL:   VerificationURL + VerToken,
			}

			if err := s.store.User().UpdateVerificationToken(expectedU.Email, VerToken); err != nil {
				s.error(w, http.StatusInternalServerError, "Failed to update verification token")
				s.logger.Errorf("%s\t%s\tError: %s", r.Method, r.URL, err.Error())
				return
			}

			body, err := json.Marshal(message)
			if err != nil {
				s.error(w, http.StatusInternalServerError, "Failed to send verification message")
				s.logger.Errorf("%s\t%s\tError: %s", r.Method, r.URL, err.Error())
				return
			}

			err = s.broker.Message().SendMessage(body, "/VerifyEmail")
			if err != nil {
				s.error(w, http.StatusUnprocessableEntity, "")
			}

			s.error(w, http.StatusUnauthorized, "User not verified")
			s.logger.Infof("%s\t%s\t  %d\tError: %s", r.Method, r.URL, http.StatusUnauthorized, errors.New("user not verified"))
			return
		}

		// добавить проверку по полю Verified
		expectedU.Sanitize()
		session, err := s.store.Session().CreateSession(uint32(expectedU.ID), GetFingerPrint(r))
		if err != nil {
			s.error(w, http.StatusInternalServerError, "Failed to create a session")
			s.logger.Errorf("%s\t%s\t  %d\tError: %s", r.Method, r.URL, http.StatusInternalServerError, err.Error()) //сделать нормалныен ошибки, чтобы их можно было сообщать пользователю
			return
		}
		resp := response{
			RefreshToken:    session.RefreshToken,
			ExpRefreshToken: int(session.ExpiresIn),
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			s.error(w, http.StatusUnprocessableEntity, "")
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
			s.error(w, http.StatusMethodNotAllowed, "Only the POST method is allowed")
			s.logger.Warnf("%s\t%s\tError: %s", r.Method, r.URL, "MethodNotAllowed")
			return
		}
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, http.StatusUnprocessableEntity, "Incorrect request fields")
			s.logger.Warnf("%s\t%s\t  %d\tError: %s", r.Method, r.URL, http.StatusUnprocessableEntity, err.Error())
			return
		}
		if err := s.store.Session().DeleteSession(GetFingerPrint(r), req.RefreshToken); err != nil {
			s.error(w, http.StatusInternalServerError, "Failed to delete a session")
			s.logger.Errorf("%s\t%s\t  %d\tError: %s", r.Method, r.URL, http.StatusInternalServerError, err.Error()) //сделать нормалныен ошибки, чтобы их можно было сообщать пользователю
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
			s.error(w, http.StatusMethodNotAllowed, "Only the GET method is allowed")
			s.logger.Warnf("%s\t%s\tError: %s", r.Method, r.URL, "MethodNotAllowed")
			return
		}
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, http.StatusUnprocessableEntity, "Incorrect request fields")
			s.logger.Warnf("%s\t%s\t  %d\tError: %s", r.Method, r.URL, http.StatusUnprocessableEntity, err.Error())
			return
		}
		session, err := s.store.Session().VerifyRefreshToken(GetFingerPrint(r), req.RefreshToken)
		if err != nil {
			w.WriteHeader(http.StatusNonAuthoritativeInfo)
			s.error(w, http.StatusUnauthorized, "The session for this user was not found")
			s.logger.Errorf("%s\t%s\tError: %s", r.Method, r.URL, err.Error())
			return
		}
		user, err := s.store.User().GetUserByID(int(session.UserId))
		if err != nil {
			s.error(w, http.StatusInternalServerError, "")
			s.logger.Errorf("%s\t%s\tError: %s", r.Method, r.URL, err.Error())
			return
		}
		payload := jwt.MapClaims{
			"sub": user.Email,
			"exp": time.Now().Add(time.Duration(s.expAccess)).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
		accessToken, err := token.SignedString([]byte(s.jwtSecretKey))
		if err != nil {
			s.error(w, http.StatusInternalServerError, "Failed to generate accessToken")
			s.logger.Errorf("%s\t%s\tError: %s", r.Method, r.URL, err.Error())
			return
		}
		resp := response{
			AccessToken:     accessToken,
			RefreshToken:    session.RefreshToken,
			ExpRefreshToken: int(session.ExpiresIn),
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			s.error(w, http.StatusUnprocessableEntity, "")
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
			s.error(w, http.StatusMethodNotAllowed, "Only the POST method is allowed")
			s.logger.Warnf("%s\t%s\tError: %s", r.Method, r.URL, "MethodNotAllowed")
			return
		}
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, http.StatusUnprocessableEntity, "Incorrect request fields")
			s.logger.Warnf("%s\t%s\tError: %s", r.Method, r.URL, "Incorrect request fields")
			return
		}
		VerToken, err := token_generator.GenerateRandomString(64)
		if err != nil {
			s.error(w, http.StatusInternalServerError, "User registration error")
			s.logger.Errorf("%s\t%s\tError: %s", r.Method, r.URL, err.Error())
			return
		}
		u := &model.User{
			Email:             req.Email,
			Password:          req.Password,
			VerificationToken: VerToken,
		}
		if err := s.store.User().Create(u); err != nil {
			s.error(w, http.StatusUnprocessableEntity, err.Error())
			s.logger.Errorf("%s\t%s\tError: %s", r.Method, r.URL, err.Error())
			return
		}
		u.Sanitize()

		message := brokerMessage{
			Email: u.Email,
			URL:   VerificationURL + VerToken,
		}

		body, err := json.Marshal(message)
		if err != nil {
			s.error(w, http.StatusInternalServerError, "Failed to send verification message")
			s.logger.Errorf("%s\t%s\tError: %s", r.Method, r.URL, err.Error())
			return
		}

		err = s.broker.Message().SendMessage(body, "/VerifyEmail")
		if err != nil {
			s.error(w, http.StatusUnprocessableEntity, "")
		}
		if err := json.NewEncoder(w).Encode(u); err != nil {
			s.error(w, http.StatusUnprocessableEntity, "")
			s.logger.Errorf("%s\t%s\tError: %s", r.Method, r.URL, err.Error())
			return
		}
		w.WriteHeader(http.StatusCreated)
		s.logger.Infof("%s\t%s", r.Method, r.URL)
	}
}

func (s *server) Verification() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			s.error(w, http.StatusMethodNotAllowed, "Only the POST method is allowed")
			s.logger.Warnf("%s\t%s\tError: %s", r.Method, r.URL, "MethodNotAllowed")
			return
		}
		token := r.PathValue("token")
		email := r.PathValue("email")

		ok, err := s.store.User().CheckVerificationToken(email, token)
		if err != nil {
			s.error(w, http.StatusInternalServerError, "Server error")
			s.logger.Errorf("%s\t%s\tError: %s", r.Method, r.URL, err.Error())
			return
		}
		err = s.store.User().SetVerify(email, ok)
		if err != nil {
			s.error(w, http.StatusInternalServerError, "Server error")
			s.logger.Errorf("%s\t%s\tError: %s", r.Method, r.URL, err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		s.logger.Infof("%s\t%s", r.Method, r.URL)
	}
}
