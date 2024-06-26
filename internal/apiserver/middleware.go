package apiserver

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/sirupsen/logrus"
)

func (s *server) PanicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		newReq := req.WithContext(req.Context())
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				log.Println(string(debug.Stack()))
			}
		}()
		next.ServeHTTP(w, newReq)
	})
}

func (s *server) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, req)
		s.logger.Logf(logrus.InfoLevel, "%s %s %s", req.Method, req.RequestURI, time.Since(start))
	})
}
