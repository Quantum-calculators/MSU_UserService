package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// test handle
func (s *APIServer) HandleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		HTML := `<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Document</title>
		</head>
		<body style="margin-left: 3vw; margin-top: 2vh;">
			<h1>Hello</h1>
		</body>
		</html>
		`
		w.Write([]byte(HTML))
	}
}

func (s *APIServer) TestHandler() http.HandlerFunc {
	type TestPostRequests struct {
		UserID int64  `json:"user_id"`
		Text   string `json:"text"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &TestPostRequests{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			fmt.Print(err)
			w.Write([]byte(err.Error()))
			return
		}
		fmt.Println(req)
	}
}
