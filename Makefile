default: compile start

compile:
	go build -o=./bin/worker ./cmd/server/main.go

start:
	./bin/worker