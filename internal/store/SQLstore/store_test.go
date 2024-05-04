package SQLstore_test

import (
	"os"
	"testing"
)

var (
	databaseURL string
)

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("DB_url")
	if databaseURL == "" {
		databaseURL = "host=localhost dbname=restapi_test sslmode=disable"
	}
	os.Exit(m.Run())
}
