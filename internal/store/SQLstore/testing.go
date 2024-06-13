package SQLstore

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
)

func TestDB(t *testing.T, databaseURL string) (*sql.DB, func(...string)) {
	t.Helper()

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}
	return db, func(tables ...string) {
		if len(tables) > 0 {
			_, err := db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE;", strings.Join(tables, ", ")))
			if err != nil {
				panic(fmt.Sprintf("the test table could not be cleared: %s", err))
			}
		}
		db.Close()
	}

}
