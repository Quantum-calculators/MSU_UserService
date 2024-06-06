package apiserver

import (
	"net/http"

	messageBroker "github.com/Quantum-calculators/MSU_UserService/internal/messageBroker"
	"github.com/Quantum-calculators/MSU_UserService/internal/store"
	"github.com/sirupsen/logrus"
)

type server struct {
	router *http.ServeMux
	logger *logrus.Logger
	store  store.Store
	Rstore store.RedisStore
	broker *messageBroker.Broker
}

func newServer(store store.Store, redisstore store.RedisStore, broker *messageBroker.Broker) *server {
	s := &server{
		router: http.NewServeMux(),
		logger: logrus.New(),
		store:  store,
		Rstore: redisstore,
	}
	s.ConfigureRouter()
	s.logger.Info("server is running")
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) ConfigureRouter() {
	s.router.HandleFunc("/", s.HandleHello())
	s.router.HandleFunc("/registration", s.Registration())
	s.router.HandleFunc("/login", s.Login())
	s.router.HandleFunc("/GAT", s.GetAccessToken())
	s.router.HandleFunc("/logout", s.Logout())
}
