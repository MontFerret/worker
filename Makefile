default: compile

compile:
	go build -o=./bin/worker ./cmd/server/main.go 