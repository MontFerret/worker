VERSION ?= $(shell git describe --tags --always --dirty)
FERRET_VERSION = $(shell go list -m all | grep github.com/MontFerret/ferret v| awk -F 'v' '{print $2}')
DIR_BIN = ./bin

default: compile start

compile:
	go build -v -o ${DIR_BIN}/worker \
	-ldflags "-X main.version=${VERSION}" \
	./main.go

start:
	./bin/worker

release:
ifeq ($(RELEASE_VERSION), )
	$(error "Release version is required (version=x)")
else ifeq ($(GITHUB_TOKEN), )
	$(error "GitHub token is required (GITHUB_TOKEN)")
else
	rm -rf ./dist && \
	git tag -a v$(RELEASE_VERSION) -m "New $(RELEASE_VERSION) version" && \
	git push origin v$(RELEASE_VERSION) && \
	goreleaser
endif