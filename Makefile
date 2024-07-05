.PHUNY: build, test, all
build: 
	go build -v cmd/apiserver/main.go

test:
	go test ./internal/store/testStore
	go test ./internal/apiserver

all: test build

.DEFAULT_GOAL := all