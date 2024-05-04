package apiserver

import (
	"net/http"

	"github.com/PepsiKingIV/Lib_REST_API_server/internal/store"
	"github.com/sirupsen/logrus"
)

type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *http.ServeMux
	store  *store.Store
}

func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: http.NewServeMux(),
	}
}

func (s *APIServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}
	s.ConfigureRouter()

	if err := s.configureStore(); err != nil {
		return err
	}
	s.logger.Info("startung api server")
	return http.ListenAndServe(s.config.ServerAddr, s.router)
}

func (s *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

func (s *APIServer) configureStore() error {
	st := store.New(s.config.Store)
	if err := st.Open(); err != nil {
		return err
	}
	s.store = st
	return nil
}

func (s *APIServer) ConfigureRouter() {
	s.router.HandleFunc("/", s.HandleHello())
	s.router.HandleFunc("/test", s.TestHandler())
}
