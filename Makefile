VERSION ?= $(shell sh versions.sh worker)
FERRET_VERSION = $(shell sh versions.sh ferret)
DIR_BIN = ./bin

default: compile start

build: vet test compile

compile:
	go build -v -o ${DIR_BIN}/worker \
	-ldflags "-X main.version=${VERSION} -X main.ferretVersion=${FERRET_VERSION}" \
	./main.go

install-tools:
	go install honnef.co/go/tools/cmd/staticcheck@latest && \
	go install golang.org/x/tools/cmd/goimports@latest && \
	go install github.com/mgechev/revive@latest

install-packages:
	go mod tidy

install: install-tools install-packages

test:
	go test ./

start:
	./bin/worker

fmt:
	go fmt ./... && \
	goimports -w -local github.com/MontFerret ./internal ./pkg main.go

lint:
	staticcheck ./... && \
	revive -config revive.toml -formatter stylish -exclude ./pkg/parser/fql/... -exclude ./vendor/... ./...

vet:
	go vet ./...
