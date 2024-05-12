package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Quantum-calculators/MSU_UserService/internal/model"
)

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

func (s *server) TestHandle() http.HandlerFunc {
	type TestPostRequests struct {
		UserID int64  `json:"user_id"`
		Text   string `json:"text"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &TestPostRequests{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			s.logger.Warnf("%s\t%s\tError: %s", r.Method, r.URL, err.Error())
			return
		}
		fmt.Println(req)
		s.logger.Infof("%s\t%s", r.Method, r.URL)
		w.WriteHeader(http.StatusOK)
	}
}

func (s *server) TestRedis() http.HandlerFunc {
	type TestPostRequests struct {
		Text string `json:"text"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &TestPostRequests{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error())) // переделать: пользователю не должно раскрываться устройство сервера
			s.logger.Warnf("%s\t%s\tError: %s", r.Method, r.URL, err.Error())
			return
		}
		s.logger.Infof("%s\t%s", r.Method, r.URL)
		w.WriteHeader(http.StatusOK)
	}
}

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
	}
}
