package apiserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Quantum-calculators/MSU_UserService/internal/messageBroker/testbroker"
	"github.com/Quantum-calculators/MSU_UserService/internal/model"
	"github.com/Quantum-calculators/MSU_UserService/internal/store/redisStore"
	"github.com/Quantum-calculators/MSU_UserService/internal/store/testStore"
	"github.com/stretchr/testify/assert"
)

func TestServer_RegistrationUser(t *testing.T) {
	store := testStore.New()
	Rstore := redisStore.New_Test()
	broker := testbroker.New()

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "registration",
			payload: map[string]interface{}{
				"email":    "test1@mail.com",
				"password": "valid_password",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "registration",
			payload: map[string]interface{}{
				"email":    "test2@mail.com",
				"password": "invalid", // слишком короткий пароль
			},
			expectedCode: http.StatusForbidden,
		},
		{
			name: "registration",
			payload: map[string]interface{}{
				"invalid_payload": "test3@mail.com", // неверное поле тела запроса
				"password":        "valid_password",
			},
			expectedCode: http.StatusForbidden,
		},
		{
			name: "registration",
			payload: map[string]interface{}{
				"invalid_payload": "test1@mail.com", // email такой же как у 1-ого пользователя
				"password":        "valid_password",
			},
			expectedCode: http.StatusForbidden,
		},
	}

	secretKey := "secret"
	s := newServer(store, Rstore, broker, 1000, secretKey)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			body, err := json.Marshal(tc.payload)
			if err != nil {
				assert.NoError(t, err)
				return
			}
			req, _ := http.NewRequest(http.MethodPost, "/registration", bytes.NewReader(body))
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_LoginUser(t *testing.T) {
	store := testStore.New()
	Rstore := redisStore.New_Test()
	broker := testbroker.New()

	testCases := []struct {
		prepare      func()
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			prepare: func() {
				u := &model.User{
					Email:    "test1@mail.com",
					Password: "valid_password",
				}
				store.User().Create(nil, u)
				store.User().SetVerify(nil, u.Email, true)
			},
			name: "Login",
			payload: map[string]interface{}{
				"email":    "test1@mail.com",
				"password": "valid_password",
			},
			expectedCode: http.StatusOK,
		},
		{
			prepare: func() {
				u := &model.User{
					Email:    "test2@mail.com",
					Password: "valid_password",
				}
				store.User().Create(nil, u)
				// store.User().SetVerify(nil, u.Email, true)  <- пользователь не прошел верификацию
			},
			name: "Login",
			payload: map[string]interface{}{
				"email":    "test2@mail.com",
				"password": "valid_password",
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			prepare: func() {
				u := &model.User{
					Email:    "test3@mail.com",
					Password: "valid_password",
				}
				store.User().Create(nil, u)
				store.User().SetVerify(nil, u.Email, true)
			},
			name: "Login",
			payload: map[string]interface{}{
				"email": "test2@mail.com", // 	<- в payload отсутствует необходимое поле
			},
			expectedCode: http.StatusUnauthorized,
		},
	}

	secretKey := "secret"
	s := newServer(store, Rstore, broker, 1000, secretKey)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare()
			rec := httptest.NewRecorder()
			body, err := json.Marshal(tc.payload)
			if err != nil {
				assert.NoError(t, err)
				return
			}
			req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_VerificationUser(t *testing.T) {
	store := testStore.New()
	Rstore := redisStore.New_Test()
	broker := testbroker.New()

	testCases := []struct {
		prepare      func()
		name         string
		payload      []string
		expectedCode int
	}{
		{
			prepare: func() {
				u := &model.User{
					Email:             "test1@mail.com",
					Password:          "valid_password",
					VerificationToken: "verification_token",
				}
				store.User().Create(nil, u)
			},
			name: "Verification",
			payload: []string{
				"verification_token",
				"test1@mail.com",
			},
			expectedCode: http.StatusOK,
		},
		{
			prepare: func() {
				u := &model.User{
					Email:             "test2@mail.com",
					Password:          "valid_password",
					VerificationToken: "verification_token",
				}
				store.User().Create(nil, u)
			},
			name: "Verification",
			payload: []string{
				"verification_token",
				"test@mail.com", // 		<- несуществующий пользователь
			},
			expectedCode: http.StatusNotFound,
		},
		{
			prepare: func() {
				u := &model.User{
					Email:             "test3@mail.com",
					Password:          "valid_password",
					VerificationToken: "verification_token",
				}
				store.User().Create(nil, u)
			},
			name: "Verification",
			payload: []string{
				"invalid_verification_token", //		<- неверный токен
				"test3@mail.com",
			},
			expectedCode: http.StatusNotAcceptable,
		},
	}

	secretKey := "secret"
	s := newServer(store, Rstore, broker, 1000, secretKey)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare()
			rec := httptest.NewRecorder()
			body, err := json.Marshal(tc.payload)
			if err != nil {
				assert.NoError(t, err)
				return
			}
			addr := fmt.Sprintf("/verification/%s/%s", tc.payload[0], tc.payload[1])
			req, _ := http.NewRequest(http.MethodGet, addr, bytes.NewReader(body))
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
