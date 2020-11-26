export GO111MODULE=on

build: deps
	go build -v -o schedulemutator cmd/main.go

deps:
	go get -v ./...