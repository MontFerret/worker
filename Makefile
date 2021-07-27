VERSION ?= $(shell sh versions.sh worker)
FERRET_VERSION = $(shell sh versions.sh ferret)
DIR_BIN = ./bin

default: compile start

compile:
	go build -v -o ${DIR_BIN}/worker \
	-ldflags "-X main.version=${VERSION} -X main.ferretVersion=${FERRET_VERSION}" \
	./main.go

test:
	go test ./

start:
	./bin/worker

fmt:
	go fmt ./...