package apiserver

import (
	"net/http"

	"github.com/Quantum-calculators/MSU_UserService/internal/store"
	"github.com/sirupsen/logrus"
)

type server struct {
	router *http.ServeMux
	logger *logrus.Logger
	store  store.Store
	Rstore store.RedisStore
}

func newServer(store store.Store, redisstore store.RedisStore) *server {
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
	s.router.HandleFunc("/test", s.TestHandle())
	s.router.HandleFunc("/testJWT", s.TestRedis())
	s.router.HandleFunc("/auth", s.Registration())
}
