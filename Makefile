.PHUNY: build
build: 
	go build -v cmd/apiserver/main.go

.DEFAULT_GOAL := build