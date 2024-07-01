package apiserver

import "errors"

var (
	ErrorOnlyPostMethod        = errors.New("only the POST method is allowed")
	ErrorRequestFields         = errors.New("incorrect request fields")
	ErrorNotFoundUserWithEmail = errors.New("there is no user with this email address")
	ErrorServer                = errors.New("server error")
	ErrorUserNotVerified       = errors.New("user not verified")
	ErrorOnlyGetMethod         = errors.New("only the GET method is allowed")
	ErrorUserUnauth            = errors.New("the session for this user was not found")
	ErrNotFound                = errors.New("records not found")
)
