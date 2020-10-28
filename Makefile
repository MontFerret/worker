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

release:
ifeq ($(RELEASE_VERSION), )
	$(error "Release version is required (version=x)")
else ifeq ($(FERRET_VERSION), )
	$(error "Ferret version is required")
else ifeq ($(GITHUB_TOKEN), )
	$(error "GitHub token is required (GITHUB_TOKEN)")
else
	rm -rf ./dist && \
	git tag -a v$(RELEASE_VERSION) -m "New $(RELEASE_VERSION) version" && \
	git push origin v$(RELEASE_VERSION) && \
	goreleaser
endif