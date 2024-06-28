package apiserver

import (
	"net/http"

	messageBroker "github.com/Quantum-calculators/MSU_UserService/internal/messageBroker"
	"github.com/Quantum-calculators/MSU_UserService/internal/store"
	"github.com/sirupsen/logrus"
)

type server struct {
	router       *http.ServeMux
	logger       *logrus.Logger
	store        store.Store
	rstore       store.RedisStore
	broker       messageBroker.Broker
	expAccess    int
	jwtSecretKey string
}

func newServer(store store.Store, redisstore store.RedisStore, broker messageBroker.Broker, ExpAccess int, JwtSecretKey string) *server {
	s := &server{
		router:       http.NewServeMux(),
		logger:       logrus.New(),
		store:        store,
		rstore:       redisstore,
		broker:       broker,
		expAccess:    ExpAccess,
		jwtSecretKey: JwtSecretKey,
	}
	s.ConfigureRouter()
	s.logger.Info("server is running")
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) ConfigureRouter() {
	s.router.HandleFunc("/hello", s.HandleHello())
	s.router.HandleFunc("/methods", s.Methods())
	s.router.HandleFunc("/registration", s.Registration())
	s.router.HandleFunc("/login", s.Login())
	s.router.HandleFunc("/get_access_token", s.AccessToken())
	s.router.HandleFunc("/logout", s.Logout())
	s.router.HandleFunc("/password_recovery", s.PasswordRecovery())
	s.router.HandleFunc("/confirmation_password_recovery", s.СonfirmationPasswordRecovery())

}
