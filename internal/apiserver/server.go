package apiserver

import (
	"net/http"

	"github.com/Quantum-calculators/MSU_UserService/internal/store"
	"github.com/sirupsen/logrus"
)

type server struct {
	router   *http.ServeMux
	logger   *logrus.Logger
	store    store.Store
	jwtstore store.JWTStore
}

func newServer(store store.Store, jwtstore store.JWTStore) *server {
	s := &server{
		router:   http.NewServeMux(),
		logger:   logrus.New(),
		store:    store,
		jwtstore: jwtstore,
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
	s.router.HandleFunc("/test", s.TestHandler())
}
