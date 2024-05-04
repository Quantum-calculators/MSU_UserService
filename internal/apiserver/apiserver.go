package apiserver

import (
	"database/sql"
	"net/http"

	"github.com/Quantum-calculators/MSU_UserService/internal/store/SQLstore"
)

func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	store := SQLstore.New(db)
	srv := newServer(store)

	return http.ListenAndServe(config.ServerAddr, srv)
}

func newDB(DatabaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", DatabaseURL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
