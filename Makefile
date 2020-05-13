default: compile start

compile:
	go build -o=./bin/worker ./cmd/server/main.go

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